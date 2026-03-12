package bootstrap

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"GotubeDownloader/internal/config"
)

// BootstrapStatus reports whether dependencies are available
type BootstrapStatus struct {
	YtDlpReady  bool   `json:"ytDlpReady"`
	FfmpegReady bool   `json:"ffmpegReady"`
	YtDlpPath   string `json:"ytDlpPath"`
	FfmpegPath  string `json:"ffmpegPath"`
	Error       string `json:"error"`
}

func ytDlpBinaryName() string {
	if runtime.GOOS == "windows" {
		return "yt-dlp.exe"
	}
	return "yt-dlp"
}

func ffmpegBinaryName() string {
	if runtime.GOOS == "windows" {
		return "ffmpeg.exe"
	}
	return "ffmpeg"
}

func getBinDir() (string, error) {
	dir, err := config.GetConfigDir()
	if err != nil {
		return "", err
	}
	binDir := filepath.Join(dir, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return "", err
	}
	return binDir, nil
}

// FindYtDlp locates the yt-dlp binary
func FindYtDlp() (string, error) {
	// 1. Check local bin dir
	binDir, err := getBinDir()
	if err == nil {
		localPath := filepath.Join(binDir, ytDlpBinaryName())
		if _, err := os.Stat(localPath); err == nil {
			return localPath, nil
		}
	}

	// 2. Check system PATH
	path, err := exec.LookPath(ytDlpBinaryName())
	if err == nil {
		return path, nil
	}

	// 3. Common homebrew paths (macOS)
	if runtime.GOOS == "darwin" {
		for _, p := range []string{"/opt/homebrew/bin/yt-dlp", "/usr/local/bin/yt-dlp"} {
			if _, err := os.Stat(p); err == nil {
				return p, nil
			}
		}
	}

	return "", fmt.Errorf("yt-dlp not found")
}

// FindFfmpeg locates the ffmpeg binary
func FindFfmpeg() (string, error) {
	// 1. Check local bin dir
	binDir, err := getBinDir()
	if err == nil {
		localPath := filepath.Join(binDir, ffmpegBinaryName())
		if _, err := os.Stat(localPath); err == nil {
			return localPath, nil
		}
	}

	// 2. Check system PATH
	path, err := exec.LookPath(ffmpegBinaryName())
	if err == nil {
		return path, nil
	}

	// 3. Common homebrew paths (macOS)
	if runtime.GOOS == "darwin" {
		for _, p := range []string{"/opt/homebrew/bin/ffmpeg", "/usr/local/bin/ffmpeg"} {
			if _, err := os.Stat(p); err == nil {
				return p, nil
			}
		}
	}

	return "", fmt.Errorf("ffmpeg not found")
}

// CheckStatus checks if yt-dlp and ffmpeg are available
func CheckStatus() BootstrapStatus {
	status := BootstrapStatus{}

	ytPath, err := FindYtDlp()
	if err == nil {
		status.YtDlpReady = true
		status.YtDlpPath = ytPath
	}

	ffPath, err := FindFfmpeg()
	if err == nil {
		status.FfmpegReady = true
		status.FfmpegPath = ffPath
	}

	return status
}

func getYtDlpDownloadURL() string {
	switch runtime.GOOS {
	case "windows":
		if runtime.GOARCH == "arm64" {
			return "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_win_arm64.exe"
		}
		return "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"
	case "darwin":
		return "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos"
	default:
		return "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux"
	}
}

// DownloadYtDlp downloads yt-dlp to the local bin directory
func DownloadYtDlp() (string, error) {
	binDir, err := getBinDir()
	if err != nil {
		return "", fmt.Errorf("failed to get bin directory: %w", err)
	}

	destPath := filepath.Join(binDir, ytDlpBinaryName())
	url := getYtDlpDownloadURL()

	if err := downloadFile(url, destPath); err != nil {
		return "", fmt.Errorf("failed to download yt-dlp: %w", err)
	}

	// Make executable on unix
	if runtime.GOOS != "windows" {
		if err := os.Chmod(destPath, 0755); err != nil {
			return "", fmt.Errorf("failed to set permissions: %w", err)
		}
	}

	// Validate
	cmd := exec.Command(destPath, "--version")
	if err := cmd.Run(); err != nil {
		os.Remove(destPath)
		return "", fmt.Errorf("downloaded yt-dlp is not functional: %w", err)
	}

	return destPath, nil
}

func downloadFile(url, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
