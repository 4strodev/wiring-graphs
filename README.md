# Wiring graphs
This package is an autowiring and DI injection library. Designed to make use of generics and type safe resolvers.
It is a container based DI framework. That means that you will need access to the container to being able to resolve
dependencies.

## Features
- Circular dependencies detection: Thanks to the usage of graphs and DFS this library is able to detect circular
dependencies. That detection is meant to be produced early avoiding to crash on production. Circular dependencies
should be detected while developing the app.

- Type safe resolving: The container only holds dependencies all the logic of resolving is encapsulated in a
`container.Resolve` function that manages the type detection successfully.

## Considerations
- Reflection: If you need an extreme performance this library is not for you. All the execution of adding and resolving
dependencies is done with reflection. That consumes time and CPU, which is not ideal for ultra performant services.

## Usage

### Define dependencies
First you need to create a container and declare dependencies.
```go
import (
    "log/slog"
    "os"

	"github.com/4strodev/wiring_graphs/pkg/container"
	"github.com/go-playground/validator/v10"
)
func main() {
	cont := container.New()

    // Must is a helper function that allows to chaing mutliple
    // depnendencies declaration without error checking. Any error will panic
	cont.Must().
		Singleton(func() *slog.Logger {
			return slog.New(slog.NewJSONHandler(os.Stdout, nil))
		}).

}
```
### Resolve dependencies
```go
import (
    "log/slog"

    "github.com/4strodev/wiring_graphs/pkg/container"
)

func main() {
	logger, err := container.Resolve[*slog.Logger](cont)
}
```

### Fill structs
You can also declare a struct with tags (or not) and resolve all its first level fields.

```go
// pending
```
