package handler

import (
	"encoding/json"
	"net/http"
)

func handleRoot() http.HandlerFunc {
	body := map[string]any{
		"service":     "Comemora",
		"description": "Serviço de lembretes de aniversários e datas especiais",
		"endpoints": []map[string]string{
			{"method": "GET", "path": "/events", "description": "Lista todos os eventos"},
			{"method": "POST", "path": "/events", "description": "Cria um novo evento"},
			{"method": "GET", "path": "/events/export", "description": "Exporta eventos em CSV"},
			{"method": "POST", "path": "/events/import", "description": "Importa eventos via CSV"},
			{"method": "GET", "path": "/trigger-check", "description": "Dispara verificação de aniversários do dia"},
			{"method": "GET", "path": "/health", "description": "Health check"},
		},
	}
	payload, _ := json.Marshal(body)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}
