package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	syslutil "github.com/anz-bank/new-sysl-playground/syslUtil"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type compileContent struct {
	Code    string
	Command []string
}

func compile(w http.ResponseWriter, r *http.Request) {
	var c compileContent
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("FATAL IO reader issue %s ", err.Error())
	}
	json.Unmarshal(body, &c)

	args := c.Command
	compileResult, err := syslutil.Execute(c.Code, args)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(compileResult); err != nil {
		panic(err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3030"
	}
	r := mux.NewRouter()
	api := r.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/compile", compile).Methods(http.MethodPost)
	api.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message": "pong"}`))
	}).Methods(http.MethodGet)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	handler := c.Handler(r)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./client")))
	log.Printf("Server is running on: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
