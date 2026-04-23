package controllers

import (
	"net/http"
	"zerotrusterp/core"
	"os"
	"syscall"
	"fmt"
	"bufio"
	"strings"
	"time"
)

func CheckRAM() string {
	var info syscall.Sysinfo_t
	err := syscall.Sysinfo(&info)
	if err != nil {
		return fmt.Sprintf("Error getting system info: %s", err)
	}

	// info.Totalram and info.Freeram are given in units of info.Unit (usually bytes)
	unit := uint64(info.Unit)
	
	total := (uint64(info.Totalram) * unit)
	free := (uint64(info.Freeram) * unit)

	// Convert bytes to Gigabytes (GB)
	totalGB := float64(total) / 1024 / 1024 / 1024
	freeGB := float64(free) / 1024 / 1024 / 1024

	return fmt.Sprintf("Total RAM: %.2f GB | Free RAM: %.2f GB", totalGB, freeGB)
}

func CheckStorage() string {
	path, err := os.Getwd()
	if err != nil {
		return fmt.Sprintf("Error getting path: %s", err)
	}

	var stat syscall.Statfs_t
	err = syscall.Statfs(path, &stat)
	if err != nil {
		return fmt.Sprintf("Error calling Statfs: %s", err)
	}

	// Total size = blocks * block size
	total := stat.Blocks * uint64(stat.Bsize)
	// Available size = available blocks * block size
	free := stat.Bavail * uint64(stat.Bsize)

	// Use Sprintf to return the formatted string
	return fmt.Sprintf("Total storage: %.2f GB | Free storage: %.2f GB", 
		float64(total)/1024/1024/1024, 
		float64(free)/1024/1024/1024)
}

func CheckCPU() string {
	// 1. Get CPU Name (Model)
	cpuName := "Unknown"
	file, err := os.Open("/proc/cpuinfo")
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "model name") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					cpuName = strings.TrimSpace(parts[1])
					break // Found it, stop scanning
				}
			}
		}
		file.Close()
	}

	// 2. Get CPU Usage (%)
	// We read /proc/stat twice with a small delay to calculate the delta
	idle1, total1 := getCPUSample()
	time.Sleep(500 * time.Millisecond)
	idle2, total2 := getCPUSample()

	idleTicks := float64(idle2 - idle1)
	totalTicks := float64(total2 - total1)
	usage := 100 * (totalTicks - idleTicks) / totalTicks

	return fmt.Sprintf("CPU: %s | Usage: %.2f%%", cpuName, usage)
}

func getCPUSample() (idle, total uint64) {
	contents, err := os.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == "cpu" {
			var user, nice, system, idleTick, iowait, irq, softirq uint64
			fmt.Sscanf(line, "cpu %d %d %d %d %d %d %d", &user, &nice, &system, &idleTick, &iowait, &irq, &softirq)
			idle = idleTick + iowait
			total = user + nice + system + idle + irq + softirq
			return
		}
	}
	return
}

func CheckDatabase(){
	
}

func Health(w http.ResponseWriter, r *http.Request){

	data := map[string]interface{}{
			"Title": "Health",
			"DatabaseName":core.GetDatabaseName(),
			"Storage":CheckStorage(),
			"RAM": CheckRAM(),
			"CPU":CheckCPU(),
			"TotalUsers":core.GetCountRecords("users"),
		}

	core.RenderPage(w,r, "apps/system/views/health.html", data)
}