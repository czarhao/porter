package utils

type Container struct {
	Name string `yaml:"name"`
	Tty  bool   `yaml:"tty"`
}

type Image struct {
	Path  string `yaml:"path"`
	Run   string `yaml:"run"`
	Layer *Layer `yaml:"layer"`
}

type Limit struct {
	MemoryLimit string `yaml:"memory"`
	CpuShare    string `yaml:"cpu_share"`
}

type Layer struct {
	RootURL   string `yaml:"root_url"`
	MntURL    string `yaml:"mnt_url"`
	WriterUrl string `yaml:"writer_url"`
	Volume    string `yaml:"volume"`
}

type Config struct {
	Container *Container `yaml:"container"`
	Image     *Image     `yaml:"image"`
	Limit     *Limit     `yaml:"limit"`
}
