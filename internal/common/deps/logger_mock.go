package deps

import (
	"context"
)

type MockLogger struct {
}

func (l *MockLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {

}

func (l *MockLogger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {

}

func (l *MockLogger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {

}
