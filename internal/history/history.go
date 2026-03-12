package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"GotubeDownloader/internal/config"
)

type HistoryEntry struct {
	ID           string    `json:"id"`
	URL          string    `json:"url"`
	Title        string    `json:"title"`
	Channel      string    `json:"channel"`
	Duration     string    `json:"duration"`
	Thumbnail    string    `json:"thumbnail"`
	FilePath     string    `json:"filePath"`
	Mode         string    `json:"mode"`
	Quality      string    `json:"quality"`
	Status       string    `json:"status"`
	FileSize     int64     `json:"fileSize"`
	DownloadedAt time.Time `json:"downloadedAt"`
}

type History struct {
	Entries []HistoryEntry `json:"entries"`
	mu      sync.RWMutex
	path    string
}

var (
	instance *History
	once     sync.Once
)

func Load() (*History, error) {
	var loadErr error
	once.Do(func() {
		dir, err := config.GetConfigDir()
		if err != nil {
			loadErr = err
			return
		}

		histPath := filepath.Join(dir, "history.json")
		instance = &History{
			Entries: []HistoryEntry{},
			path:    histPath,
		}

		data, err := os.ReadFile(histPath)
		if err != nil {
			if os.IsNotExist(err) {
				return
			}
			loadErr = err
			return
		}

		if err := json.Unmarshal(data, &instance.Entries); err != nil {
			loadErr = err
			return
		}
	})
	return instance, loadErr
}

func (h *History) Add(entry HistoryEntry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Prepend (newest first)
	h.Entries = append([]HistoryEntry{entry}, h.Entries...)

	// Keep max 500 entries
	if len(h.Entries) > 500 {
		h.Entries = h.Entries[:500]
	}

	return h.save()
}

func (h *History) GetAll() []HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.Entries
}

func (h *History) Clear() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Entries = []HistoryEntry{}
	return h.save()
}

func (h *History) Remove(id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i, e := range h.Entries {
		if e.ID == id {
			h.Entries = append(h.Entries[:i], h.Entries[i+1:]...)
			break
		}
	}

	return h.save()
}

func (h *History) save() error {
	data, err := json.MarshalIndent(h.Entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0644)
}
