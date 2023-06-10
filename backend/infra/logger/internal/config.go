package internal

const (
	DefaultMaxSize    = 10 // MB
	DefaultMaxBackups = 2
	DefaultFileMode   = "0640"
)

type Config struct {
	Level        Level
	Encoder      Encoder
	Format       CustomFormat
	OutputPaths  []string
	RotateConfig RotateConfig
}

type RotateConfig struct {
	Disable    bool
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	FileMode   string
}
