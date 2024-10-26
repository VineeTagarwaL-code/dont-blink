package utils

import (
	"log"
	"net/http"
)

func CheckOrigin(r *http.Request) bool {
	return true
}

func PrintError(err string) {
	log.Print(err)
}
