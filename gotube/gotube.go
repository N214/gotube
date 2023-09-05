package gotube

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	//"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gorilla/mux"
)

// Not used anymore since we use the default functional framework from google
func Run() (err error) {
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

	infoLog.Printf("Starting server on localhost%s", port)

	if err := svr.ListenAndServe(); err != nil {
		errorLog.Fatal("Error starting server: ", err)
		return err
	}
	return
}

func init() {
	// logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)


	// Initialize a new instance of our application struct, containing the // dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog, 
	}

	//functions.HTTP("webhook", app.webhook)
	//funcframework.RegisterHTTPFunction("webhook", app.webhook)
	funcframework.RegisterHTTPFunctionContext(context.TODO(), "webhook", app.webhook)
	port := ":8089"

	// export myFunc as env variable
	os.Setenv("FUNCTION_TARGET", "webhook")
	infoLog.Printf("Starting server on localhost%s", port)
}