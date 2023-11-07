package config

type Config interface {
	Validate() error
	Complete() error
}

type Registration struct {
	Scenario string     `koanf:"scenario,omitempty"`
	Etcd     EtcdConfig `koanf:"etcd,omitempty"`
}

type MongoConfig struct {
	Uri      string `koanf:"uri"`
	Database string `koanf:"database"`
}

type LogConfig struct {
	Enabled       bool   `koanf:"enabled"`
	LogLevel      string `koanf:"logLevel"`
	LogPath       string `koanf:"logPath"`
	OutputConsole bool   `koanf:"outputToConsole"`
	FileName      string `koanf:"fileName"`
	MaxSizeInMB   int    `koanf:"maxSizeInMB"`
	MaxAgeInDays  int    `koanf:"maxAgeInDays"`
	MaxBackups    int    `koanf:"maxBackups"`
	Compress      bool   `koanf:"compress"`
}

type EtcdConfig struct {
	RefreshSeconds        uint     `koanf:"refreshSeconds"`
	ConnectTimeoutSeconds uint     `koanf:"connectTimeoutSeconds"`
	Endpoints             []string `koanf:"endpoints"`
}

type HttpSetting struct {
	Port    uint   `koanf:"port"`
	Address string `koanf:"address"`
	Proxy   string `koanf:"proxy"`
}

type TaskPoolSetting struct {
	Capacity int `koanf:"capacity"`
}

type RedisConfig struct {
	Address                  string `koanf:"address,omitempty"`
	Password                 string `koanf:"password,omitempty"`
	DefaultDb                int    `koanf:"defaultDb,omitempty"`
	PoolSize                 int    `koanf:"poolSize,omitempty"`
	PoolTimeout              int    `koanf:"poolTimeout"`
	ReadTimeout              int    `koanf:"readTimeout"`
	WriteTimeout             int    `koanf:"writeTimeout"`
	AutoCreateConsumerGroups bool   `koanf:"autoCreateConsumerGroups"`
}

type ServerConfig struct {
	ApplicationName string           `koanf:"applicationName"`
	Http            *HttpSetting     `koanf:"http"`
	Registration    *Registration    `koanf:"registration"`
	Redis           *RedisConfig     `koanf:"redis"`
	Mongo           *MongoConfig     `koanf:"mongodb"`
	LogSetting      *LogConfig       `koanf:"logConfig"`
	TaskPoolSetting *TaskPoolSetting `koanf:"taskPool"`
}

func (s ServerConfig) GetServerConfig() *ServerConfig {
	return &s
}
func (s ServerConfig) Validate() error {
	return nil
}
func (s ServerConfig) Complete() error {
	return nil
}
