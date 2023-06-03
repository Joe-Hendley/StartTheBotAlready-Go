package taunt

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Joe-Hendley/StartTheBotAlready-Go/internal/voice"
)

var RequestChan chan<- voice.Request
var buffers map[int]*voice.Buffer

func init() {
	filePaths := getFilePaths()
	buffers = populateBuffers(filePaths)
}

func play(tauntID int, guildID, channelID string) {
	if RequestChan == nil {
		log.Fatal("taunt called before voice channel registered")
	}

	buffer, ok := buffers[tauntID]
	if !ok {
		log.Fatalf("taunt [%d] not found", tauntID)
	}

	RequestChan <- voice.NewRequest(buffer, guildID, channelID)
}

func getFilePaths() []string {
	dir, _ := os.Getwd()
	path := filepath.Join(dir, "audio", "taunts")

	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	paths := make([]string, 0, 42)

	for _, file := range entries {
		if strings.HasSuffix(file.Name(), ".dca") {
			paths = append(paths, filepath.Join(path, file.Name()))
		}
	}

	return paths
}

func populateBuffers(paths []string) map[int]*voice.Buffer {
	buffers := make(map[int]*voice.Buffer, 42)
	for idx, path := range paths {
		buffer, err := voice.LoadSound(path)
		if err != nil {
			log.Fatalf("failed to load sound at path %s: %v", path, err)
		}
		buffers[idx+1] = buffer
	}
	return buffers
}
