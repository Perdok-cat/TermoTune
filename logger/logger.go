package logger

import (
    "fmt"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// Custom error type for TermoTune
type TermoTuneError struct {
    Message   string
    Component string
    Code      int
}

func (e *TermoTuneError) Error() string {
    return fmt.Sprintf("TermoTune[%s]: %s (code: %d)", e.Component, e.Message, e.Code)
}

// Initialize logger
func init() {
    config := zap.NewDevelopmentConfig()
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    
    var err error
    Logger, err = config.Build()
    if err != nil {
        panic(fmt.Sprintf("Failed to initialize logger: %v", err))
    }
}

// NewTermoTuneError creates a new TermoTune error
func NewTermoTuneError(message string) *TermoTuneError {
    return &TermoTuneError{
        Message:   message,
        Component: "TermoTune",
        Code:      1,
    }
}

// NewTermoTuneErrorWithComponent creates a new TermoTune error with component
func NewTermoTuneErrorWithComponent(message, component string) *TermoTuneError {
    return &TermoTuneError{
        Message:   message,
        Component: component,
        Code:      1,
    }
}

// NewTermoTuneErrorWithCode creates a new TermoTune error with code
func NewTermoTuneErrorWithCode(message string, code int) *TermoTuneError {
    return &TermoTuneError{
        Message:   message,
        Component: "TermoTune",
        Code:      code,
    }
}

// Logging functions
func LogError(err error) {
    Logger.Error("Error occurred", zap.Error(err))
}

func LogErrorWithFields(err error, fields ...zap.Field) {
    Logger.Error("Error occurred", append([]zap.Field{zap.Error(err)}, fields...)...)
}

func LogInfo(msg string, fields ...zap.Field) {
    Logger.Info(msg, fields...)
}

func LogWarn(msg string, fields ...zap.Field) {
    Logger.Warn(msg, fields...)
}

func LogDebug(msg string, fields ...zap.Field) {
    Logger.Debug(msg, fields...)
}

func LogFatal(msg string, fields ...zap.Field) {
    Logger.Fatal(msg, fields...)
}

// Convenience functions for common operations
func LogDatabaseError(operation string, err error) {
    Logger.Error("Database operation failed",
        zap.String("operation", operation),
        zap.Error(err),
    )
}

func LogMusicOperation(operation, musicName string, err error) {
    if err != nil {
        Logger.Error("Music operation failed",
            zap.String("operation", operation),
            zap.String("music", musicName),
            zap.Error(err),
        )
    } else {
        Logger.Info("Music operation successful",
            zap.String("operation", operation),
            zap.String("music", musicName),
        )
    }
}

func LogPlaylistOperation(operation, playlistName string, err error) {
    if err != nil {
        Logger.Error("Playlist operation failed",
            zap.String("operation", operation),
            zap.String("playlist", playlistName),
            zap.Error(err),
        )
    } else {
        Logger.Info("Playlist operation successful",
            zap.String("operation", operation),
            zap.String("playlist", playlistName),
        )
    }
}

// Cleanup function
func Sync() {
    if Logger != nil {
        Logger.Sync()
    }
}