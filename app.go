package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"GotubeDownloader/internal/bootstrap"
	"GotubeDownloader/internal/config"
	"GotubeDownloader/internal/downloader"
	"GotubeDownloader/internal/history"
)

// App struct – Wails binding
type App struct {
	ctx     context.Context
	queue   *downloader.Queue
	history *history.History
	config  *config.Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Warning: failed to load config:", err)
		cfg = config.DefaultConfig()
	}
	a.config = cfg

	// Ensure download directory exists
	os.MkdirAll(cfg.DownloadDir, 0755)

	// Load history
	h, err := history.Load()
	if err != nil {
		fmt.Println("Warning: failed to load history:", err)
	}
	a.history = h

	// Initialize queue
	a.queue = downloader.NewQueue(cfg.MaxConcurrent, a.onQueueUpdate, a.onDownloadComplete)

	// Start shell (tray icon + global hotkey)
	a.initShell()
}

// shutdown is called when the app exits
func (a *App) shutdown(ctx context.Context) {
	a.cleanupShell()
}

// ===== BOOTSTRAP =====

func (a *App) CheckBootstrap() bootstrap.BootstrapStatus {
	return bootstrap.CheckStatus()
}

func (a *App) InstallYtDlp() bootstrap.BootstrapStatus {
	path, err := bootstrap.DownloadYtDlp()
	if err != nil {
		return bootstrap.BootstrapStatus{
			YtDlpReady: false,
			Error:      err.Error(),
		}
	}
	status := bootstrap.CheckStatus()
	status.YtDlpPath = path
	status.YtDlpReady = true
	return status
}

// ===== WINDOW =====

func (a *App) HideWindow() {
	a.HideOverlay()
}

func (a *App) ShowWindow() {
	a.ShowOverlay()
}

// ===== ANALYSIS =====

func (a *App) AnalyzeURL(url string) (*downloader.VideoInfo, error) {
	if url == "" {
		return nil, fmt.Errorf("URL is empty")
	}
	return downloader.Analyze(url)
}

func (a *App) AnalyzePlaylistURL(url string) (*downloader.VideoInfo, error) {
	if url == "" {
		return nil, fmt.Errorf("URL is empty")
	}
	return downloader.AnalyzePlaylist(url)
}

// ===== DOWNLOADS =====

func (a *App) StartDownload(url string, mode string, quality string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("URL is empty")
	}

	outputDir := a.config.DownloadDir
	if mode == "" {
		mode = a.config.DefaultMode
	}
	if quality == "" {
		quality = "high"
	}

	opts := downloader.DownloadOptions{
		URL:       url,
		Mode:      mode,
		Quality:   quality,
		OutputDir: outputDir,
	}

	// Queue immediately with minimal info (instant feedback)
	info := &downloader.VideoInfo{
		URL:   url,
		Title: url,
		Type:  "video",
	}

	id := a.queue.Add(opts, info)

	// Fetch metadata in background then update item
	go func() {
		meta, err := downloader.Analyze(url)
		if err == nil && meta != nil {
			a.queue.UpdateMeta(id, meta)
		}
	}()

	return id, nil
}

func (a *App) StartBatchDownload(urls []string, mode string, quality string) ([]string, error) {
	if len(urls) == 0 {
		return nil, fmt.Errorf("no URLs provided")
	}
	if quality == "" {
		quality = "high"
	}
	outputDir := a.config.DownloadDir
	ids := a.queue.AddBatch(urls, mode, quality, outputDir)
	return ids, nil
}

func (a *App) CancelDownload(id string) {
	a.queue.Cancel(id)
}

func (a *App) RetryDownload(id string) {
	a.queue.Retry(id)
}

func (a *App) RemoveDownload(id string) {
	a.queue.Remove(id)
}

func (a *App) ClearCompleted() {
	a.queue.ClearCompleted()
}

func (a *App) CancelAllDownloads() {
	a.queue.CancelAll()
}

func (a *App) GetQueue() []*downloader.QueueItem {
	return a.queue.GetAll()
}

// ===== HISTORY =====

