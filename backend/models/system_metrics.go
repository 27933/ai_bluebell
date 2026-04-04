package models

import "time"

// SystemMetrics 系统性能指标
type SystemMetrics struct {
	CPUUsage    float64   `json:"cpu_usage"`    // CPU使用率(百分比)
	MemoryUsage float64   `json:"memory_usage"` // 内存使用率(百分比)
	DiskUsage   float64   `json:"disk_usage"`   // 磁盘使用率(百分比)
	ActiveUsers int64     `json:"active_users"` // 活跃用户数
	Timestamp   time.Time `json:"timestamp"`    // 记录时间
}

// SystemMetricsRequest 系统指标查询请求
type SystemMetricsRequest struct {
	StartTime int64  `json:"start_time" form:"start_time"` // 开始时间戳(秒)
	EndTime   int64  `json:"end_time" form:"end_time"`     // 结束时间戳(秒)
	MetricType string `json:"metric_type" form:"metric_type"` // 指标类型(cpu/memory/disk)
}