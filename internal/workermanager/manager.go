package workermanager

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"api-gateway/internal/grpcclient"
)

// WorkerManagerInterface defines the methods for worker management.
type WorkerManagerInterface interface {
	SelectWorker(sipID string) (string, error)
	NewGRPCClient(workerAddress string) (*grpcclient.GRPCClient, error)
	Unbind(sipID string)
}

// WorkerManager manages worker selection and tracking.
type WorkerManager struct {
	RegistryAddress   string
	SipIDToWorkerAddr map[string]string
	WorkerUsage       map[string]int
	mu                sync.Mutex
}

// NewWorkerManager creates a new WorkerManager instance.
func NewWorkerManager(registryAddress string) *WorkerManager {
	return &WorkerManager{
		RegistryAddress:   registryAddress,
		SipIDToWorkerAddr: make(map[string]string),
		WorkerUsage:       make(map[string]int),
	}
}

// GetHealthyWorkers retrieves a list of healthy workers from the registry.
func (wm *WorkerManager) GetHealthyWorkers() ([]string, error) {
	log.Println("WorkerManager: Fetching healthy workers from registry")
	resp, err := http.Get(wm.RegistryAddress + "/workers/healthy")
	if err != nil {
		log.Printf("WorkerManager: Error fetching healthy workers: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var workers []string
	if err := json.NewDecoder(resp.Body).Decode(&workers); err != nil {
		log.Printf("WorkerManager: Error decoding workers: %v", err)
		return nil, err
	}

	log.Printf("WorkerManager: Retrieved %d healthy workers", len(workers))
	return workers, nil
}

// SelectWorker selects a worker for the given SIP ID.
func (wm *WorkerManager) SelectWorker(sipID string) (string, error) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	log.Printf("WorkerManager: Selecting worker for SIP ID %s", sipID)
	workerAddress, exists := wm.SipIDToWorkerAddr[sipID]
	if exists {
		log.Printf("WorkerManager: Found existing worker for SIP ID %s: %s", sipID, workerAddress)
		// Check if the existing worker is healthy
		if wm.isWorkerHealthy(workerAddress) {
			return workerAddress, nil
		} else {
			log.Printf("WorkerManager: Worker %s is not healthy", workerAddress)
			return "", errors.New("phone not ready. Please try again later or unregister and register again")
		}
	}

	workers, err := wm.GetHealthyWorkers()
	if err != nil || len(workers) == 0 {
		log.Println("WorkerManager: No healthy workers available")
		return "", errors.New("no healthy workers available")
	}

	// Select a new worker if there are new healthy workers
	leastUsedWorker := workers[0]
	minUsage := wm.WorkerUsage[leastUsedWorker]

	for _, worker := range workers {
		usage, ok := wm.WorkerUsage[worker]
		if !ok {
			usage = 0 // Default to 0 if the worker is not in WorkerUsage
		}
		if usage < minUsage {
			leastUsedWorker = worker
			minUsage = usage
		}
	}

	workerAddress = leastUsedWorker
	wm.SipIDToWorkerAddr[sipID] = workerAddress
	wm.WorkerUsage[workerAddress]++

	log.Printf("WorkerManager: Selected worker %s for SIP ID %s", workerAddress, sipID)
	return workerAddress, nil
}

// isWorkerHealthy checks if a worker is healthy.
func (wm *WorkerManager) isWorkerHealthy(address string) bool {
	log.Printf("WorkerManager: Checking health of worker at %s", address)
	resp, err := http.Get(address + "/healthcheck")
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("WorkerManager: Worker %s is not healthy", address)
		return false
	}
	return true
}

// NewGRPCClient creates a new gRPC client for the given worker address.
func (wm *WorkerManager) NewGRPCClient(workerAddress string) (*grpcclient.GRPCClient, error) {
	log.Printf("WorkerManager: Creating new gRPC client for worker at %s", workerAddress)
	return grpcclient.NewGRPCClient(workerAddress)
}

// Unbind removes the mapping of a SIP ID to a worker address.
func (wm *WorkerManager) Unbind(sipID string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	log.Printf("WorkerManager: Unbinding SIP ID %s", sipID)

	if workerAddr, exists := wm.SipIDToWorkerAddr[sipID]; exists {
		delete(wm.SipIDToWorkerAddr, sipID)
		wm.WorkerUsage[workerAddr]--
		if wm.WorkerUsage[workerAddr] == 0 {
			delete(wm.WorkerUsage, workerAddr)
		}
		log.Printf("WorkerManager: Worker %s usage decremented", workerAddr)
	}
}
