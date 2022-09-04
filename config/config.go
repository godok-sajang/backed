package config

import (
	"flag"
	"fmt"
	"godok/util"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Profile string
	EnvFile string
)

func Config(key string) string {
	err := godotenv.Load(EnvFile)
	if err != nil {
		log.Fatal(err)
	}
	return os.Getenv(key)
}

func Init() {
	wordPtr := flag.String("profile", "", "")
	flag.Parse()

	Profile = "dev"
	if wordPtr != nil && util.Contain(*wordPtr, []string{"prod", "test"}) {
		Profile = fmt.Sprintf("%s", *wordPtr)
	}
	EnvFile = fmt.Sprintf("env/%s.env", Profile)

	fmt.Println("Active profile: ", Profile)
}
