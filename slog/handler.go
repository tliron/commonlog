package slog

import (
	contextpkg "context"
	"log/slog"
)

var MOCK_HANDLER = NewMockHandler()

//
// MockHandler
//

type MockHandler struct{}

func NewMockHandler() *MockHandler {
	return new(MockHandler)
}

// ([slog.Handler] interface)
func (self *MockHandler) Enabled(context contextpkg.Context, level slog.Level) bool {
	return false
}

// ([slog.Handler] interface)
func (self *MockHandler) Handle(context contextpkg.Context, record slog.Record) error {
	return nil
}

// ([slog.Handler] interface)
func (self *MockHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return self
}

// ([slog.Handler] interface)
func (self *MockHandler) WithGroup(name string) slog.Handler {
	return self
}
