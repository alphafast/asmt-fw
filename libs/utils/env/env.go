package env

import (
	"os"
	"strconv"
)

func RequiredEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("[env.RequiredEnv]: missing required environment variable: " + key)
	}

	return val
}

func ToInt(val string) int {
	if intNum, err := strconv.Atoi(val); err != nil {
		panic("[env.ToInt]: failed to convert string to int")
	} else {
		return intNum
	}
}
