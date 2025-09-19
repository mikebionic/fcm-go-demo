package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "path/filepath"
    
    firebase "firebase.google.com/go/v4"
    "firebase.google.com/go/v4/messaging"
    "google.golang.org/api/option"
)

// NotificationRequest represents the incoming request structure
type NotificationRequest struct {
    Token string `json:"token"`
    Title string `json:"title"`
    Body  string `json:"body"`
    Data  map[string]string `json:"data,omitempty"`
}

// NotificationResponse represents the response structure
type NotificationResponse struct {
    Success   bool   `json:"success"`
    MessageID string `json:"messageId,omitempty"`
    Error     string `json:"error,omitempty"`
}

type FCMServer struct {
    messagingClient *messaging.Client
}

func NewFCMServer() (*FCMServer, error) {
    ctx := context.Background()
    opt := option.WithCredentialsFile("secrets/serviceAccountKey.json")
    app, err := firebase.NewApp(ctx, nil, opt)
    if err != nil {
        return nil, fmt.Errorf("error initializing Firebase app: %v", err)
    }

    client, err := app.Messaging(ctx)
    if err != nil {
        return nil, fmt.Errorf("error getting Messaging client: %v", err)
    }

    return &FCMServer{
        messagingClient: client,
    }, nil
}

func (s *FCMServer) sendNotificationHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Content-Type", "application/json")
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    if r.Method != "POST" {
        response := NotificationResponse{
            Success: false,
            Error:   "Method not allowed. Use POST",
        }
        w.WriteHeader(http.StatusMethodNotAllowed)
        json.NewEncoder(w).Encode(response)
        return
    }

    var req NotificationRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response := NotificationResponse{
            Success: false,
            Error:   "Invalid JSON format",
        }
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(response)
        return
    }

    // Validate required fields
    if req.Token == "" {
        response := NotificationResponse{
            Success: false,
            Error:   "Device token is required",
        }
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(response)
        return
    }

    if req.Title == "" || req.Body == "" {
        response := NotificationResponse{
            Success: false,
            Error:   "Title and body are required",
        }
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(response)
        return
    }

    // Create the FCM message
    message := &messaging.Message{
        Notification: &messaging.Notification{
            Title: req.Title,
            Body:  req.Body,
        },
        Token: req.Token,
    }

    // Add custom data if provided
    if req.Data != nil && len(req.Data) > 0 {
        message.Data = req.Data
    }

    // Send the message
    ctx := context.Background()
    messageID, err := s.messagingClient.Send(ctx, message)
    if err != nil {
        log.Printf("Error sending FCM message: %v", err)
        response := NotificationResponse{
            Success: false,
            Error:   fmt.Sprintf("Failed to send notification: %v", err),
        }
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(response)
        return
    }

    log.Printf("Successfully sent FCM message. ID: %s", messageID)
    response := NotificationResponse{
        Success:   true,
        MessageID: messageID,
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func staticFileHandler(w http.ResponseWriter, r *http.Request) {
    if filepath.IsAbs(r.URL.Path) || filepath.Clean(r.URL.Path) != r.URL.Path {
        http.Error(w, "Invalid path", http.StatusBadRequest)
        return
    }

    path := r.URL.Path
    if path == "/" {
        path = "/index.html"
    }

    filePath := filepath.Join("web", path)
    http.ServeFile(w, r, filePath)
}

func main() {
    server, err := NewFCMServer()
    if err != nil {
        log.Fatalf("Failed to initialize FCM server: %v", err)
    }

    http.HandleFunc("/send-notification", server.sendNotificationHandler)

    fs := http.FileServer(http.Dir("./web"))
    http.Handle("/", fs)

    port := ":8000"
    log.Printf("🚀 FCM Server starting on port %s", port)
    log.Printf("📱 View notifications: http://localhost:8000/")
    
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}