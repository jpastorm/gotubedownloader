package downloader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// QueueItemStatus represents the state of a download
type QueueItemStatus string

const (
	StatusQueued      QueueItemStatus = "queued"
	StatusDownloading QueueItemStatus = "downloading"
	StatusCompleted   QueueItemStatus = "completed"
	StatusFailed      QueueItemStatus = "failed"
	StatusCancelled   QueueItemStatus = "cancelled"
)

// QueueItem represents a single download in the queue
type QueueItem struct {
	ID          string          `json:"id"`
	URL         string          `json:"url"`
	Title       string          `json:"title"`
	Channel     string          `json:"channel"`
	Thumbnail   string          `json:"thumbnail"`
	Duration    string          `json:"duration"`
	Status      QueueItemStatus `json:"status"`
	Progress    float64         `json:"progress"`
	Speed       string          `json:"speed"`
	ETA         string          `json:"eta"`
	Downloaded  string          `json:"downloaded"`
	Total       string          `json:"total"`
	Error       string          `json:"error"`
	Mode        string          `json:"mode"`
	OutputDir   string          `json:"outputDir"`
	FilePath    string          `json:"filePath"`
	LogLines    []string        `json:"logLines"`
	CreatedAt   time.Time       `json:"createdAt"`
	CompletedAt *time.Time      `json:"completedAt"`
	Options     DownloadOptions `json:"options"`
	cancel      context.CancelFunc
}

// Queue manages concurrent downloads
type Queue struct {
	items       []*QueueItem
	mu          sync.RWMutex
	maxWorkers  int
	activeCount int
	onUpdate    func(item *QueueItem)
	onComplete  func(item *QueueItem)
	sem         chan struct{}
	ctx         context.Context
	cancelAll   context.CancelFunc
}

// NewQueue creates a new download queue
func NewQueue(maxWorkers int, onUpdate func(*QueueItem), onComplete func(*QueueItem)) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	q := &Queue{
		items:      []*QueueItem{},
		maxWorkers: maxWorkers,
		onUpdate:   onUpdate,
		onComplete: onComplete,
		sem:        make(chan struct{}, maxWorkers),
		ctx:        ctx,
		cancelAll:  cancel,
	}
	return q
}

// UpdateMaxWorkers changes the max concurrent downloads
func (q *Queue) UpdateMaxWorkers(n int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.maxWorkers = n
	q.sem = make(chan struct{}, n)
}

// Add adds one or more items to the queue
func (q *Queue) Add(opts DownloadOptions, info *VideoInfo) string {
	q.mu.Lock()

	item := &QueueItem{
		ID:        uuid.New().String(),
		URL:       opts.URL,
		Title:     info.Title,
		Channel:   info.Channel,
		Thumbnail: info.Thumbnail,
		Duration:  info.DurationStr,
		Status:    StatusQueued,
		Mode:      opts.Mode,
		OutputDir: opts.OutputDir,
		Options:   opts,
		CreatedAt: time.Now(),
		LogLines:  []string{},
	}

	q.items = append(q.items, item)
	q.mu.Unlock()

	// Start processing
	go q.processItem(item)

	return item.ID
}

// UpdateMeta updates an item's metadata (title, channel, thumbnail, duration) after background analysis
func (q *Queue) UpdateMeta(id string, info *VideoInfo) {
	q.mu.Lock()
	for _, item := range q.items {
		if item.ID == id {
			if info.Title != "" && item.Title == item.URL {
				item.Title = info.Title
			}
			if info.Channel != "" && item.Channel == "" {
				item.Channel = info.Channel
			}
			if info.Thumbnail != "" && item.Thumbnail == "" {
				item.Thumbnail = info.Thumbnail
			}
			if info.DurationStr != "" && item.Duration == "" {
				item.Duration = info.DurationStr
			}
			q.mu.Unlock()
			q.notifyUpdate(item)
			return
		}
	}
	q.mu.Unlock()
}

// AddBatch adds multiple URLs to the queue
func (q *Queue) AddBatch(urls []string, mode string, quality string, outputDir string) []string {
	ids := []string{}
	for _, url := range urls {
		opts := DownloadOptions{
			URL:       url,
			Mode:      mode,
			Quality:   quality,
			OutputDir: outputDir,
		}
		info := &VideoInfo{
			URL:   url,
			Title: fmt.Sprintf("Loading... %s", url),
		}
		id := q.Add(opts, info)
		ids = append(ids, id)
	}
	return ids
}

