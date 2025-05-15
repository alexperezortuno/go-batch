package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/alexperezortuno/go-batch/internal/config"
)

// Niveles de log
const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
)

// Logger es la estructura principal del logger
type Logger struct {
	Level      string
	File       *os.File
	Caller     bool
	TimeFormat string
	Output     io.Writer
	Mu         sync.Mutex
	Fields     []interface{}
}

// NewLogger crea una nueva instancia de Logger
func NewLogger(cfg config.Config) *Logger {
	l := &Logger{
		Level:      strings.ToUpper(cfg.Logger.Level),
		Caller:     cfg.Logger.Caller,
		TimeFormat: cfg.Logger.TimeFormat,
	}

	if l.TimeFormat == "" {
		l.TimeFormat = time.RFC3339
	}

	// Configurar salida
	if cfg.Logger.FilePath != "" {
		dir := filepath.Dir(cfg.Logger.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("No se pudo crear directorio para logs: %v. Usando stdout", err)
			l.Output = os.Stdout
		} else {
			file, err := os.OpenFile(cfg.Logger.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Printf("No se pudo abrir archivo de log: %v. Usando stdout", err)
				l.Output = os.Stdout
			} else {
				l.File = file
				l.Output = io.MultiWriter(os.Stdout, file)
			}
		}
	} else {
		l.Output = os.Stdout
	}

	return l
}

// Close cierra los recursos del logger
func (l *Logger) Close() error {
	if l.File != nil {
		return l.File.Close()
	}
	return nil
}

// Debug registra un mensaje de nivel DEBUG
func (l *Logger) Debug(msg string, fields ...interface{}) {
	if l.shouldLog(LevelDebug) {
		l.log(LevelDebug, msg, fields...)
	}
}

// Info registra un mensaje de nivel INFO
func (l *Logger) Info(msg string, fields ...interface{}) {
	if l.shouldLog(LevelInfo) {
		l.log(LevelInfo, msg, fields...)
	}
}

// Warn registra un mensaje de nivel WARN
func (l *Logger) Warn(msg string, fields ...interface{}) {
	if l.shouldLog(LevelWarn) {
		l.log(LevelWarn, msg, fields...)
	}
}

// Error registra un mensaje de nivel ERROR
func (l *Logger) Error(msg string, fields ...interface{}) {
	if l.shouldLog(LevelError) {
		l.log(LevelError, msg, fields...)
	}
}

// shouldLog determina si un mensaje debe ser registrado según el nivel
func (l *Logger) shouldLog(level string) bool {
	levels := map[string]int{
		LevelDebug: 4,
		LevelInfo:  3,
		LevelWarn:  2,
		LevelError: 1,
	}

	return levels[level] <= levels[l.Level]
}

// log escribe el mensaje en el output
func (l *Logger) log(level, msg string, fields ...interface{}) {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	now := time.Now().Format(l.TimeFormat)
	callerInfo := ""

	if l.Caller {
		if _, file, line, ok := runtime.Caller(2); ok {
			callerInfo = fmt.Sprintf(" [%s:%d]", filepath.Base(file), line)
		}
	}

	// Formatear campos adicionales
	var formattedFields string
	if len(fields) > 0 {
		var sb strings.Builder
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				if i > 0 {
					sb.WriteString(" ")
				}
				sb.WriteString(fmt.Sprintf("%v=%v", fields[i], fields[i+1]))
			}
		}
		formattedFields = " " + sb.String()
	}

	logMsg := fmt.Sprintf("%s [%s]%s %s%s\n", now, level, callerInfo, msg, formattedFields)

	if _, err := l.Output.Write([]byte(logMsg)); err != nil {
		log.Printf("Error escribiendo log: %v", err)
	}
}

// WithFields crea un nuevo logger con campos adicionales
func (l *Logger) WithFields(fields ...interface{}) *Logger {
	return &Logger{
		Level:      l.Level,
		Output:     l.Output,
		File:       l.File,
		Caller:     l.Caller,
		TimeFormat: l.TimeFormat,
	}
}
