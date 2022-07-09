package main

import (
	"log"
	"net/http"

	"github.com/nanato-okajima/attendance_management/database"
	"github.com/nanato-okajima/attendance_management/logger"
)

func main() {
	router := settingRoute()
	logger.SetupLogger()

	if err := database.SetupDB(); err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}
