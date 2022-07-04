package database

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrConfigDriverRequired       = errors.New("config driver is required")
	ErrConfigHostRequired         = errors.New("config host is required")
	ErrConfigPortRequired         = errors.New("config port is required")
	ErrConfigUsernameRequired     = errors.New("config username is required")
	ErrConfigPasswordRequired     = errors.New("config password is required")
	ErrConfigDatabaseNameRequired = errors.New("config database name is required")
)

type DbConnection struct {
	db     *gorm.DB
	Config Config
}

func New(opts ...ConfigOption) (map[string]*gorm.DB, error) {
	conn := map[string]*gorm.DB{}

	for _, fn := range opts {
		if nil != fn {
			cfg := defaultDatabaseConfig()
			fn(&cfg)
			if err := checkRequiredDatabaseConfig(cfg.Driver, cfg.Host, cfg.Port, cfg.Username, cfg.Password,
				cfg.Name); err != nil {
				return nil, err
			}

			db, err := cfg.connectDatabase()
			if err != nil {
				return nil, err
			}
			conn[cfg.Name] = db
		}
	}



	return conn, nil
}

func (r *DbConnection) Conn() *gorm.DB {
	return r.db
}

func checkRequiredDatabaseConfig(driver, host, port, username, password, name string) error {
	if driver == "" {
		return ErrConfigDriverRequired
	}
	if host == "" {
		return ErrConfigHostRequired
	}
	if port == "" {
		return ErrConfigPortRequired
	}
	if username == "" {
		return ErrConfigUsernameRequired
	}
	if password == "" {
		return ErrConfigPasswordRequired
	}
	if name == "" {
		return ErrConfigDatabaseNameRequired
	}
	return nil
}
