package cgroups

// 配置 内存 CPU 限制
type ResourceConfig struct {
	// 内存限制 CPU核心数 CPU时间片权重
	MemoryLimit, CpuShare string
}
