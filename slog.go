package kafka

import (
	"fmt"
	"log/slog"
)

type SlogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger() *SlogLogger {
	return &SlogLogger{logger: slog.Default()}
}

func (s *SlogLogger) With(args ...interface{}) LoggerInterface {
	if len(args)%2 != 0 {
		// If odd number of args, add empty string as key
		args = append(args, "")
	}

	attrs := make([]any, 0, len(args))
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			key, ok := args[i].(string)
			if !ok {
				key = fmt.Sprintf("%v", args[i])
			}
			attrs = append(attrs, key, args[i+1])
		}
	}

	return &SlogLogger{logger: s.logger.With(attrs...)}
}

func (s *SlogLogger) Debug(args ...interface{}) {
	s.logger.Debug(fmt.Sprint(args...))
}

func (s *SlogLogger) Info(args ...interface{}) {
	s.logger.Info(fmt.Sprint(args...))
}

func (s *SlogLogger) Warn(args ...interface{}) {
	s.logger.Warn(fmt.Sprint(args...))
}

func (s *SlogLogger) Error(args ...interface{}) {
	s.logger.Error(fmt.Sprint(args...))
}

func (s *SlogLogger) Debugf(format string, args ...interface{}) {
	s.logger.Debug(fmt.Sprintf(format, args...))
}

func (s *SlogLogger) Infof(format string, args ...interface{}) {
	s.logger.Info(fmt.Sprintf(format, args...))
}

func (s *SlogLogger) Warnf(format string, args ...interface{}) {
	s.logger.Warn(fmt.Sprintf(format, args...))
}

func (s *SlogLogger) Errorf(format string, args ...interface{}) {
	s.logger.Error(fmt.Sprintf(format, args...))
}

func (s *SlogLogger) Infow(msg string, keysAndValues ...interface{}) {
	attrs := s.buildAttrs(keysAndValues...)
	s.logger.Info(msg, attrs...)
}

func (s *SlogLogger) Errorw(msg string, keysAndValues ...interface{}) {
	attrs := s.buildAttrs(keysAndValues...)
	s.logger.Error(msg, attrs...)
}

func (s *SlogLogger) Warnw(msg string, keysAndValues ...interface{}) {
	attrs := s.buildAttrs(keysAndValues...)
	s.logger.Warn(msg, attrs...)
}

func (s *SlogLogger) buildAttrs(keysAndValues ...interface{}) []any {
	if len(keysAndValues)%2 != 0 {
		keysAndValues = append(keysAndValues, "")
	}

	attrs := make([]any, 0, len(keysAndValues))
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key, ok := keysAndValues[i].(string)
			if !ok {
				key = fmt.Sprintf("%v", keysAndValues[i])
			}
			attrs = append(attrs, key, keysAndValues[i+1])
		}
	}
	return attrs
}
