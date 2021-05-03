package ngin

import "fmt"

const (
	KB = 1 << 10 // 1 KB
	MB = 1 << 20 // 1 MB
	GB = 1 << 30 // 1 GB
)

var (
	KB_STR = fmt.Sprintf("kb: %db\n", KB)
	MB_STR = fmt.Sprintf("mb: %dkb, %db\n", MB/KB, MB)
	GB_STR = fmt.Sprintf("gb: %dmb, %dkb, %db\n", GB/MB, GB/KB, GB)
)
