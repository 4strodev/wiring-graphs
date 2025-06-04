package container_test

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/4strodev/wiring/pkg/container"
	"github.com/4strodev/wiring/pkg/internal/testutils"
	"github.com/stretchr/testify/require"
)


func TestDetectCircularDependencies(t *testing.T) {
	cont := container.New()

	err := cont.AddDependencies(func(s io.Writer) (*slog.Logger, error) {
		return slog.New(slog.NewJSONHandler(s, nil)), errors.New("an artificial error")
	}, func(logger *slog.Logger) testutils.MyService {
		return testutils.MyService{}
	}, func() io.Writer {
		return os.Stdout
	})
	require.NoError(t, err, "container should allow add dependencies without errors")

	_, err = cont.DetectCircularDependencies()
	require.NoError(t, err, "circular dependencies should not be detected")
}

func TestDetectCircularDependencies_SelfReference(t *testing.T) {
	cont := container.New()

	cont.AddDependency(func(s testutils.MyService) testutils.MyService {
		return testutils.MyService{}
	})

	_, err := cont.DetectCircularDependencies()
	require.Error(t, err, "should detect self reference")
}

func TestResolve(t *testing.T) {
	cont := container.New()

	cont.AddDependency(testutils.NewService)

	_, err := container.Resolve[testutils.MyService](cont)
	require.NoError(t, err, "testutils.MyService should be resolved")
}

func TestResolve_WithDependencies(t *testing.T) {
	cont := container.New()

	dependantResolver := func(s testutils.MyService) *bytes.Buffer {
		return bytes.NewBuffer([]byte{})
	}

	err := cont.AddDependencies(dependantResolver, testutils.NewService)
	require.NoError(t, err)

	_, err = container.Resolve[testutils.MyService](cont)
	require.NoError(t, err, "testutils.MyService should be resolved")

	buf, err := container.Resolve[*bytes.Buffer](cont)
	require.NoError(t, err)
	require.NotNil(t, buf)
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

	dependantResolver := func(service testutils.MyService) *bytes.Buffer {
		return bytes.NewBuffer([]byte{})
	}

	err := cont.AddTokenDependency("buffer", dependantResolver)
	require.NoError(t, err)

	cont.AddDependency(testutils.NewService)

	buf, err := container.ResolveToken[*bytes.Buffer](cont, "buffer")
	require.NoError(t, err)
	require.NotNil(t, buf)
}
