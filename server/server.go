package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"time"

	"github.com/itomofumi/go-gin-xorm-starter/factory"
	"github.com/itomofumi/go-gin-xorm-starter/model"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/itomofumi/go-gin-xorm-starter/infra"
	"github.com/itomofumi/go-gin-xorm-starter/util"
	"github.com/itomofumi/gognito/auth"
	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"xorm.io/xorm"
)

const (
	ipEnv                = "IP"
	portEnv              = "PORT"
	shutdownTimeoutEnv   = "SHUTDOWN_TIMEOUT"
	cognitoRegionEnv     = "COGNITO_REGION"
	cognitoUserPoolIDEnv = "COGNITO_USER_POOL_ID"
)

// DB Engine の初期化
func setupDBEngine(logLevel logrus.Level) (*xorm.Engine, error) {
	dbOptions := infra.LoadMySQLConfigEnv()
	log.Printf("MySQL Connection String: %v", dbOptions.FormatDSN())
	engine, err := infra.InitMySQLEngine(dbOptions)
	if err != nil {
		return nil, err
	}
	sqlLogWriter := &lumberjack.Logger{
		Filename:   path.Join(os.Getenv("LOG_DIR"), "server_sql.log"),
		MaxSize:    10,   // megabytes
		MaxBackups: 100,  // default: not to remove old logs
		Compress:   true, // disabled by default
	}

	loggerSQL := logrus.New()
	loggerSQL.Level = logLevel
	loggerSQL.Out = io.MultiWriter(os.Stdout, sqlLogWriter)

	engine.SetLogger(engine.Logger())

	return engine, nil
}

// Start starts api server
// func Start(serverOptions Options) error {
func Start() error {
	logger := util.GetLogger()

	// ログの出力設定
	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = logrus.DebugLevel
		logger.Warnf("LOG_LEVEL is not set.")
	}

	logDir := os.Getenv("LOG_DIR")

	// db engine 初期化
	engine, err := setupDBEngine(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Println("engine closed")
		engine.Close()
	}()

	// key-value store initialization.
	kvsClient := infra.NewKVSClient()
	defer func() {
		log.Println("kvsClient closed")
		kvsClient.Close()
	}()

	accessLogWriter := &lumberjack.Logger{
		Filename:   path.Join(logDir, "server_access.log"),
		MaxSize:    10,   // megabytes
		MaxBackups: 100,  // default: not to remove old logs
		Compress:   true, // disabled by default
	}

	loggerAccess := logrus.New()
	loggerAccess.Level = logLevel
	loggerAccess.Out = io.MultiWriter(os.Stdout, accessLogWriter)

	// Gin エラーログ
	ginErrorLogWriter := &lumberjack.Logger{
		Filename:   path.Join(logDir, "error.log"),
		MaxSize:    10,   // megabytes
		MaxBackups: 100,  // default: not to remove old logs
		Compress:   true, // disabled by default
	}

	gin.DefaultErrorWriter = io.MultiWriter(os.Stderr, ginErrorLogWriter)

	// service factoryの初期化
	factory := factory.NewService(engine, kvsClient)

	// override gin validator
	binding.Validator = &model.StructValidator{}

	// Ginの初期化
	r := gin.Default()

	// middlewareのロード
	r.Use(LogMiddleware(loggerAccess, time.RFC3339, false))
	r.Use(CORSMiddleware())
	r.Use(ServiceKeyMiddleware(factory))

	// auth middlewareの準備
	authenticator, err := auth.New(
		&auth.UserPool{
			Region: os.Getenv(cognitoRegionEnv),
			PoolID: os.Getenv(cognitoUserPoolIDEnv),
		},
		&auth.Option{
			NoVerification: DisableVerification,
		})

	if err != nil {
		return err
	}

	r.Use(SetAuth(authenticator))

	defineRoutes(r)

	ip := os.Getenv(ipEnv)

	port := os.Getenv(portEnv)
	if port == "" {
		port = "3000"
	}

	// parse SHUTDOWN_TIMEOUT ENV
	var shutdownTimeout int
	shutdownTimeoutStr := os.Getenv(shutdownTimeoutEnv)
	if shutdownTimeout, err = strconv.Atoi(shutdownTimeoutStr); err != nil {
		logger.Warnf("%v expects int value, but %v was given.", shutdownTimeoutEnv, shutdownTimeoutStr)
		logger.Infof("use default 5 [sec] for %v", shutdownTimeoutEnv)
		shutdownTimeout = 5
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", ip, port),
		Handler: r,
	}

	go func() {
		// Start server
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for "interrupt" or "kill" signal to gracefully shutdown.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	sig := <-quit

	logger.Printf("Shutdown Server with Signal %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeout)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown:", err)
	}
	logger.Println("Server exiting")

	return nil
}
