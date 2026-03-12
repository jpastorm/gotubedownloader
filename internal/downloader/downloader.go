package downloader

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"GotubeDownloader/internal/bootstrap"
	"GotubeDownloader/internal/config"
)

// VideoInfo holds analyzed video/playlist metadata
type VideoInfo struct {
	ID          string      `json:"id"`
	URL         string      `json:"url"`
	Title       string      `json:"title"`
	Channel     string      `json:"channel"`
	Duration    float64     `json:"duration"`
	DurationStr string      `json:"durationStr"`
	Thumbnail   string      `json:"thumbnail"`
	Type        string      `json:"type"`
	VideoCount  int         `json:"videoCount"`
	Entries     []VideoInfo `json:"entries,omitempty"`
	IsLive      bool        `json:"isLive"`
	Description string      `json:"description"`
}

// DownloadProgress holds real-time download progress
type DownloadProgress struct {
	ID         string  `json:"id"`
	Status     string  `json:"status"`
	Percent    float64 `json:"percent"`
	Speed      string  `json:"speed"`
	ETA        string  `json:"eta"`
	Downloaded string  `json:"downloaded"`
	Total      string  `json:"total"`
	Title      string  `json:"title"`
	LogLine    string  `json:"logLine"`
	FilePath   string  `json:"filePath"`
}

// DownloadOptions configures a download
type DownloadOptions struct {
	URL           string `json:"url"`
	Mode          string `json:"mode"`
	Quality       string `json:"quality"`
	OutputDir     string `json:"outputDir"`
	Subtitles     bool   `json:"subtitles"`
	SubtitleLang  string `json:"subtitleLang"`
	EmbedSubs     bool   `json:"embedSubs"`
	EmbedMetadata bool   `json:"embedMetadata"`
}

// ytDlpJSON represents the raw yt-dlp JSON output
type ytDlpJSON struct {
	ID            string      `json:"id"`
	Title         string      `json:"title"`
	Uploader      string      `json:"uploader"`
	Channel       string      `json:"channel"`
	Duration      float64     `json:"duration"`
	Thumbnail     string      `json:"thumbnail"`
	Type          string      `json:"_type"`
	Entries       []ytDlpJSON `json:"entries"`
	WebpageURL    string      `json:"webpage_url"`
	Description   string      `json:"description"`
	IsLive        bool        `json:"is_live"`
	PlaylistCount int         `json:"playlist_count"`
}

var progressRegex = regexp.MustCompile(`\[download\]\s+(\d+\.?\d*)%\s+of\s+~?\s*(\S+)\s+at\s+(\S+)\s+ETA\s+(\S+)`)
var mergeRegex = regexp.MustCompile(`\[Merger\]|Merging formats`)
var destinationRegex = regexp.MustCompile(`\[(?:download|ExtractAudio|Merger)\]\s+(?:Destination:\s*|Merging formats into ")([^"\r\n]+)`)

