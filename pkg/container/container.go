package container

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/4strodev/wiring/pkg/internal/collections/graph"
	"github.com/4strodev/wiring/pkg/resolver"
)

type Container struct {
	graph.Graph[resolver.DependencyResolver[any]]
	typeIndex  map[reflect.Type]*graph.Node[resolver.DependencyResolver[any]]
	tokenIndex map[string]*graph.Node[resolver.DependencyResolver[any]]
	connected  bool
}

// Retuns a new container and sets a default dependency that allows
// to inject this container as a parameter on the resolvers
func New() *Container {
	container := &Container{
		Graph:      graph.NewGraph[resolver.DependencyResolver[any]](),
		typeIndex:  make(map[reflect.Type]*graph.Node[resolver.DependencyResolver[any]]),
		tokenIndex: make(map[string]*graph.Node[resolver.DependencyResolver[any]]),
	}

	// Allow resolvers to inject container
	container.AddDependency(func() *Container {
		return container
	})

	return container
}

func (c *Container) AddDependency(res any) error {
	c.connected = false
	if !resolver.IsValid(res) {
		return errors.New("invalid resolver")
	}

	builder := resolver.DependencyResolver[any]{
		Resolver: res,
	}

	resType := builder.Type()
	_, exists := c.typeIndex[resType]
	if exists {
		return fmt.Errorf("dependency for this type already exists: %v", resType)
	}

	node := graph.NewNode(builder)
	c.Add(node)
	c.typeIndex[resType] = node

	return nil
}

func (c *Container) AddTokenDependency(token string, res any) error {
	c.connected = false
	if !resolver.IsValid(res) {
		return errors.New("invalid resolver")
	}

	builder := resolver.DependencyResolver[any]{
		Resolver: res,
	}

	_, exists := c.tokenIndex[token]
	if exists {
		return fmt.Errorf("dependency for token type already exists: %s", token)
	}

	node := graph.NewNode(builder)
	c.Add(node)
	c.tokenIndex[token] = node

	return nil
}

func (c *Container) AddDependencies(resolvers ...any) error {
	c.connected = false
	for _, res := range resolvers {
		if !resolver.IsValid(res) {
			return errors.New("invalid resolver")
		}

		builder := resolver.DependencyResolver[any]{
			Resolver: res,
		}

		resType := builder.Type()
		_, exists := c.typeIndex[resType]
		if exists {
			return fmt.Errorf("dependency for this type already exists: %v", resType)
		}

		node := graph.NewNode(builder)
		c.Add(node)
		c.typeIndex[resType] = node
	}

	return nil
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
		return nil, fmt.Errorf("dependency for %v not found", t)
	}

	return node, nil
}

// setConnections stablishes connections between nodes and
// look for circular dependencies
func (c *Container) setConnections() error {
	for node := range c.Graph.GetNodes() {
		resType := reflect.TypeOf(node.Val.Resolver)
		for i := 0; i < resType.NumIn(); i++ {
			dependencyType := resType.In(i)

			if node.Val.Type() == dependencyType {
				return fmt.Errorf("circular depndency found: %v", []*graph.Node[resolver.DependencyResolver[any]]{
					node,
				})
			}

			dependencyNode, err := c.getNodeFor(dependencyType)
			if err != nil {
				return err
			}

			if dependencyNode.IsConnectedWith(node) {
				return fmt.Errorf("circular dependency found: %v",
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
		return fmt.Errorf("circular dependency found: %v", cicle)
	}

	c.connected = true
	return nil
}

func (c Container) resolve(t reflect.Type) (resolvedValue reflect.Value, err error) {
	node, ok := c.typeIndex[t]
	if !ok {
		err = fmt.Errorf("dependency not found for type %v", t)
		return
	}

	dependencyResolver := node.Val
	inputTypes := dependencyResolver.Input()
	inputArgs := []reflect.Value{}

	for _, inputType := range inputTypes {
		arg, err := c.resolve(inputType)
		if err != nil {
			return resolvedValue, err
		}

		inputArgs = append(inputArgs, arg)
	}

	return resolver.Execute(dependencyResolver, inputArgs)
}

func (c Container) resolveToken(token string) (resolvedValue reflect.Value, err error) {
	node, ok := c.tokenIndex[token]
	if !ok {
		err = fmt.Errorf("dependency not found for token '%s'", token)
		return
	}

	dependencyResolver := node.Val
	inputTypes := dependencyResolver.Input()
	inputArgs := []reflect.Value{}

	for _, inputType := range inputTypes {
		arg, err := c.resolve(inputType)
		if err != nil {
			return resolvedValue, err
		}

		inputArgs = append(inputArgs, arg)
	}

	return resolver.Execute(dependencyResolver, inputArgs)
}