func (a *App) GetHistory() []history.HistoryEntry {
	if a.history == nil {
		return []history.HistoryEntry{}
	}
	return a.history.GetAll()
}

func (a *App) ClearHistory() error {
	if a.history == nil {
		return nil
	}
	return a.history.Clear()
}

func (a *App) RemoveHistoryEntry(id string) error {
	if a.history == nil {
		return nil
	}
	return a.history.Remove(id)
}

// ===== SETTINGS =====

func (a *App) GetSettings() *config.Config {
	return a.config
}

func (a *App) SaveSettings(cfg config.Config) error {
	if err := config.Save(&cfg); err != nil {
		return err
	}
	a.config = &cfg

	// Update queue concurrency
	a.queue.UpdateMaxWorkers(cfg.MaxConcurrent)

	// Ensure download dir exists
	os.MkdirAll(cfg.DownloadDir, 0755)

	return nil
}

// ===== FILE OPS =====

func (a *App) SelectFolder() (string, error) {
	dir, err := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select Download Folder",
	})
	if err != nil {
		return "", err
	}
	return dir, nil
}

func (a *App) OpenFolder(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", filepath.Dir(path))
	case "darwin":
		cmd = exec.Command("open", filepath.Dir(path))
	default:
		cmd = exec.Command("xdg-open", filepath.Dir(path))
	}
	return cmd.Start()
}

func (a *App) OpenFile(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}

// OpenURL opens a URL in the default browser.
func (a *App) OpenURL(url string) {
	wailsRuntime.BrowserOpenURL(a.ctx, url)
}

func (a *App) GetDiagnostics() string {
	status := bootstrap.CheckStatus()
	cfg := a.config

	diag := fmt.Sprintf(`GotubeDownloader Diagnostics
================================
OS: %s/%s
Time: %s

yt-dlp: %v (%s)
ffmpeg: %v (%s)

Config:
  Download Dir: %s
  Max Concurrent: %d
  Fragment Concurrent: %d
  Default Mode: %s
  Proxy: %s
  Speed Limit: %s

Queue: %d items
`,
		runtime.GOOS, runtime.GOARCH,
		time.Now().Format(time.RFC3339),
		status.YtDlpReady, status.YtDlpPath,
		status.FfmpegReady, status.FfmpegPath,
		cfg.DownloadDir,
		cfg.MaxConcurrent,
		cfg.FragmentConcurrent,
		cfg.DefaultMode,
		cfg.Proxy,
		cfg.SpeedLimit,
		len(a.queue.GetAll()),
	)

	return diag
}

// ===== CALLBACKS =====

func (a *App) onQueueUpdate(item *downloader.QueueItem) {
	wailsRuntime.EventsEmit(a.ctx, "queue:update", item)
}

func (a *App) onDownloadComplete(item *downloader.QueueItem) {
	// Add to history
	if a.history != nil {
		entry := history.HistoryEntry{
			ID:           item.ID,
			URL:          item.URL,
			Title:        item.Title,
			Channel:      item.Channel,
			Duration:     item.Duration,
			Thumbnail:    item.Thumbnail,
			FilePath:     item.FilePath,
			Mode:         item.Mode,
			Status:       "completed",
			DownloadedAt: time.Now(),
		}
		a.history.Add(entry)
	}

	// Emit completion event
	wailsRuntime.EventsEmit(a.ctx, "download:complete", item)

	// Native macOS notification
	if a.config.NotifyOnComplete {
		title := item.Title
		if title == "" || title == item.URL {
			title = "Download"
		}
		subtitle := "Download Complete"
		if item.Mode == "audio" {
			subtitle = "Audio Ready"
		}
		// Sanitize for AppleScript: escape backslashes then double-quotes
		safeTitle := strings.ReplaceAll(title, `\`, `\\`)
		safeTitle = strings.ReplaceAll(safeTitle, `"`, `\"`)
		script := fmt.Sprintf(
			`display notification "%s" with title "GoTube" subtitle "%s"`,
			safeTitle, subtitle,
		)
		go exec.Command("osascript", "-e", script).Run()
	}
}
