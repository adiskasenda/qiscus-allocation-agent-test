package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AvailableAgentsResponse struct {
	Data struct {
		Agents []struct {
			Id                   int    `json:"id"`
			Email                string `json:"email"`
			IsAvailable          bool   `json:"is_available"`
			CurrentCustomerCount int    `json:"current_customer_count"`
		} `json:"agents"`
	} `json:"data"`
}

type APIService interface {
	// GetAvailableAgents(room_id string) ([]byte, error)
	GetAvailableAgents(room_id string) (*AvailableAgentsResponse, error)
	AssignAgent(roomID string, agentID int) error
}

type apiService struct {
	BaseUrl   string
	SecretKey string
	AppCode   string
}

func NewAPIService(baseUrl, secretKey, appCode string) APIService {
	return &apiService{
		BaseUrl:   baseUrl,
		SecretKey: secretKey,
		AppCode:   appCode,
	}
}

func (s *apiService) GetAvailableAgents(room_id string) (*AvailableAgentsResponse, error) {
	endpoint := fmt.Sprintf("%s/api/v2/admin/service/available_agents?room_id=%s", s.BaseUrl, room_id)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Add("Qiscus-Secret-Key", s.SecretKey)
	req.Header.Add("Qiscus-App-Id", s.AppCode)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var parsed AvailableAgentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	return &parsed, nil
	// return io.ReadAll(resp.Body)
}

func (s *apiService) AssignAgent(roomID string, agentID int) error {
	endpoint := fmt.Sprintf("%s/api/v1/admin/service/assign_agent", s.BaseUrl)

	// Buat payload JSON
	payload := map[string]interface{}{
		"room_id":  roomID,
		"agent_id": agentID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Buat request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Qiscus-Secret-Key", s.SecretKey)
	req.Header.Set("Qiscus-App-Id", s.AppCode)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Jika status bukan 2xx, return error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to assign agent: status %d", resp.StatusCode)
	}

	return nil
}
