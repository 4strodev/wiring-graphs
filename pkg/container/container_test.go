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

type MyDeps struct {
	Service MyService
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

func TestDetectCircularDependencies_SelfReference(t *testing.T) {
	cont := container.New()

	cont.AddDependency(func(s MyService) MyService {
		return MyService{}
	})

	_, err := cont.DetectCircularDependencies()
	require.Error(t, err, "should detect self reference")
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

func TestFill(t *testing.T) {
	cont := container.New()

	dependantResolver := func(s MyService) *bytes.Buffer {
		return bytes.NewBuffer([]byte{})
	}

	err := cont.AddDependencies(dependantResolver, NewService)
	require.NoError(t, err)

	var deps MyDeps
	err = cont.Fill(&deps)
	require.NoError(t, err, "struct should be filled correctly")
}

func TestResolveToken(t *testing.T) {
	cont := container.New()

	dependantResolver := func() *bytes.Buffer {
		return bytes.NewBuffer([]byte{})
	}

	err := cont.AddTokenDependency("buffer", dependantResolver)
	require.NoError(t, err)

	buf, err := container.ResolveToken[*bytes.Buffer](cont, "buffer")
	require.NoError(t, err)
	require.NotNil(t, buf)
}

func TestResolveToken_WithDependencies(t *testing.T) {
	cont := container.New()

	dependantResolver := func(service MyService) *bytes.Buffer {
		return bytes.NewBuffer([]byte{})
	}

	err := cont.AddTokenDependency("buffer", dependantResolver)
	require.NoError(t, err)

	cont.AddDependency(NewService)

	buf, err := container.ResolveToken[*bytes.Buffer](cont, "buffer")
	require.NoError(t, err)
	require.NotNil(t, buf)
}
