package main

import (
	"github.com/gorilla/mux"

	"github.com/nanato-okajima/attendance_management/database"
	"github.com/nanato-okajima/attendance_management/handler"
	"github.com/nanato-okajima/attendance_management/service"
)

func settingRoute() *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/v1").Subrouter()

	ar := database.NewAttendanceDatabase()
	as := service.NewAttendanceService(ar)
	ah := handler.NewAttendanceHandler(as)

	s.HandleFunc("/attendance/register", ah.RegisterHandler).Methods("POST")
	s.HandleFunc("/attendance/list", ah.ListHandler).Methods("GET")
	s.HandleFunc("/attendance/{id}", ah.UpdateHandler).Methods("PUT")
	s.HandleFunc("/attendance/{id}", ah.DeleteHandler).Methods("DELETE")

	return r
}
