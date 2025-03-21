package configuration

import (
	"os"
	"strconv"
)

type Config struct{
    APIPort		string	    //port to expose
    DB			DBConfig
    Env			string	    //dev or prod
}

type DBConfig struct{
    Addr    string
    Token   string
    Org     string
    Bucket  string
}

func GetConfig() Config {
	return Config{
		APIPort: GetString("API_PORT", ":8080"),
		DB: DBConfig{
			Addr:   GetString("DB_ADDR", "http://localhost:8086"),
			Token:  GetString("DB_TOKEN", "mytoken"),
			Org:    GetString("DB_ORG", "my-org"),
			Bucket: GetString("DB_BUCKET", "my-bucket"),
		},
		Env: GetString("ENV", "prod"),
	}
}

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return valAsInt
}

func GetBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}

	return boolVal
}
