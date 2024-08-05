package handlers

import (
	"encoding/json"
	"net/http"

	"api-gateway/internal/grpcclient"
	"api-gateway/internal/workermanager"
	"api-gateway/proto"
)

func ListPhonesHandler(workerManager *workermanager.WorkerManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Assuming the SIP ID is passed as a query parameter for listing phones
		sipID := r.URL.Query().Get("sipId")
		if sipID == "" {
			http.Error(w, "Missing sipId parameter", http.StatusBadRequest)
			return
		}

		workerAddress, err := workerManager.SelectWorker(sipID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		client, err := grpcclient.NewGRPCClient(workerAddress)
		if err != nil {
			http.Error(w, "Failed to create gRPC client", http.StatusInternalServerError)
			return
		}

		req := &proto.ListPhonesRequest{SipId: sipID}
		resp, err := client.ListPhones(r.Context(), req)
		if err != nil {
			http.Error(w, "Failed to list phones", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
