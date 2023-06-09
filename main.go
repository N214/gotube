package main

import (
	"os"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	// logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)


	// Initialize a new instance of our application struct, containing the // dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog, 
	}

	// router
	r := mux.NewRouter()
	r.HandleFunc("/webhooks", app.webhook).Methods("POST")
	port := ":8089"


	// Better to create a new server and use our own error log logger with
	svr := &http.Server{
		Addr: port,
		ErrorLog: errorLog,
		Handler: r,
	}

	infoLog.Printf("Starting server on %s", port)
	if err := svr.ListenAndServe(); err != nil {
		errorLog.Fatal("Error starting server: ", err)
	}
}