package main

import (
	"errors"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/4strodev/wiring/pkg/container"
)

type MyService struct {
}

func main() {
	cont := container.New()

	err := cont.AddDependencies(func(s io.Writer) (*slog.Logger, error) {
		return slog.New(slog.NewJSONHandler(s, nil)), errors.New("an artificial error")
	}, func(logger *slog.Logger) MyService {
		return MyService{}
	}, func() io.Writer {
		return os.Stdout
	})

	cont.Graph.DetectCircularRelations()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("container works")
}
