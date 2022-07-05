package main

import (
	"log"
	"net/http"
)

func main() {
	settingRoute()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