// Analyze fetches metadata for a given URL — ultra fast
// Uses --no-playlist to skip playlist extraction for single video URLs
// For explicit playlist analysis, use AnalyzePlaylist
func Analyze(url string) (*VideoInfo, error) {
	ytDlpPath, err := bootstrap.FindYtDlp()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp not found: %w", err)
	}

	cfg := config.Get()

	// --dump-json --no-playlist: fastest possible, single video only
	// --no-download --no-warnings: skip downloading, skip warnings
	// --socket-timeout 10: fail fast on slow connections
	args := []string{
		"--dump-json",
		"--no-download",
		"--no-warnings",
		"--no-playlist",
		"--no-check-certificates",
		"--socket-timeout", "10",
	}

	if cfg.Proxy != "" {
		args = append(args, "--proxy", cfg.Proxy)
	}

	args = append(args, url)

	ctx, cancel := context.WithTimeout(context.Background(), 15*1000000000) // 15 seconds max
	defer cancel()

	cmd := exec.CommandContext(ctx, ytDlpPath, args...)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("analysis failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	var raw ytDlpJSON
	if err := json.Unmarshal(output, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return convertToVideoInfo(raw, url), nil
}

// AnalyzePlaylist fetches metadata for a playlist URL
func AnalyzePlaylist(url string) (*VideoInfo, error) {
	ytDlpPath, err := bootstrap.FindYtDlp()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp not found: %w", err)
	}

	cfg := config.Get()

	args := []string{
		"--dump-single-json",
		"--no-download",
		"--no-warnings",
		"--flat-playlist",
		"--no-check-certificates",
		"--socket-timeout", "15",
	}

	if cfg.Proxy != "" {
		args = append(args, "--proxy", cfg.Proxy)
	}

	args = append(args, url)

	ctx, cancel := context.WithTimeout(context.Background(), 30*1000000000)
	defer cancel()

	cmd := exec.CommandContext(ctx, ytDlpPath, args...)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("analysis failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	var raw ytDlpJSON
	if err := json.Unmarshal(output, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return convertToVideoInfo(raw, url), nil
}

func convertToVideoInfo(raw ytDlpJSON, originalURL string) *VideoInfo {
	url := raw.WebpageURL
	if url == "" {
		url = originalURL
	}

	channel := raw.Channel
	if channel == "" {
		channel = raw.Uploader
	}

	info := &VideoInfo{
		ID:          raw.ID,
		URL:         url,
		Title:       raw.Title,
		Channel:     channel,
		Duration:    raw.Duration,
		DurationStr: formatDuration(raw.Duration),
		Thumbnail:   raw.Thumbnail,
		IsLive:      raw.IsLive,
		Description: raw.Description,
	}

	// If yt-dlp reports it as a playlist with entries, treat as playlist
	if raw.Type == "playlist" && len(raw.Entries) > 0 {
		info.Type = "playlist"
		info.VideoCount = len(raw.Entries)
		info.Entries = make([]VideoInfo, 0, len(raw.Entries))
		for _, entry := range raw.Entries {
			vi := convertToVideoInfo(entry, "")
			info.Entries = append(info.Entries, *vi)
		}
		// Use first entry info if playlist title is generic
		if info.Title == "" && len(info.Entries) > 0 {
			info.Title = "Playlist"
		}
		if info.Thumbnail == "" && len(info.Entries) > 0 {
			info.Thumbnail = info.Entries[0].Thumbnail
		}
	} else {
		info.Type = "video"
		info.VideoCount = 1
	}

	return info
}

func formatDuration(seconds float64) string {
	if seconds <= 0 {
		return ""
	}
	h := int(seconds) / 3600
	m := (int(seconds) % 3600) / 60
	s := int(seconds) % 60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%d:%02d", m, s)
}

// StartDownload runs a yt-dlp download and sends progress via callback
func StartDownload(ctx context.Context, opts DownloadOptions, onProgress func(DownloadProgress)) error {
	ytDlpPath, err := bootstrap.FindYtDlp()
	if err != nil {
		return fmt.Errorf("yt-dlp not found: %w", err)
	}

	cfg := config.Get()
	args := buildDownloadArgs(opts, cfg)
	args = append(args, opts.URL)

	cmd := exec.CommandContext(ctx, ytDlpPath, args...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start download: %w", err)
	}

	// Read stdout and stderr concurrently
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			progress := parseProgressLine(line, opts.URL)
			if progress != nil {
				onProgress(*progress)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			progress := parseProgressLine(line, opts.URL)
			if progress != nil {
				onProgress(*progress)
			}
		}
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("download cancelled")
		}
		return fmt.Errorf("download failed: %w", err)
	}

	return nil
}

