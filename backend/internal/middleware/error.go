package middleware // <--- THIS IS WHAT IS MISSING

import (
    "log/slog"
    "net/http"
)


func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                slog.Error("panic recovered",
                    "error", err,
                    "method", r.Method,
                    "path", r.URL.Path,
                    "remote_addr", r.RemoteAddr,
                )
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}