package infra

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

// LoadMySQLConfigEnv initializes MySQL config using Environment Variables.
func LoadMySQLConfigEnv() *mysql.Config {
	conf := &mysql.Config{
		Net:                  "tcp",
		Addr:                 os.Getenv("DATABASE_HOST"),
		DBName:               os.Getenv("DATABASE_NAME"),
		User:                 os.Getenv("DATABASE_USER"),
		Passwd:               os.Getenv("DATABASE_PASSWORD"),
		AllowNativePasswords: true,
	}
	return conf
}

// InitMySQLEngine initialize xorm engine for mysql
func InitMySQLEngine(conf *mysql.Config) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("mysql", conf.FormatDSN())
	if err != nil {
		return nil, err
	}

	charset, ok := conf.Params["charset"]
	if !ok {
		charset = "utf8mb4"
	}

	engine.Charset(charset)
	engine.SetMapper(core.GonicMapper{})
	engine.ShowSQL(true)
	engine.StoreEngine("InnoDb")

	logLevel, err := parseLogLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		return nil, err
	}
	engine.SetLogLevel(logLevel)

	return engine, nil
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

// RunSQLFile runs sql file.
func RunSQLFile(mysqlConnectionString, sqlFilepath string) error {

	var err error
	engine, err := xorm.NewEngine("mysql", mysqlConnectionString)
	if err != nil {
		return err
	}
	defer engine.Close()

	engine.SetConnMaxLifetime(time.Second)
	engine.ShowSQL(true)
	engine.Logger().SetLevel(core.LOG_INFO)

	err = engine.Ping()
	if err != nil {
		return err
	}

	file, err := os.Open(sqlFilepath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = engine.Import(file)

	if err != nil {
		if err.Error() == "not an error" {
			err = nil
		} else {
			return err
		}
	}

	return nil
}

// EscapeMySQLString prevents from SQL-injection.
func EscapeMySQLString(value string) string {
	replace := map[string]string{
		"\\":   "\\\\",
		"'":    `\'`,
		"\\0":  "\\\\0",
		"\n":   "\\n",
		"\r":   "\\r",
		`"`:    `\"`,
		"\x1a": "\\Z",
	}

	for b, a := range replace {
		value = strings.Replace(value, b, a, -1)
	}

	return value
}
