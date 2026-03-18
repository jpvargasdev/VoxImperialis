package system

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

// StatusResult holds aggregated system metrics.
type StatusResult struct {
	Hostname string
	Uptime   string
	Load     string
	Memory   string
	Disk     string
}

// GetSystemStatus collects hostname, uptime, load average, memory, and disk info.
func GetSystemStatus() (StatusResult, error) {
	hostname, _ := os.Hostname()

	uptime, err := getUptime()
	if err != nil {
		uptime = "unavailable"
	}

	load, err := getLoadAvg()
	if err != nil {
		load = "unavailable"
	}

	mem, err := getMemory()
	if err != nil {
		mem = "unavailable"
	}

	disk, err := getDiskUsage()
	if err != nil {
		disk = "unavailable"
	}

	return StatusResult{
		Hostname: hostname,
		Uptime:   uptime,
		Load:     load,
		Memory:   mem,
		Disk:     disk,
	}, nil
}

// getUptime reads /proc/uptime and returns a human-readable duration.
func getUptime() (string, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "", err
	}
	var totalSec float64
	if _, err := fmt.Sscanf(strings.TrimSpace(string(data)), "%f", &totalSec); err != nil {
		return "", fmt.Errorf("parse /proc/uptime: %w", err)
	}
	total := int(totalSec)
	days := total / 86400
	hours := (total % 86400) / 3600
	minutes := (total % 3600) / 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes), nil
	}
	return fmt.Sprintf("%dh %dm", hours, minutes), nil
}

// getLoadAvg reads /proc/loadavg and returns the 1/5/15-minute averages.
func getLoadAvg() (string, error) {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return "", err
	}
	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return "", fmt.Errorf("unexpected /proc/loadavg format")
	}
	return fmt.Sprintf("%s %s %s", fields[0], fields[1], fields[2]), nil
}

// getMemory reads /proc/meminfo and returns used/total in GiB.
func getMemory() (string, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return "", err
	}
	values := make(map[string]uint64)
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSuffix(fields[0], ":")
		var val uint64
		fmt.Sscanf(fields[1], "%d", &val) //nolint:errcheck
		values[key] = val
	}
	total := values["MemTotal"]
	available := values["MemAvailable"]
	if total == 0 {
		return "", fmt.Errorf("MemTotal not found in /proc/meminfo")
	}
	used := total - available
	// /proc/meminfo values are in kB; convert to GiB
	return fmt.Sprintf("%.1fG / %.1fG",
		float64(used)/(1024*1024),
		float64(total)/(1024*1024),
	), nil
}

// getDiskUsage returns used/total GiB and percentage for the root filesystem.
func getDiskUsage() (string, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil {
		return "", err
	}
	bsize := uint64(stat.Bsize)
	total := stat.Blocks * bsize
	free := stat.Bfree * bsize
	if total == 0 {
		return "0G / 0G (0%)", nil
	}
	used := total - free
	pct := int(used * 100 / total)
	return fmt.Sprintf("%.1fG / %.1fG (%d%%)",
		float64(used)/(1024*1024*1024),
		float64(total)/(1024*1024*1024),
		pct,
	), nil
}
