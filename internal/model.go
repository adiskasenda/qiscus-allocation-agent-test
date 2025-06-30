package internal

type Agent struct {
	Email     string
	Available bool
	ChatCount int
}

type ChatWebhook struct {
	RoomID string `json:"room_id"`
}
