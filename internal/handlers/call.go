package handlers

import (
	"encoding/json"
	"net/http"

	"api-gateway/internal/grpcclient"
	"api-gateway/internal/workermanager"
	"api-gateway/proto"
)

func CallHandler(workerManager *workermanager.WorkerManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req proto.CallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		workerAddress, err := workerManager.SelectWorker(req.SipId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		client, err := grpcclient.NewGRPCClient(workerAddress)
		if err != nil {
			http.Error(w, "Failed to create gRPC client", http.StatusInternalServerError)
			return
		}

		resp, err := client.Call(r.Context(), &req)
		if err != nil {
			http.Error(w, "Failed to initiate call", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
