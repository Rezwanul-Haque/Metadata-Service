package helpers

import (
	"github.com/joho/godotenv"
	"github.com/rezwanul-haque/Metadata-Service/src/logger"
	"log"
	"os"
)

func IsEmpty(value string) bool {
	if value == "" {
		return true
	}
	return false
}

func ByteEmpty(s []byte) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}

func Paginate(pageNum int, pageSize int, sliceLength int) (int, int) {
	start := (pageNum - 1) * pageSize

	if start > sliceLength {
		start = sliceLength
	}

	end := pageNum * pageSize
	if end > sliceLength {
		end = sliceLength
	}

	return start, end
}

func GoDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		logger.Error("Error loading .env file", err)
		foundKey := os.Getenv(key)
		if foundKey != "" {
			return foundKey
		} else {
			log.Fatalf("Required environment variables not found.")
		}
	}

	return os.Getenv(key)
}
