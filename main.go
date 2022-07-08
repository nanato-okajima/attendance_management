package main

import (
	"log"
	"net/http"

	"github.com/nanato-okajima/attendance_management/database"
)

func main() {
	router := settingRoute()

	if err := database.SetupDB(); err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}
