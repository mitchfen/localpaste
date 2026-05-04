package paste

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func cleanup() {
	os.Remove(pasteFile)
}

func TestWrite(t *testing.T) {
	defer cleanup()

	content := "test paste content"
	paste, err := Write(content)

	if err != nil {
		t.Fatalf("Write() failed: %v", err)
	}

	if paste.Content != content {
		t.Errorf("Expected content %q, got %q", content, paste.Content)
	}

	if paste.ExpiresAt.Before(time.Now()) {
		t.Error("Paste expiration time should be in the future")
	}
}

func TestRead(t *testing.T) {
	defer cleanup()

	content := "test content"
	_, err := Write(content)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	paste, err := Read()
	if err != nil {
		t.Fatalf("Read() failed: %v", err)
	}

	if paste == nil {
		t.Fatal("Expected paste, got nil")
	}

	if paste.Content != content {
		t.Errorf("Expected content %q, got %q", content, paste.Content)
	}
}

func TestReadNonExistent(t *testing.T) {
	defer cleanup()

	paste, err := Read()
	if err != nil {
		t.Fatalf("Read() should not error on missing file: %v", err)
	}

	if paste != nil {
		t.Fatal("Expected nil for non-existent paste")
	}
}

func TestExpiresIn(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		shouldBe func(string) bool
	}{
		{
			name:     "expired",
			duration: -1 * time.Minute,
			shouldBe: func(s string) bool { return s == "expired" },
		},
		{
			name:     "less than a minute",
			duration: 30 * time.Second,
			shouldBe: func(s string) bool { return s == "< 1m" },
		},
		{
			name:     "5 minutes",
			duration: 5 * time.Minute,
			shouldBe: func(s string) bool { return s == "4m" || s == "5m" }, // Account for test execution time
		},
		{
			name:     "60 minutes",
			duration: 60 * time.Minute,
			shouldBe: func(s string) bool { return s == "59m" || s == "1h" }, // Account for test execution time
		},
		{
			name:     "2 hours",
			duration: 2 * time.Hour,
			shouldBe: func(s string) bool { return s == "1h" || s == "2h" }, // Account for test execution time
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paste := &Paste{
				Content:   "test",
				ExpiresAt: time.Now().Add(tt.duration),
			}

			result := paste.ExpiresIn()
			if !tt.shouldBe(result) {
				t.Errorf("Got unexpected value: %q", result)
			}
		})
	}
}

func TestWriteAndReadPersistence(t *testing.T) {
	defer cleanup()

	content := "multi-line\ncontent\ntest"
	p1, err := Write(content)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	p2, err := Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if p2 == nil {
		t.Fatal("Expected paste after read")
	}

	if p1.Content != p2.Content {
		t.Errorf("Content mismatch: %q != %q", p1.Content, p2.Content)
	}

	// Check expiration is approximately the same (within 1 second for test execution)
	diff := p1.ExpiresAt.Sub(p2.ExpiresAt)
	if diff < 0 {
		diff = -diff
	}
	if diff > time.Second {
		t.Errorf("Expiration time drift too large: %v", diff)
	}
}

func TestExpiredPasteNotRead(t *testing.T) {
	defer cleanup()

	// Create a paste and manually set it to expired
	paste := &Paste{
		Content:   "expired content",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	data, _ := json.Marshal(paste)
	os.WriteFile(pasteFile, data, 0644)

	// Try to read - should return nil and delete the file
	result, err := Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if result != nil {
		t.Fatal("Expected nil for expired paste")
	}

	// File should be deleted
	if _, err := os.Stat(pasteFile); !os.IsNotExist(err) {
		t.Error("Expired paste file should be deleted")
	}
}
