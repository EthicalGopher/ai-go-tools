package tts

import (
	"fmt"
	"os"
	"path/filepath"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
)

func TTS(input, filename, dir string) (string, error) {
	// Ensure the output directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// Construct full output path
	fullPath := filepath.Join(dir, filename+".mp3")

	// Delete existing file if it exists
	if _, err := os.Stat(fullPath); err == nil {
		if err := os.Remove(fullPath); err != nil {
			return "", fmt.Errorf("failed to remove existing file: %v", err)
		}
	}

	// Initialize TTS
	speech := htgotts.Speech{
		Folder:   dir,
		Language: voices.English, // Or make this a parameter if needed
		Handler:  &handlers.Native{},
	}

	// Generate speech file
	if _, err := speech.CreateSpeechFile(input, filename); err != nil {
		return "", fmt.Errorf("failed to create speech: %v", err)
	}

	return "success", nil
}
