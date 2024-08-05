package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"api-gateway/internal/middleware"
	"api-gateway/internal/workermanager"
	"api-gateway/proto"
)

func RegisterPhoneHandler(workerManager workermanager.WorkerManagerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetRequestIDFromContext(r.Context())
		logger := middleware.GetLogger()

		logger.Debug(requestID, "RegisterPhoneHandler: Received request")
		var req proto.RegisterPhoneRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Debug(requestID, "RegisterPhoneHandler: Invalid request payload: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		logger.Debug(requestID, "RegisterPhoneHandler: Selecting worker for SIP ID %s", req.SipId)
		workerAddress, err := workerManager.SelectWorker(req.SipId)
		if err != nil {
			logger.Debug(requestID, "RegisterPhoneHandler: Error selecting worker: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Debug(requestID, "RegisterPhoneHandler: Creating gRPC client for worker at %s", workerAddress)
		client, err := workerManager.NewGRPCClient(workerAddress)
		if err != nil {
			logger.Debug(requestID, "RegisterPhoneHandler: Error creating gRPC client: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Debug(requestID, "RegisterPhoneHandler: Registering phone for SIP ID %s", req.SipId)
		resp, err := client.RegisterPhone(context.Background(), &req)
		if err != nil {
			logger.Debug(requestID, "RegisterPhoneHandler: Error registering phone: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Debug(requestID, "RegisterPhoneHandler: Successfully registered phone")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Debug(requestID, "RegisterPhoneHandler: Error encoding response: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
