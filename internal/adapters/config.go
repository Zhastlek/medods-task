package adapters

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURL       string
	NameDB         string
	NameCollection string
	AccessKey      string
	RefreshKey     string
}

func Configuration() Config {
	err := godotenv.Load(filepath.Join(".", ".env"))
	if err != nil {
		fmt.Printf(".env file load error: %v\n", err)
	}
	return Config{
		MongoURL:       os.Getenv("mongoURL"),
		NameDB:         os.Getenv("nameDB"),
		NameCollection: os.Getenv("nameCollection"),
		AccessKey:      os.Getenv("access_key"),
	}
}
