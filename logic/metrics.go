package logic

import (
	"bluebell/models"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// GetSystemMetrics 获取系统性能指标
func GetSystemMetrics() (*models.SystemMetrics, error) {
	metrics := &models.SystemMetrics{
		Timestamp: time.Now(),
	}

	// 1. 获取CPU使用率
	cpuUsage, err := getCPUUsage()
	if err != nil {
		zap.L().Error("getCPUUsage failed", zap.Error(err))
		cpuUsage = 0
	}
	metrics.CPUUsage = cpuUsage

	// 2. 获取内存使用率
	memoryUsage, err := getMemoryUsage()
	if err != nil {
		zap.L().Error("getMemoryUsage failed", zap.Error(err))
		memoryUsage = 0
	}
	metrics.MemoryUsage = memoryUsage

	// 3. 获取磁盘使用率
	diskUsage, err := getDiskUsage()
	if err != nil {
		zap.L().Error("getDiskUsage failed", zap.Error(err))
		diskUsage = 0
	}
	metrics.DiskUsage = diskUsage

	// 4. 获取活跃用户数（这里简化处理，统计最近5分钟有操作的用户）
	activeUsers, err := getActiveUsers(5 * time.Minute)
	if err != nil {
		zap.L().Error("getActiveUsers failed", zap.Error(err))
		activeUsers = 0
	}
	metrics.ActiveUsers = activeUsers

	return metrics, nil
}

// getCPUUsage 获取CPU使用率（从/proc/stat读取）
func getCPUUsage() (float64, error) {
	data, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return 0, fmt.Errorf("failed to read /proc/stat: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == "cpu" {
			if len(fields) < 8 {
				return 0, fmt.Errorf("invalid cpu stat format")
			}

			// 解析CPU时间
			user, _ := strconv.ParseFloat(fields[1], 64)
			nice, _ := strconv.ParseFloat(fields[2], 64)
			system, _ := strconv.ParseFloat(fields[3], 64)
			idle, _ := strconv.ParseFloat(fields[4], 64)
			iowait, _ := strconv.ParseFloat(fields[5], 64)
			irq, _ := strconv.ParseFloat(fields[6], 64)
			softirq, _ := strconv.ParseFloat(fields[7], 64)

			total := user + nice + system + idle + iowait + irq + softirq
			if total == 0 {
				return 0, nil
			}

			// 计算使用率（非空闲时间占比）
			usage := (total - idle - iowait) / total * 100
			return usage, nil
		}
	}

	return 0, fmt.Errorf("cpu stat not found")
}

// getMemoryUsage 获取内存使用率（从/proc/meminfo读取）
func getMemoryUsage() (float64, error) {
	data, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, fmt.Errorf("failed to read /proc/meminfo: %v", err)
	}

	var total, available uint64
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			total, _ = strconv.ParseUint(fields[1], 10, 64)
		case "MemAvailable:":
			available, _ = strconv.ParseUint(fields[1], 10, 64)
		}
	}

	if total == 0 {
		return 0, fmt.Errorf("failed to get memory info")
	}

	// 计算使用率
	used := total - available
	usagePercent := float64(used) / float64(total) * 100
	return usagePercent, nil
}

// getDiskUsage 获取磁盘使用率
func getDiskUsage() (float64, error) {
	// 获取当前工作目录的磁盘使用情况
	wd, err := os.Getwd()
	if err != nil {
		return 0, err
	}

	var stat syscall.Statfs_t
	if err := syscall.Statfs(wd, &stat); err != nil {
		return 0, err
	}

	// 计算磁盘使用率
	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)
	used := total - free

	if total == 0 {
		return 0, fmt.Errorf("invalid disk stats")
	}

	usagePercent := float64(used) / float64(total) * 100
	return usagePercent, nil
}

// getActiveUsers 获取活跃用户数（最近指定时间内有操作的用户）
func getActiveUsers(duration time.Duration) (int64, error) {
	// 这里简化处理，实际应用中应该：
	// 1. 从Redis获取在线用户列表
	// 2. 或者从数据库查询最近活跃的用户
	// 3. 或者使用WebSocket连接数

	// 临时返回一个模拟值
	return 42, nil
}

// GetSystemMetricsHistory 获取系统性能历史数据（简化版）
func GetSystemMetricsHistory(startTime, endTime int64, metricType string) ([]interface{}, error) {
	// 这里简化处理，返回模拟的历史数据
	// 实际应用中应该：
	// 1. 从数据库查询历史记录
	// 2. 或者从时序数据库（如InfluxDB）查询
	// 3. 或者从Prometheus查询

	// 生成模拟数据
	var data []interface{}
	interval := (endTime - startTime) / 10 // 分成10个点

	for i := 0; i < 10; i++ {
		timestamp := startTime + int64(i)*interval
		value := 20.0 + float64(i)*2.0 // 模拟递增的数据

		data = append(data, map[string]interface{}{
			"timestamp": timestamp,
			"value":     value,
		})
	}

	return data, nil
}