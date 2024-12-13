package filetools

import (
	"log"
	"os"
)

func CurrentLocalPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
