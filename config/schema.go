package config

type Service_ struct {
	Name     string `yaml:name`
	TTL      int64  `yaml:ttl`
	Interval int64  `yaml:interval`
	Address  string `yaml:address`
}

type Logger_ struct {
	Level string `yaml:level`
	Dir   string `yaml:dir`
}

type ConfigSchema_ struct {
	Service  Service_  `yaml:service`
	Logger   Logger_   `yaml:logger`
	Database Database_ `yaml:database`
}

type SQLite_ struct {
	Path string `yaml:path`
}

type MySQL_ struct {
	Address      string `yaml:address`
	User         string `yaml:user`
	Password     string `yaml:password`
	DB           string `yaml:db`
	MaxIdleTime  int    `yaml:maxIdelTime`
	MaxLifeTime  int    `yaml:maxLifeTime`
	MaxIdleConns int    `yaml:maxIdleConns`
	MaxOpenConns int    `yaml:maxOpenConns`
}

type Database_ struct {
	Driver string  `yaml:driver`
	MySQL  MySQL_  `yaml:mysql`
	SQLite SQLite_ `yaml:sqlite`
}
