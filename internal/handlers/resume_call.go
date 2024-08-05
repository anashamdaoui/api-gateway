package handlers

import (
	"encoding/json"
	"net/http"

	"api-gateway/internal/grpcclient"
	"api-gateway/internal/workermanager"
	"api-gateway/proto"

	"github.com/gorilla/mux"
)

func ResumeCallHandler(workerManager *workermanager.WorkerManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		callID := vars["callID"]

		workerAddress, err := workerManager.SelectWorker(callID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		client, err := grpcclient.NewGRPCClient(workerAddress)
		if err != nil {
			http.Error(w, "Failed to create gRPC client", http.StatusInternalServerError)
			return
		}

		req := &proto.CallActionRequest{CallId: callID}
		resp, err := client.ResumeCall(r.Context(), req)
		if err != nil {
			http.Error(w, "Failed to resume call", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
