package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"api-gateway/internal/workermanager"
	"api-gateway/proto"

	"github.com/gorilla/mux"
)

func AnswerCallHandler(workerManager workermanager.WorkerManagerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		callID := vars["callID"]

		client, err := workerManager.NewGRPCClient("dummy-worker-address") // Provide the actual worker address
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req := &proto.CallActionRequest{
			CallId: callID,
		}

		resp, err := client.AnswerCall(context.Background(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
