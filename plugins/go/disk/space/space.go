package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
)

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

// disk usage of path/disk
func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

func main() {
	var path, in string
	var atleast, free float64
	flag.StringVar(&path, "path", "/", "The path to check disk usage for")
	flag.StringVar(&in, "in", "gb", "The size in which to measure. b|kb|mb|gb")
	flag.Float64Var(&atleast, "at-least", 1, "At least how much space is free before warning")

	flag.Parse()

	disk := DiskUsage(path)

	switch in {
	case "b":
		free = float64(disk.Free) / float64(B)
	case "kb":
		free = float64(disk.Free) / float64(KB)
	case "mb":
		free = float64(disk.Free) / float64(MB)
	default:
		free = float64(disk.Free) / float64(GB)
	}

	if float64(free) <= float64(atleast) {
		fmt.Printf("%g free disk space available. At least %g free space required", free, atleast)
		os.Exit(1)
	} else {
		fmt.Printf("%g free disk space available", free)
		os.Exit(0)
	}
}
