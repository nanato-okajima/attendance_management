package main

import (
	"github.com/gorilla/mux"

	"github.com/nanato-okajima/attendance_management/handler"
)

func settingRoute() *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/v1").Subrouter()

	s.HandleFunc("/attendance/register", handler.RegisterHandler).Methods("POST")
	s.HandleFunc("/attendance/list", handler.ListHandler).Methods("GET")
	s.HandleFunc("/attendance/{id}", handler.UpdateHandler).Methods("PUT")
	s.HandleFunc("/attendance/{id}", handler.DeleteHandler).Methods("DELETE")

	return r
}
