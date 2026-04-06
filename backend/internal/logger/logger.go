import (
    "log/slog"
    "os"
)

func initLogger() {
    // In production (Railway), we want JSON. In local dev, maybe Text.
    var handler slog.Handler
    if os.Getenv("RAILWAY_ENVIRONMENT") != "" {
        handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
    } else {
        handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
    }

    logger := slog.New(handler)
    slog.SetDefault(logger)
}