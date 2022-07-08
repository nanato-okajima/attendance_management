package main

import (
	"github.com/gorilla/mux"

	"github.com/nanato-okajima/attendance_management/handler"
)

func settingRoute() *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/v1").Subrouter()

	s.HandleFunc("/attendance/register", handler.AttendanceRegisterHandler).Methods("POST")
	s.HandleFunc("/attendance/list", handler.AttendanceListHandler).Methods("GET")
	s.HandleFunc("/attendance/{id}", handler.AttendanceUpdateHandler).Methods("PUT")
	s.HandleFunc("/attendance/{id}", handler.AttendanceDeleteHandler).Methods("DELETE")

	return r
}
