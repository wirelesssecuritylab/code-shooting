package sql

type configDTO struct {
	ID         string                   `yaml:"id"`
	User       string                   `yaml:"user"`
	Password   string                   `yaml:"password"`
	Host       string                   `yaml:"host"`
	Port       int                      `yaml:"port"`
	DBName     string                   `yaml:"dbName"`
	ConnParams map[string]interface{}   `yaml:"connParams"`
	Plugins    []map[string]interface{} `yaml:"plugins"`
}
