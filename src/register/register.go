package register

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"business2api/src/logger"
	"business2api/src/pool"
)

// ==================== 注册与刷新 ====================

var (
	DataDir       string
	TargetCount   int
	MinCount      int
	CheckInterval time.Duration
	Threads       int
	Headless      bool   // 注册无头模式
	Proxy         string // 代理
)

var IsRegistering int32

// Stats 注册统计（导出）
var Stats = &RegisterStats{}

// 注册统计
type RegisterStats struct {
	Total     int       `json:"total"`
	Success   int       `json:"success"`
	Failed    int       `json:"failed"`
	LastError string    `json:"lastError"`
	UpdatedAt time.Time `json:"updatedAt"`
	mu        sync.RWMutex
}

func (s *RegisterStats) AddSuccess() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Total++
	s.Success++
	s.UpdatedAt = time.Now()
}

func (s *RegisterStats) AddFailed(err string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Total++
	s.Failed++
	s.LastError = err
	s.UpdatedAt = time.Now()
}

func (s *RegisterStats) Get() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return map[string]interface{}{
		"total":      s.Total,
		"success":    s.Success,
		"failed":     s.Failed,
		"last_error": s.LastError,
		"updated_at": s.UpdatedAt,
	}
}

// 注册结果
type RegisterResult struct {
	Success  bool   `json:"success"`
	Email    string `json:"email"`
	Error    string `json:"error"`
	NeedWait bool   `json:"needWait"`
}

// StartRegister 启动注册任务
func StartRegister(count int) error {
	if !atomic.CompareAndSwapInt32(&IsRegistering, 0, 1) {
		return fmt.Errorf("注册进程已在运行")
	}

	// 获取数据目录的绝对路径
	dataDirAbs, _ := filepath.Abs(DataDir)
	if err := os.MkdirAll(dataDirAbs, 0755); err != nil {
		atomic.StoreInt32(&IsRegistering, 0)
		return fmt.Errorf("创建数据目录失败: %w", err)
	}

	// 使用配置的线程数
	threads := Threads
	if threads <= 0 {
		threads = 1
	}
	for i := 0; i < threads; i++ {
		go NativeRegisterWorker(i+1, dataDirAbs)
	}

	// 监控进度
	go func() {
		for {
			time.Sleep(10 * time.Second)
			pool.Pool.Load(DataDir)
			if pool.Pool.TotalCount() >= TargetCount {
				logger.Info("✅ 已达到目标账号数: %d，停止注册", pool.Pool.TotalCount())
				atomic.StoreInt32(&IsRegistering, 0)
				return
			}
		}
	}()

	return nil
}

// PoolMaintainer 号池维护器
func PoolMaintainer() {
	interval := CheckInterval
	if interval < time.Minute {
		interval = 30 * time.Minute
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	CheckAndMaintainPool()

	for range ticker.C {
		CheckAndMaintainPool()
	}
}

// CheckAndMaintainPool 检查并维护号池
func CheckAndMaintainPool() {
	pool.Pool.Load(DataDir)

	readyCount := pool.Pool.ReadyCount()
	pendingCount := pool.Pool.PendingCount()
	totalCount := pool.Pool.TotalCount()

	logger.Info("📊 号池检查: ready=%d, pending=%d, total=%d, 目标=%d, 最小=%d",
		readyCount, pendingCount, totalCount, TargetCount, MinCount)

	if totalCount < TargetCount {
		needCount := TargetCount - totalCount
		logger.Info("⚠️ 账号数未达目标，需要注册 %d 个", needCount)
		if err := StartRegister(needCount); err != nil {
			logger.Error("❌ 启动注册失败: %v", err)
		}
	}
}
