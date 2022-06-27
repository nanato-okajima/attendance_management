package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, "ほげほげ"); err != nil {
			log.Fatalf("起動に失敗 %#v", err)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
