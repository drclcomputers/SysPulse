package v2

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}

type Logger struct {
	mu       sync.Mutex
	out      io.Writer
	file     *os.File
	level    LogLevel
	dir      string
	filename string
}

type LogEntry struct {
	Level     LogLevel
	Message   string
	Timestamp time.Time
	File      string
	Line      int
	Function  string
}

func New(logDir string, level LogLevel) (*Logger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	filename := fmt.Sprintf("syspulse_%s.log", time.Now().Format("2006-01-02"))
	filepath := filepath.Join(logDir, filename)

	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	return &Logger{
		out:      io.MultiWriter(os.Stderr, file),
		file:     file,
		level:    level,
		dir:      logDir,
		filename: filename,
	}, nil
}

func (l *Logger) log(level LogLevel, message string) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	pc, file, line, _ := runtime.Caller(2)
	fn := runtime.FuncForPC(pc).Name()

	entry := LogEntry{
		Level:     level,
		Message:   message,
		Timestamp: time.Now(),
		File:      file,
		Line:      line,
		Function:  fn,
	}

	fmt.Fprintf(l.out, "[%s] %s - %s:%d - %s: %s\n",
		entry.Timestamp.Format("2006-01-02 15:04:05"),
		entry.Level,
		filepath.Base(entry.File),
		entry.Line,
		entry.Function,
		entry.Message,
	)
}

func (l *Logger) Debug(message string) { l.log(DEBUG, message) }
func (l *Logger) Info(message string)  { l.log(INFO, message) }
func (l *Logger) Warn(message string)  { l.log(WARN, message) }
func (l *Logger) Error(message string) { l.log(ERROR, message) }
func (l *Logger) Fatal(message string) {
	l.log(FATAL, message)
	os.Exit(1)
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func (l *Logger) RotateLog() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		if err := l.file.Close(); err != nil {
			return fmt.Errorf("failed to close current log file: %v", err)
		}
	}

	filename := fmt.Sprintf("syspulse_%s.log", time.Now().Format("2006-01-02"))
	if filename == l.filename {
		return nil
	}

	filepath := filepath.Join(l.dir, filename)
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open new log file: %v", err)
	}

	l.file = file
	l.out = io.MultiWriter(os.Stderr, file)
	l.filename = filename

	return nil
}
