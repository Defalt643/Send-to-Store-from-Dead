package main

// Resource TODO
type Resource struct {
	Rabbit Rabbit
}

// Struct for search result from Elasticsearch
type SpeechRecognitionDetail struct {
	Link         string   `json:"link"`
	Text         []string `json:"text"`
	CreatedAt    int64    `json:"created_at"`
	UpdatedAt    int64    `json:"updated_at"`
	Status       string   `json:"status"`
	ErrorMessage string   `json:"error_message"`
}

type IncommingMessage struct {
	CreatedTime   int64    `json:"created_time"`
	Text          []string `json:"text"`
	ID            string   `json:"_id"`
	ChannelTypeID string   `json:"channel_type_id"`
}

type YoutubeDetal struct {
	MessageFromVideo []string `json:"message_from_video"`
}
