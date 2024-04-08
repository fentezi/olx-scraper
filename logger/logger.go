package logger

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func (p *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}
	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", " ")
		if err != nil {
			return err
		}

	}
	timeStr := r.Time.Format("[15:04:05]")
	msg := color.CyanString(r.Message)

	p.l.Println(timeStr, level, msg, color.WhiteString(string(b)))

	return nil
}

func NewPrettyHandler(out io.Writer, opts slog.HandlerOptions) *PrettyHandler {
	return &PrettyHandler{
		Handler: slog.NewTextHandler(out, &opts),
		l:       log.New(out, "", 0),
	}
}

func Logger() *slog.Logger {
	env := getEnvVariable("ENV")

	var handler slog.Handler
	switch env {
	case "local":
		handler = newConsoleHandler(slog.LevelInfo)
	case "prod":
		handler = newFileHandler()
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}

	return slog.New(handler)
}

func getEnvVariable(key string) string {
	err := godotenv.Load()
	if err != nil {
		return ""
	}

	return os.Getenv(key)
}

func newConsoleHandler(level slog.Level) slog.Handler {
	return NewPrettyHandler(os.Stdout, slog.HandlerOptions{Level: level})
}

func newFileHandler() slog.Handler {
	file, err := os.OpenFile("slog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug})
}