func buildDownloadArgs(opts DownloadOptions, cfg *config.Config) []string {
	// Include quality in filename to differentiate downloads at different qualities
	qualityTag := ""
	if opts.Quality != "" && opts.Quality != "source" {
		qualityTag = fmt.Sprintf(" [%s]", opts.Quality)
	}

	outTemplate := fmt.Sprintf("%s/%%(title)s%s.%%(ext)s", opts.OutputDir, qualityTag)

	args := []string{
		"--newline",
		"--no-warnings",
		"--no-playlist",
		"-o", outTemplate,
	}

	if opts.Mode == "audio" {
		args = append(args, "-x", "--audio-format", "mp3")
		switch opts.Quality {
		case "low":
			args = append(args, "--audio-quality", "9") // ~64-96 kbps
		case "med":
			args = append(args, "--audio-quality", "5") // ~128-192 kbps
		case "high":
			args = append(args, "--audio-quality", "2") // ~192-256 kbps
		default: // source
			args = append(args, "--audio-quality", "0") // best available
		}
	} else {
		switch opts.Quality {
		case "low":
			args = append(args, "-f", "bestvideo[height<=480][ext=mp4]+bestaudio[ext=m4a]/best[height<=480]",
				"--merge-output-format", "mp4")
		case "med":
			args = append(args, "-f", "bestvideo[height<=720][ext=mp4]+bestaudio[ext=m4a]/best[height<=720]",
				"--merge-output-format", "mp4")
		case "high":
			args = append(args, "-f", "bestvideo[height<=1080][ext=mp4]+bestaudio[ext=m4a]/best[height<=1080]",
				"--merge-output-format", "mp4")
		default: // source
			args = append(args, "-f", "bestvideo[ext=mp4]+bestaudio[ext=m4a]/bestvideo+bestaudio/best",
				"--merge-output-format", "mp4")
		}
	}

	if cfg.SkipDuplicates {
		args = append(args, "--no-overwrites")
		// Per-quality archive so different qualities aren't considered duplicates
		archivePath := fmt.Sprintf("%s/.download_archive_%s_%s", opts.OutputDir, opts.Mode, opts.Quality)
		args = append(args, "--download-archive", archivePath)
	} else {
		// Unique filenames via epoch seconds — never collide
		args[4] = fmt.Sprintf("%s/%%(title)s%s %%(epoch)s.%%(ext)s", opts.OutputDir, qualityTag)
		args = append(args, "--force-overwrites")
	}
	if cfg.ContinueDownloads {
		args = append(args, "-c")
	}
	if cfg.FragmentConcurrent > 0 {
		args = append(args, "--concurrent-fragments", fmt.Sprintf("%d", cfg.FragmentConcurrent))
	}
	if cfg.SpeedLimit != "" {
		args = append(args, "-r", cfg.SpeedLimit)
	}
	if cfg.Proxy != "" {
		args = append(args, "--proxy", cfg.Proxy)
	}

	// Subtitles
	if opts.Subtitles || cfg.DownloadSubtitles {
		lang := opts.SubtitleLang
		if lang == "" {
			lang = cfg.SubtitleLang
		}
		args = append(args, "--write-subs", "--sub-langs", lang)
		if opts.EmbedSubs || cfg.EmbedSubtitles {
			args = append(args, "--embed-subs")
		}
	}

	if opts.EmbedMetadata || cfg.EmbedMetadata {
		args = append(args, "--embed-thumbnail", "--embed-metadata", "--embed-chapters")
	}

	return args
}

func parseProgressLine(line string, url string) *DownloadProgress {
	line = strings.TrimRight(line, "\r\n")
	p := &DownloadProgress{
		ID:      url,
		LogLine: line,
	}

	matches := progressRegex.FindStringSubmatch(line)
	if matches != nil {
		p.Status = "downloading"
		fmt.Sscanf(matches[1], "%f", &p.Percent)
		p.Total = matches[2]
		p.Speed = matches[3]
		p.ETA = matches[4]
		return p
	}

	if mergeRegex.MatchString(line) {
		p.Status = "merging"
		p.Percent = 100
		// Capture merged file path
		if destMatch := destinationRegex.FindStringSubmatch(line); destMatch != nil {
			p.FilePath = destMatch[1]
		}
		return p
	}

	// Capture destination file path and extract title from filename
	if destMatch := destinationRegex.FindStringSubmatch(line); destMatch != nil {
		p.FilePath = destMatch[1]
		p.Title = extractTitleFromPath(destMatch[1])
	}

	if strings.Contains(line, "has already been downloaded") {
		p.Status = "complete"
		p.Percent = 100
		return p
	}

	if strings.Contains(line, "[download]") && strings.Contains(line, "100%") {
		p.Status = "downloading"
		p.Percent = 100
		return p
	}

	if strings.Contains(line, "ERROR") || strings.Contains(line, "error") {
		p.Status = "error"
		return p
	}

	// Return log line for non-progress output
	if strings.TrimSpace(line) != "" {
		p.Status = "info"
		return p
	}

	return nil
}

// extractTitleFromPath extracts a clean video title from a file path like
// "/path/to/Some Video Title [high].mp4" → "Some Video Title"
func extractTitleFromPath(filePath string) string {
	base := filepath.Base(filePath)
	// Remove extension
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	// Remove quality tag like " [high]", " [med]", " [low]", " [source]"
	qualityTag := regexp.MustCompile(`\s*\[(high|med|low|source)\]$`)
	name = qualityTag.ReplaceAllString(name, "")
	return strings.TrimSpace(name)
}
