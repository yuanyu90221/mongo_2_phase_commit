package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func GetConnectString() string {
	DB_HOST := os.Getenv("MONGO_HOST")
	DB_PORT := os.Getenv("PORT")
	DB_NAME := os.Getenv("DB_NAME")
	DB_USERBAME := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	DB_PASSWD := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", DB_USERBAME, DB_PASSWD, DB_HOST, DB_PORT, DB_NAME)
}
