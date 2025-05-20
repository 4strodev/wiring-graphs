package container_test

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/4strodev/wiring/pkg/container"
	"github.com/stretchr/testify/require"
)

type MyService struct {
}

func (s MyService) SayHi() string {
	return "hi!"
}

func NewService() MyService {
	return MyService{}
}

func TestDetectCircularDependencies(t *testing.T) {
	cont := container.New()

	err := cont.AddDependencies(func(s io.Writer) (*slog.Logger, error) {
		return slog.New(slog.NewJSONHandler(s, nil)), errors.New("an artificial error")
	}, func(logger *slog.Logger) MyService {
		return MyService{}
	}, func() io.Writer {
		return os.Stdout
	})
	require.NoError(t, err, "container should allow add dependencies without errors")

	_, err = cont.DetectCircularDependencies()
	require.NoError(t, err, "circular dependencies should not be detected")
}

func TestResolve(t *testing.T) {
	cont := container.New()

	cont.AddDependency(NewService)

	_, err := container.Resolve[MyService](cont)
	require.NoError(t, err, "MyService should be resolved")
}

func TestResolve_WithDependencies(t *testing.T) {
	cont := container.New()

	dependantResolver := func(s MyService) *bytes.Buffer {
		return bytes.NewBuffer([]byte{})
	}

	err := cont.AddDependencies(dependantResolver, NewService)
	require.NoError(t, err)

	_, err = container.Resolve[MyService](cont)
	require.NoError(t, err, "MyService should be resolved")

	buf, err := container.Resolve[*bytes.Buffer](cont)
	require.NoError(t, err)
	require.NotNil(t, buf)
}
