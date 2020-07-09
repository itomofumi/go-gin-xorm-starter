package infra

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"xorm.io/core"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
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
	showSQL := os.Getenv("SHOW_SQL")
	if showSQL == "0" || showSQL == "false" {
		engine.ShowSQL(false)
	} else {
		engine.ShowSQL(true)
	}

	logLevel, err := parseLogLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		return nil, err
	}
	engine.SetLogLevel(logLevel)

	return engine, nil
}

// parseLogLevel parses level string into xorm's LogLevel
func parseLogLevel(lvl string) (log.LogLevel, error) {
	switch strings.ToLower(lvl) {
	case "panic", "fatal", "error":
		return log.LOG_ERR, nil
	case "warn", "warning":
		return log.LOG_WARNING, nil
	case "info":
		return log.LOG_INFO, nil
	case "debug":
		return log.LOG_DEBUG, nil
	}
	return log.LOG_DEBUG, fmt.Errorf("cannot parse \"%v\" into go-xorm/core.LogLevel", lvl)
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
	engine.ShowSQL(false)
	engine.Logger().SetLevel(log.LOG_WARNING)

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

var escapeReplace = []struct {
	Key      string
	Replaced string
}{
	{"\\", "\\\\"},
	{`'`, `\'`},
	{"\\0", "\\\\0"},
	{"\n", "\\n"},
	{"\r", "\\r"},
	{`"`, `\"`},
	{"\x1a", "\\Z"},
}

// EscapeMySQLString prevents from SQL-injection.
func EscapeMySQLString(value string) string {
	for _, r := range escapeReplace {
		value = strings.Replace(value, r.Key, r.Replaced, -1)
	}

	return value
}
