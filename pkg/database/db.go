package database

import (
	"database/sql"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpgx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strconv"
	"strings"
	"sync"
	"time"
)

// singleton instance of database connection.
var (
	dbInstance        *gorm.DB
	dbOnce            sync.Once
	templatePostgres  = "host={host} port={port} user={username} dbname={name} password={password} {options}"
	templateMysql     = "{username}:{password}@({host}:{port})/{name}?{options}"
	templateSqlServer = "sqlserver://{username}:{password}@{host}:{port}?database={name}&{options}"

	optionPlaceholders = map[string]string{
		"{username}": "username",
		"{password}": "password",
		"{host}":     "host",
		"{name}":     "name",
		"{options}":  "options",
	}
	maxOpenConn     = 25
	maxIdleConn     = 25
	maxLifeTimeConn = 300
	maxIdleTimeConn = 300
)

// Conn alias for DB().
func Conn() *gorm.DB {
	return DB()
}

// DB creates a new instance of gorm.DB if a connection is not established.
// return singleton instance.
func DB() *gorm.DB {
	if dbInstance == nil {
		dbOnce.Do(func() {
			openDB()
		})
	}
	return dbInstance
}

// openDB initialize gorm DB.
func openDB() {
	dbConfig, err := beego.AppConfig.GetSection("database")
	if err != nil {
		panic(err)
	}

	var dbDebug = true
	var logLevel = logger.Info

	if debug, err := beego.AppConfig.Bool("database::debug"); err == nil {
		dbDebug = debug
	}

	if !dbDebug {
		logLevel = logger.Silent
	}

	if db, err := sql.Open("nrpgx", buildDsn(templatePostgres, dbConfig)); err != nil {
		panic(err)
	} else {
		gormDB, err := gorm.Open(
			getDialectWithExistingConnection(dbConfig, db),
			&gorm.Config{
				SkipDefaultTransaction: true,
				PrepareStmt:            true,
				Logger:                 logger.Default.LogMode(logLevel),
			},
		)
		if err != nil {
			panic("cannot open database.")
		}
		dbInstance = gormDB
		sqlDb, err := dbInstance.DB()
		if err != nil {
			panic(err)
		}

		if parse, err := strconv.Atoi(dbConfig["maxopenconn"]); err == nil {
			maxOpenConn = parse
		}
		if parse, err := strconv.Atoi(dbConfig["maxidleconn"]); err == nil {
			maxIdleConn = parse
		}
		if parse, err := strconv.Atoi(dbConfig["maxlifetimeconn"]); err == nil {
			maxLifeTimeConn = parse
		}
		if parse, err := strconv.Atoi(dbConfig["maxidletimeconn"]); err == nil {
			maxIdleTimeConn = parse
		}
		sqlDb.SetMaxOpenConns(maxOpenConn)
		sqlDb.SetMaxIdleConns(maxIdleConn)
		sqlDb.SetConnMaxLifetime(time.Duration(maxLifeTimeConn) * time.Second)
		sqlDb.SetConnMaxIdleTime(time.Duration(maxIdleTimeConn) * time.Second)
	}
}

func getDialect(dbConfig map[string]string) gorm.Dialector {

	switch dbConfig["driver"] {
	case "postgres":
		return postgres.Open(buildDsn(templatePostgres, dbConfig))
	case "mssql":
		return mysql.Open(buildDsn(templateMysql, dbConfig))
	case "mysql":
		return sqlserver.Open(buildDsn(templateSqlServer, dbConfig))
	default:
		return postgres.Open(buildDsn(templatePostgres, dbConfig))
	}
}

func getDialectWithExistingConnection(dbConfig map[string]string, db *sql.DB) gorm.Dialector {

	switch dbConfig["driver"] {
	case "postgres":
		return postgres.New(postgres.Config{Conn: db})
	case "mssql":
		return mysql.New(mysql.Config{Conn: db})
	case "mysql":
		return sqlserver.New(sqlserver.Config{Conn: db})
	default:
		return postgres.New(postgres.Config{Conn: db})
	}
}

func buildDsn(template string, dbConfig map[string]string) string {
	for k, v := range optionPlaceholders {
		template = strings.Replace(template, k, dbConfig[v], 1)
	}
	return strings.Replace(template, "{port}", dbConfig["port"], 1)
}
