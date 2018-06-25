package infra

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	// load mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

// MySQLConnectionOptions defines mysql options
type MySQLConnectionOptions struct {
	Database  string
	Host      string
	Port      uint16
	User      string
	Password  string
	Charset   string
	Collation string
}

// InitMySQLEngine initialize xorm engine for mysql
func InitMySQLEngine(options *MySQLConnectionOptions) (*xorm.Engine, error) {
	if options == nil {
		options = NewMySQLConnectionOptionsWithENV()
	}
	engine, err := xorm.NewEngine("mysql", options.String())
	if err != nil {
		return nil, err
	}

	engine.SetMapper(core.GonicMapper{})
	engine.ShowSQL(true)
	engine.Charset(options.Charset)
	engine.StoreEngine("InnoDb")

	logLevel, err := parseLogLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		return nil, err
	}
	engine.SetLogLevel(logLevel)

	return engine, nil
}

// String builds mysql connection string
func (options MySQLConnectionOptions) String() string {
	if options.User == "" {
		options.User = "root"
	}
	if options.Port == 0 {
		options.Port = 3306
	}
	if options.Host == "" {
		options.Host = "localhost"
	}

	connStr := make([]byte, 0)
	connStr = append(connStr, options.User...)
	if options.Password != "" {
		connStr = append(connStr, ':')
		connStr = append(connStr, options.Password...)
	}
	connStr = append(connStr, "@tcp("...)
	connStr = append(connStr, options.Host...)
	connStr = append(connStr, ':')
	connStr = append(connStr, strconv.Itoa(int(options.Port))...)
	connStr = append(connStr, ")/"...)

	if options.Database != "" {
		connStr = append(connStr, options.Database...)
	}

	if options.Charset != "" || options.Collation != "" {
		connStr = append(connStr, '?')

		params := make([]string, 0)
		if options.Charset != "" {
			params = append(params, "charset="+options.Charset)
		}
		if options.Collation != "" {
			params = append(params, "collation="+options.Collation)
		}

		connStr = append(connStr, strings.Join(params, "&")...)
	}
	return string(connStr)
}

// NewMySQLConnectionOptionsWithENV loads mysql options from ENVIRONMENT VARIABLE
func NewMySQLConnectionOptionsWithENV() *MySQLConnectionOptions {
	return &MySQLConnectionOptions{
		Host:      os.Getenv("DATABASE_HOST"),
		Database:  os.Getenv("DATABASE_NAME"),
		User:      os.Getenv("DATABASE_USER"),
		Password:  os.Getenv("DATABASE_PASSWORD"),
		Charset:   "utf8mb4",
		Collation: "utf8mb4_general_ci",
	}
}

// ParseLogLevel parses level string into xorm's LogLevel
func parseLogLevel(lvl string) (core.LogLevel, error) {
	switch strings.ToLower(lvl) {
	case "panic", "fatal", "error":
		return core.LOG_ERR, nil
	case "warn", "warning":
		return core.LOG_WARNING, nil
	case "info":
		return core.LOG_INFO, nil
	case "debug":
		return core.LOG_DEBUG, nil
	}
	return core.LOG_DEBUG, fmt.Errorf("cannot parse \"%v\" into go-xorm/core.LogLevel", lvl)
}
