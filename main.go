package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

// Define structs for Slack and Google Chat messages
type Message struct {
	Text string `json:"text"`
	// Add other fields as needed
}

// List of allowed space IDs
var allowedSpaceIDs map[string]struct{}

// Initialize the list of allowed space IDs from environment variable
func init() {
	ids := os.Getenv("ALLOWED_SPACE_IDS")
	if ids != "" {
		idList := strings.Split(ids, ",")
		allowedSpaceIDs = make(map[string]struct{}, len(idList))
		for _, id := range idList {
			allowedSpaceIDs[id] = struct{}{}
		}
	}
}

// Handle Slack webhook requests
func slackHandler(w http.ResponseWriter, r *http.Request) {
	var msg Message

	switch r.Method {
	case http.MethodPost:
		if r.Header.Get("Content-Type") == "application/json" {
			err := json.NewDecoder(r.Body).Decode(&msg)
			if err != nil {
				http.Error(w, "Error decoding JSON request", http.StatusBadRequest)
				log.Printf("Error decoding JSON request: %s", err)
				return
			}
		} else {
			r.ParseForm()
			payload := r.FormValue("payload")
			err := json.Unmarshal([]byte(payload), &msg)
			if err != nil {
				http.Error(w, "Error decoding Slack message", http.StatusBadRequest)
				log.Printf("Error decoding Slack message: %s", err)
				return
			}
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// If allowedSpaceIDs is nil or empty, skip space ID check
	if len(allowedSpaceIDs) > 0 {
		// Extract space ID from URL path
		parts := strings.Split(r.URL.Path, "/")
		spaceID := parts[len(parts)-2]

		// Check if the space ID is allowed
		if _, ok := allowedSpaceIDs[spaceID]; !ok {
			http.Error(w, "Space ID not allowed", http.StatusForbidden)
			log.Printf("Space ID %s not allowed", spaceID)
			return
		}
	}

	// Translate Slack message to Google Chat format
	googleChatMsg := Message{
		Text: msg.Text,
	}

	// Convert Google Chat message to JSON
	buf, err := json.Marshal(googleChatMsg)
	if err != nil {
		http.Error(w, "Error encoding Google Chat message", http.StatusInternalServerError)
		log.Printf("Error encoding Google Chat message: %s", err)
		return
	}

	// Build Google Chat webhook URL with path and query parameters from incoming request
	googleChatWebhookURL := "https://chat.googleapis.com" + r.URL.Path
	if r.URL.RawQuery != "" {
		googleChatWebhookURL += "?" + r.URL.RawQuery
	}

	// Send translated message to Google Chat webhook
	_, err = http.Post(googleChatWebhookURL, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		http.Error(w, "Error sending message to Google Chat webhook", http.StatusInternalServerError)
		log.Printf("Error sending message to Google Chat webhook: %s", err)
		return
	}

	// Respond to Slack with success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message translated and sent to Google Chat"))
}

// HealthCheck handler for checking the status
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	// Set up HTTP handlers
	http.HandleFunc("/v1/spaces/*", slackHandler)
	http.HandleFunc("/healthcheck", healthCheckHandler)

	// Get the port from environment variable, default to 8080 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start HTTP server
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
