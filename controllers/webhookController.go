package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"qiscus-test/services"
)

type AgentController struct {
	webhookService services.AllocationService
}

func WebHookController(s *services.AllocationService) *AgentController {
	return &AgentController{webhookService: *s}
}

func (c *AgentController) Webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	prettyJSON, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
	} else {
		fmt.Println("Webhook received JSON:\n", string(prettyJSON))
	}

	roomID := payload["room_id"].(string)
	if roomID == "" {
		http.Error(w, "room_id is required", http.StatusBadRequest)
		return
	}

	err = c.webhookService.EnqueueChat(roomID)
	if err != nil {
		http.Error(w, "Failed to enqueue chat", http.StatusInternalServerError)
		return
	}
}
