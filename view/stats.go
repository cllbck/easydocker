//Docker
//Copyright 2012-2017 Docker, Inc.
//
//This product includes software developed at Docker, Inc. (https://www.docker.com).
//
//This product contains software (https://github.com/kr/pty) developed
//by Keith Rarick, licensed under the MIT License.
//
//The following is courtesy of our legal counsel:
//
//
//Use and transfer of Docker may be subject to certain restrictions by the
//United States and other governments.
//It is your responsibility to ensure that your use and/or transfer does not
//violate applicable laws.
//
//For more information, please see https://www.bis.doc.gov
//
//See also https://www.apache.org/dev/crypto.html and/or seek legal counsel.

package view

import (
	docker "github.com/fsouza/go-dockerclient"
	"runtime"
)

func calculateCPUPercentUnix(stats docker.Stats) float64 {
	var (
		previousCPU    = stats.PreCPUStats.CPUUsage.TotalUsage
		previousSystem = stats.PreCPUStats.SystemCPUUsage
		cpuPercent     = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		// calculate the change for the entire system between readings
		systemDelta = float64(stats.CPUStats.SystemCPUUsage) - float64(previousSystem)
		onlineCPUs  = float64(stats.CPUStats.OnlineCPUs)
	)

	if onlineCPUs == 0.0 {
		onlineCPUs = float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
	}
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}
	return cpuPercent
}

func calculateCPUPercentWindows(stats docker.Stats) float64 {
	// Max number of 100ns intervals between the previous time read and now
	possIntervals := uint64(stats.Read.Sub(stats.PreRead).Nanoseconds()) // Start with number of ns intervals
	possIntervals /= 100                                                 // Convert to number of 100ns intervals
	possIntervals *= uint64(stats.NumProcs)                              // Multiple by the number of processors

	// Intervals used
	intervalsUsed := stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage

	// Percentage avoiding divide-by-zero
	if possIntervals > 0 {
		return float64(intervalsUsed) / float64(possIntervals) * 100.0
	}
	return 0.00
}

func calculateBlockIOUnix(stats docker.Stats) (float64, float64) {
	var blkRead, blkWrite uint64
	for _, bioEntry := range stats.BlkioStats.IOServiceBytesRecursive {
		if len(bioEntry.Op) == 0 {
			continue
		}
		switch bioEntry.Op[0] {
		case 'r', 'R':
			blkRead = blkRead + bioEntry.Value
		case 'w', 'W':
			blkWrite = blkWrite + bioEntry.Value
		}
	}
	return float64(blkRead), float64(blkWrite)
}

func calculateNetwork(stats docker.Stats) (float64, float64) {
	var rx, tx float64

	for _, v := range stats.Networks {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
	}
	return rx, tx
}

// calculateMemUsageUnixNoCache calculate memory usage of the container.
// Page cache is intentionally excluded to avoid misinterpretation of the output.
func calculateMemUsageUnixNoCache(stats docker.Stats) float64 {
	return float64(stats.MemoryStats.Usage - stats.MemoryStats.Stats.Cache)
}

func calculateMemPercentUnixNoCache(limit float64, usedNoCache float64) float64 {
	// MemoryStats.Limit will never be 0 unless the container is not running and we haven't
	// got any data from cgroup
	if limit != 0 {
		return usedNoCache / limit * 100.0
	}
	return 0
}

func calculateMem(stats docker.Stats) float64 {
	if runtime.GOOS == "windows" {
		return float64(stats.MemoryStats.PrivateWorkingSet)
	} else {
		return calculateMemUsageUnixNoCache(stats)
	}
}

func calculateMemoryPercentage(stats docker.Stats) float64 {
	var memPercent float64
	if runtime.GOOS != "windows" {
		mem := calculateMemUsageUnixNoCache(stats)
		memLimit := calculateMemoryLimit(stats)
		memPercent = calculateMemPercentUnixNoCache(memLimit, mem)
	}
	return memPercent
}

func calculateMemoryLimit(stats docker.Stats) float64 {
	var memLimit float64
	if runtime.GOOS != "windows" {
		memLimit = float64(stats.MemoryStats.Limit)
	}
	return memLimit
}

func calculateCpuPercent(stats docker.Stats) float64 {
	var cpuPercent float64
	if runtime.GOOS != "windows" {
		cpuPercent = calculateCPUPercentUnix(stats)
	} else {
		cpuPercent = calculateCPUPercentWindows(stats)
	}
	return cpuPercent * 100
}

func calculateBlockIO(stats docker.Stats) (float64, float64) {
	if runtime.GOOS != "windows" {
		return calculateBlockIOUnix(stats)
	} else {
		return float64(stats.StorageStats.ReadSizeBytes), float64(stats.StorageStats.WriteSizeBytes)
	}
}
