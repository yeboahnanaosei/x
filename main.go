package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", index)
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "5000"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
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
	w.Write([]byte("Hi, " + userName + "! Welcome to RTS Solutions. " + userPhone))
}
