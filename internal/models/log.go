package models

type RequestLog struct {
    Method string `json:"method"`
    Path   string `json:"path"`
    IP     string `json:"ip"`
}

type ResponseLog struct {
    Status int `json:"status"`
}

type LogEntry struct {
    Timestamp  string      `json:"timestamp"`
    DurationMs int64       `json:"duration_ms"`
    Request    RequestLog  `json:"request"`
    Response   ResponseLog `json:"response"`
}