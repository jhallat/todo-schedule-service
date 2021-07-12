package health

import (
	"encoding/json"
	"log"
	"net/http"
)

const healthPath = "/health"

func SetupHealth() {
	handleHealth := http.HandlerFunc(healthHandler)
	http.Handle(healthPath, handleHealth)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		checks := make([]Check, 0)
		health := Health {Status: "UP", Checks: checks}
		healthJson, err := json.Marshal(health)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(healthJson)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}