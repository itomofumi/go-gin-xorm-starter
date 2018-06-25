package util

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv apply .env to ENVIRONMENT VARIABLE.
func LoadEnv() {
	// just print messages to stdout before logger setup.
	if _, found := os.LookupEnv("ENV_FILE"); !found {
		fmt.Println("no env file specified. try to load default .env.")
		os.Setenv("ENV_FILE", ".env")
	}

	envfile := os.Getenv("ENV_FILE")
	err := godotenv.Load(envfile)
	if err != nil {
		fmt.Printf("no env file loaded %v\n", err)
	} else {
		fmt.Printf("env file loaded: %v\n", envfile)
	}
}
