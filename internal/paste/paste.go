package paste

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

const pasteFile = "paste.json"

var fileMutex sync.Mutex

type Paste struct {
	Content   string    `json:"content"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (p *Paste) ExpiresIn() string {
	remaining := time.Until(p.ExpiresAt)
	if remaining <= 0 {
		return "expired"
	}
	minutes := int(remaining.Minutes())
	if minutes < 1 {
		return "< 1m"
	}
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}
	hours := int(remaining.Hours())
	return fmt.Sprintf("%dh", hours)
}

func Read() (*Paste, error) {
	data, err := os.ReadFile(pasteFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var paste Paste
	if err := json.Unmarshal(data, &paste); err != nil {
		return nil, err
	}

	// Check if expired
	if time.Now().After(paste.ExpiresAt) {
		os.Remove(pasteFile)
		return nil, nil
	}

	return &paste, nil
}

func Write(content string) (*Paste, error) {
	paste := &Paste{
		Content:   content,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	data, err := json.Marshal(paste)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(pasteFile, data, 0644); err != nil {
		return nil, err
	}

	return paste, nil
}

func GetMutex() *sync.Mutex {
	return &fileMutex
}

func StartCleanupRoutine() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			fileMutex.Lock()
			paste, _ := Read()
			if paste != nil && time.Now().After(paste.ExpiresAt) {
				os.Remove(pasteFile)
			}
			fileMutex.Unlock()
		}
	}()
}