func (q *Queue) processItem(item *QueueItem) {
	// Wait for a slot
	select {
	case q.sem <- struct{}{}:
	case <-q.ctx.Done():
		return
	}
	defer func() { <-q.sem }()

	q.mu.Lock()
	if item.Status == StatusCancelled {
		q.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(q.ctx)
	item.cancel = cancel
	item.Status = StatusDownloading
	q.mu.Unlock()

	q.notifyUpdate(item)

	var lastFilePath string
	err := StartDownload(ctx, item.Options, func(p DownloadProgress) {
		q.mu.Lock()
		item.Progress = p.Percent
		item.Speed = p.Speed
		item.ETA = p.ETA
		item.Downloaded = p.Downloaded
		item.Total = p.Total
		if p.Title != "" {
			item.Title = p.Title
		}
		if p.FilePath != "" {
			lastFilePath = p.FilePath
		}
		if p.LogLine != "" {
			item.LogLines = append(item.LogLines, p.LogLine)
			// Keep max 100 log lines
			if len(item.LogLines) > 100 {
				item.LogLines = item.LogLines[len(item.LogLines)-100:]
			}
		}
		q.mu.Unlock()
		q.notifyUpdate(item)
	})

	q.mu.Lock()
	if err != nil {
		if item.Status != StatusCancelled {
			item.Status = StatusFailed
			item.Error = err.Error()
		}
	} else {
		item.Status = StatusCompleted
		item.Progress = 100
		now := time.Now()
		item.CompletedAt = &now
		if lastFilePath != "" {
			item.FilePath = lastFilePath
		}
		// Fallback: scan output directory for the newest matching file
		if item.FilePath == "" {
			item.FilePath = findDownloadedFile(item.OutputDir, item.Mode)
		}
	}
	q.mu.Unlock()

	q.notifyUpdate(item)
	if item.Status == StatusCompleted && q.onComplete != nil {
		q.onComplete(item)
	}
}

// findDownloadedFile scans dir for the most recently modified file matching mode (mp4/mp3)
func findDownloadedFile(dir string, mode string) string {
	ext := ".mp4"
	if mode == "audio" {
		ext = ".mp3"
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	var newest string
	var newestTime time.Time
	for _, e := range entries {
		if e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		name := strings.ToLower(e.Name())
		if !strings.HasSuffix(name, ext) {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.ModTime().After(newestTime) {
			newestTime = info.ModTime()
			newest = filepath.Join(dir, e.Name())
		}
	}
	return newest
}

func (q *Queue) notifyUpdate(item *QueueItem) {
	if q.onUpdate != nil {
		q.onUpdate(item)
	}
}

// Cancel cancels a specific download
func (q *Queue) Cancel(id string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, item := range q.items {
		if item.ID == id {
			item.Status = StatusCancelled
			if item.cancel != nil {
				item.cancel()
			}
			break
		}
	}
}

// Retry retries a failed download
func (q *Queue) Retry(id string) {
	q.mu.Lock()
	var item *QueueItem
	for _, i := range q.items {
		if i.ID == id {
			item = i
			break
		}
	}

	if item == nil || (item.Status != StatusFailed && item.Status != StatusCancelled) {
		q.mu.Unlock()
		return
	}

	item.Status = StatusQueued
	item.Progress = 0
	item.Error = ""
	item.Speed = ""
	item.ETA = ""
	item.LogLines = []string{}
	q.mu.Unlock()

	go q.processItem(item)
}

// Remove removes a download from the queue
func (q *Queue) Remove(id string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for i, item := range q.items {
		if item.ID == id {
			if item.cancel != nil {
				item.cancel()
			}
			q.items = append(q.items[:i], q.items[i+1:]...)
			break
		}
	}
}

// GetAll returns all queue items
func (q *Queue) GetAll() []*QueueItem {
	q.mu.RLock()
	defer q.mu.RUnlock()

	result := make([]*QueueItem, len(q.items))
	copy(result, q.items)
	return result
}

// ClearCompleted removes all completed/failed/cancelled items
func (q *Queue) ClearCompleted() {
	q.mu.Lock()
	defer q.mu.Unlock()

	filtered := []*QueueItem{}
	for _, item := range q.items {
		if item.Status == StatusQueued || item.Status == StatusDownloading {
			filtered = append(filtered, item)
		}
	}
	q.items = filtered
}

// CancelAll cancels all active downloads
func (q *Queue) CancelAll() {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, item := range q.items {
		if item.Status == StatusDownloading || item.Status == StatusQueued {
			item.Status = StatusCancelled
			if item.cancel != nil {
				item.cancel()
			}
		}
	}
}
