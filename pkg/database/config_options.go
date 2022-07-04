package database

type ConfigOption func(*Config)

func ConfigDriverName(driverName string) ConfigOption {
	return func(cfg *Config) {
		cfg.Driver = driverName
		cfg.TemplateDsn = templateDsn[driverName]
	}
}

func ConfigHost(host string) ConfigOption {
	return func(cfg *Config) { cfg.Host = host }
}

func ConfigPort(port string) ConfigOption {
	return func(cfg *Config) { cfg.Port = port }
}

func ConfigUsername(username string) ConfigOption {
	return func(cfg *Config) { cfg.Username = username }
}

func ConfigPassword(password string) ConfigOption {
	return func(cfg *Config) { cfg.Password = password }
}

func ConfigDebugEnabled(enabled bool) ConfigOption {
	return func(cfg *Config) { cfg.Debug = enabled }
}

func ConfigMaxOpenConnection(value int) ConfigOption {
	return func(cfg *Config) { cfg.MaxOpenConnection = value }
}

func ConfigMaxIdleConnection(value int) ConfigOption {
	return func(cfg *Config) { cfg.MaxIdleConnection = value }
}

func ConfigMaxLifeTimeConnection(value int) ConfigOption {
	return func(cfg *Config) { cfg.MaxLifeTimeConnection = value }
}

func ConfigMaxIdleTimeConnection(value int) ConfigOption {
	return func(cfg *Config) { cfg.MaxIdleTimeConnection = value }
}



