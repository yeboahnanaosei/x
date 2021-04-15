package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(":9900", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		body := []byte("Method Not Allowed")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		w.Write(body)
		return
	}

	userName := r.FormValue("ProfileName")
	userPhone := r.FormValue("WaId")
	w.Write([]byte("Hi, "+userName+"! Welcome to RTS Solutions. "+ userPhone))
}