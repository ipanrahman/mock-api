package log

import (
	"fmt"
	"log/slog"
	"mock-api/internal/common/stacktrace"
	"strings"
)

// Level is the default log severity level.
var Level = new(slog.LevelVar)

// Replace is the default replace attribute.
// https://sourcegraph.com/github.com/uber-go/zap/-/blob/zapcore/entry.go?L117
var Replace = func(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.SourceKey {
		if source, ok := a.Value.Any().(*slog.Source); ok {
			idx := strings.LastIndexByte(source.File, '/')
			if idx == -1 {
				return a
			}
			// Find the penultimate separator.
			idx = strings.LastIndexByte(source.File[:idx], '/')
			if idx == -1 {
				return a
			}
			source.File = source.File[idx+1:]
		}
	}
	return a
}

// Initializes the slog configuration.
func init() {
	Level.Set(slog.LevelInfo)
}

func Error(err error) slog.Attr {
	return slog.String("error", fmt.Sprintf("%+v", err))
}

func Stack(key string) slog.Attr {
	stack := stacktrace.TakeStacktrace(20 /* n */, 3 /* skip */)
	return slog.Any(key, stack)
}
