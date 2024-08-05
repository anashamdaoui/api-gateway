package routers

import (
	"api-gateway/internal/handlers"
	"api-gateway/internal/workermanager"

	"github.com/gorilla/mux"
)

// InitRoutes initializes the HTTP routes for the API gateway.
func InitRoutes(router *mux.Router, registryAddress string) {
	// Create a new WorkerManager instance
	workerManager := workermanager.NewWorkerManager(registryAddress)

	// Set up the route handlers
	router.HandleFunc("/phones/register", handlers.RegisterPhoneHandler(workerManager)).Methods("POST")
	router.HandleFunc("/phones/unregister", handlers.UnregisterPhoneHandler(workerManager)).Methods("POST")
	router.HandleFunc("/phones/call", handlers.CallHandler(workerManager)).Methods("POST")
	router.HandleFunc("/phones/answer/{callID}", handlers.AnswerCallHandler(workerManager)).Methods("POST")
	router.HandleFunc("/phones/hangup/{callID}", handlers.HangupCallHandler(workerManager)).Methods("POST")
	router.HandleFunc("/phones/hold/{callID}", handlers.HoldCallHandler(workerManager)).Methods("POST")
	router.HandleFunc("/phones/resume/{callID}", handlers.ResumeCallHandler(workerManager)).Methods("POST")
	router.HandleFunc("/phones/list", handlers.ListPhonesHandler(workerManager)).Methods("GET")
}
