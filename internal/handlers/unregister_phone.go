package handlers

import (
	"encoding/json"
	"net/http"

	"api-gateway/internal/grpcclient"
	"api-gateway/internal/workermanager"
	"api-gateway/proto"
)

func UnregisterPhoneHandler(workerManager *workermanager.WorkerManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req proto.UnregisterPhoneRequest
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

		resp, err := client.UnregisterPhone(r.Context(), &req)
		if err != nil {
			http.Error(w, "Failed to unregister phone", http.StatusInternalServerError)
			return
		}

		workerManager.Unbind(req.SipId)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
