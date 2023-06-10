package logger

type configDto struct {
	Level        string           `yaml:"level"`
	Encoder      string           `yaml:"encoder"`
	Format       string           `yaml:"format"`
	OutputPaths  []string         `yaml:"outputPaths"`
	RotateConfig *rotateConfigDto `yaml:"rotateConfig"`
}

type rotateConfigDto struct {
	MaxSize    int  `yaml:"maxSize"`
	MaxBackups int  `yaml:"maxBackups"`
	MaxAge     int  `yaml:"maxAge"`
	Compress   bool `yaml:"compress"`
	FileMode   string  `yaml:"fileMode"`
}
