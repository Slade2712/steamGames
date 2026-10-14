
package logs

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"

    "steam_games/internal/models"
)

// Protège le fichier contre les écritures simultanées.
var logMutex sync.Mutex

// BuildLogPath construit le chemin selon la date et l'heure.

func BuildLogPath(now time.Time) string {
    year := fmt.Sprintf("%04d", now.Year())
    month := fmt.Sprintf("%02d", int(now.Month()))
    day := fmt.Sprintf("%02d", now.Day())

    minute := (now.Minute() / 10) * 10
    hourMinute := fmt.Sprintf("%02d-%02d", now.Hour(), minute)

    return filepath.Join("logs", year, month, day, hourMinute)
}


// CreateLogDirectory crée les dossiers manquants.
func CreateLogDirectory(now time.Time) (string, error) {
    path := BuildLogPath(now)

    err := os.MkdirAll(path, 0755)
    if err != nil {
        return "", err
    }

    return path, nil
}

// SaveLog ajoute une requête dans le fichier de la minute.

func SaveLog(entry models.LogEntry) error {
    logMutex.Lock()
    defer logMutex.Unlock()

    path, err := CreateLogDirectory(time.Now())
    if err != nil {
        return err
    }

    filePath := filepath.Join(path, "requests.json")

    entries := []models.LogEntry{}

    data, err := os.ReadFile(filePath)
    if err == nil {
        err = json.Unmarshal(data, &entries)
        if err != nil {
            return err
        }
    } else if !os.IsNotExist(err) {
        return err
    }

    entries = append(entries, entry)

    data, err = json.MarshalIndent(entries, "", "    ")
    if err != nil {
        return err
    }

    return os.WriteFile(filePath, data, 0644)
}
