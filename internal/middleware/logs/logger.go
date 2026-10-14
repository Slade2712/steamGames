
package logs

import (
    "log"
    "net/http"
    "time"

    "steam_games/internal/models"
    logService "steam_games/internal/services/logs"
)

type LogWriter struct {
    http.ResponseWriter
    Status int
}

func (w *LogWriter) WriteHeader(status int) {
    w.Status = status
    w.ResponseWriter.WriteHeader(status)
}

func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        writer := &LogWriter{
            ResponseWriter: w,
            Status: http.StatusOK,
        }

        next.ServeHTTP(writer, r)

        entry := models.LogEntry{
            Timestamp: start.Format(time.RFC3339),
            DurationMs: time.Since(start).Milliseconds(),
            Request: models.RequestLog{
                Method: r.Method,
                Path: r.URL.Path,
                IP: r.RemoteAddr,
            },
            Response: models.ResponseLog{
                Status: writer.Status,
            },
        }

        err := logService.SaveLog(entry)
        if err != nil {
            log.Println(err)
        }
    })
}
