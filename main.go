package main

import (
	"log"
	"net/http"
	"os"

	syslutil "github.com/anz-bank/new-sysl-playground/syslUtil"

	"github.com/gorilla/mux"
)

func compile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	args := []string{"sd", "-o", "./tmp/call-login-sequence.png", "-s", "MobileApp <- Login", "./tmp/hello.sysl"}
	syslutil.Execute(args)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "post called"}`))
}

func main() {
	port := os.Getenv("PORT")
	r := mux.NewRouter()
	api := r.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/compile", compile).Methods(http.MethodPost)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(".")))
	log.Printf("Server is running on: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
