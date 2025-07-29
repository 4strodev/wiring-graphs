package container

import (
	"reflect"

	"github.com/4strodev/wiring_graphs/pkg/errors"
	"github.com/4strodev/wiring_graphs/pkg/internal/collections/graph"
	"github.com/4strodev/wiring_graphs/pkg/resolver"
)

type resolverConfig struct {
	singleton  bool
	resolved   bool
	node       *graph.Node[resolver.DependencyResolver[any]]
	savedValue reflect.Value
}

type Container struct {
	graph.Graph[resolver.DependencyResolver[any]]
	typeIndex  map[reflect.Type]*resolverConfig
	tokenIndex map[string]*resolverConfig
	connected  bool
}

// Retuns a new container and sets a default dependency that allows
// to inject this container as a parameter on the resolvers
func New() *Container {
	container := &Container{
		Graph:      graph.NewGraph[resolver.DependencyResolver[any]](),
		typeIndex:  make(map[reflect.Type]*resolverConfig),
		tokenIndex: make(map[string]*resolverConfig),
	}

	// Allow resolvers to inject container
	container.Dependencies(func() *Container {
		return container
	})

	return container
}

func (c *Container) Dependencies(resolvers ...any) error {
	c.connected = false
	for _, res := range resolvers {
		config, err := buildConfig(res)
		if err != nil {
			return err
		}

		resType := config.node.Val.Type()
		_, exists := c.typeIndex[resType]
		if exists {
			return errors.Errorf(errors.E_REDECLARED_DEPENDENCY, "dependency for this type already exists: %v", resType)
		}

		c.Add(config.node)
		c.typeIndex[resType] = config
	}

	return nil
}

func (c *Container) Singleton(resolvers ...any) error {
	c.connected = false
	for _, res := range resolvers {
		config, err := buildConfig(res)
		if err != nil {
			return err
		}

		resType := config.node.Val.Type()
		_, exists := c.typeIndex[resType]
		if exists {
			return errors.Errorf(errors.E_REDECLARED_DEPENDENCY, "dependency for this type already exists: %v", resType)
		}

		c.Add(config.node)
		config.singleton = true
		c.typeIndex[resType] = config
	}

	return nil
}

func (c *Container) TokenSingleton(dependencies map[string]any) error {
	c.connected = false
	for token, res := range dependencies {
		config, err := buildConfig(res)
		if err != nil {
			return err
		}

		_, exists := c.tokenIndex[token]
		if exists {
			return errors.Errorf(errors.E_REDECLARED_DEPENDENCY, "dependency for token type already exists: %s", token)
		}

		c.Add(config.node)
		config.singleton = true
		c.tokenIndex[token] = config
	}

	return nil
}

func (c *Container) Token(dependencies map[string]any) error {
	c.connected = false
	for token, res := range dependencies {
		config, err := buildConfig(res)
		if err != nil {
			return err
		}

		_, exists := c.tokenIndex[token]
		if exists {
			return errors.Errorf(errors.E_REDECLARED_DEPENDENCY, "dependency for token already exists: %s", token)
		}

		c.Add(config.node)
		c.tokenIndex[token] = config
	}

	return nil
}

func buildConfig(res any) (*resolverConfig, error) {
	if !resolver.IsValid(res) {
		return nil, errors.Errorf(errors.E_INVALID_RESOLVER, "Invalid resolver")
	}

	builder := resolver.DependencyResolver[any]{
		Resolver: res,
	}

	node := graph.NewNode(builder)
	config := resolverConfig{
		node: node,
	}

	return &config, nil
}

func (c *Container) ensureNodesConnected() error {
	if c.connected {
		return nil
	}

	return c.setConnections()
}

func (c *Container) DetectCircularDependencies() ([]*graph.Node[resolver.DependencyResolver[any]], error) {
	err := c.ensureNodesConnected()
	if err != nil {
		return nil, err
	}
	path, detected := c.Graph.DetectCircularRelations()
	if detected {
		return path, nil
	}

	return nil, nil
}

func (c Container) getNodeFor(t reflect.Type) (*graph.Node[resolver.DependencyResolver[any]], error) {
	node, ok := c.typeIndex[t]
	if !ok {
		return nil, errors.Errorf(errors.E_DEPENDENCY_NOT_FOUND, "dependency for %v not fonud", t)
	}

	return node.node, nil
}

// setConnections stablishes connections between nodes and
// look for circular dependencies
func (c *Container) setConnections() error {
	for node := range c.Graph.GetNodes() {
		resType := reflect.TypeOf(node.Val.Resolver)
		for i := 0; i < resType.NumIn(); i++ {
			dependencyType := resType.In(i)

			if node.Val.Type() == dependencyType {
				return errors.Errorf(
					errors.E_CIRCULAR_DEPENDENCY,
					"circular dependency found: %v",
					[]*graph.Node[resolver.DependencyResolver[any]]{
						node,
					})
			}

			dependencyNode, err := c.getNodeFor(dependencyType)
			if err != nil {
				return err
			}

			if dependencyNode.IsConnectedWith(node) {
				return errors.Errorf(
					errors.E_CIRCULAR_DEPENDENCY,
					"circular dependency found: %v",
					[]*graph.Node[resolver.DependencyResolver[any]]{
						node,
						dependencyNode,
						node,
					},
				)
			}

			c.Connect(dependencyNode, node, graph.OUT)
		}
	}

	cicle, hasCicle := c.DetectCircularRelations()
	if hasCicle {
		return errors.Errorf(errors.E_CIRCULAR_DEPENDENCY, "circular dependency found: %v", cicle)
	}

	c.connected = true
	return nil
}

func (c Container) resolve(t reflect.Type) (resolvedValue reflect.Value, err error) {
	node, ok := c.typeIndex[t]
	if !ok {
		err = errors.Errorf(errors.E_DEPENDENCY_NOT_FOUND, "dependency not found for type %v", t)
		return
	}

	if node.singleton && node.resolved {
		resolvedValue = node.savedValue
		return
	}

	dependencyResolver := node.node.Val
	inputTypes := dependencyResolver.Input()
	inputArgs := []reflect.Value{}

	for _, inputType := range inputTypes {
		arg, err := c.resolve(inputType)
		if err != nil {
			return resolvedValue, err
		}

		inputArgs = append(inputArgs, arg)
	}

	resolvedValue, err = resolver.Execute(dependencyResolver, inputArgs)
	if err != nil {
		return
	}
	node.resolved = true

	if node.singleton {
		node.savedValue = resolvedValue
	}

	return
}

func (c Container) resolveToken(token string) (resolvedValue reflect.Value, err error) {
	node, ok := c.tokenIndex[token]
	if !ok {
		err = errors.Errorf(errors.E_DEPENDENCY_NOT_FOUND, "dependency not found for token '%s'", token)
		return
	}

	if node.singleton && node.resolved {
		resolvedValue = node.savedValue
		return
	}

	dependencyResolver := node.node.Val
	inputTypes := dependencyResolver.Input()
	inputArgs := []reflect.Value{}

	for _, inputType := range inputTypes {
		arg, err := c.resolve(inputType)
		if err != nil {
			return resolvedValue, err
		}

		inputArgs = append(inputArgs, arg)
	}

	resolvedValue, err = resolver.Execute(dependencyResolver, inputArgs)
	if err != nil {
		return
	}
	node.resolved = true

	if node.singleton {
		node.savedValue = resolvedValue
	}

	return
}
