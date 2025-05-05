package container

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/4strodev/wiring/pkg/collections/graph"
	"github.com/4strodev/wiring/pkg/resolver"
)

type Container struct {
	graph.Graph[resolver.DependencyBuilder[any]]
	typeIndex map[reflect.Type]*graph.Node[resolver.DependencyBuilder[any]]
}

func New() *Container {
	return &Container{
		Graph:     graph.NewGraph[resolver.DependencyBuilder[any]](),
		typeIndex: make(map[reflect.Type]*graph.Node[resolver.DependencyBuilder[any]]),
	}
}

func (c *Container) AddDependency(res any) error {
	if !resolver.IsValid(res) {
		return errors.New("invalid resolver")
	}

	builder := resolver.DependencyBuilder[any]{
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

	return c.setConnections()
}

func (c *Container) AddDependencies(resolvers ...any) error {
	for _, res := range resolvers {
		if !resolver.IsValid(res) {
			return errors.New("invalid resolver")
		}

		builder := resolver.DependencyBuilder[any]{
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

	return c.setConnections()
}

func (c Container) getNodeFor(t reflect.Type) (*graph.Node[resolver.DependencyBuilder[any]], error) {
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
			dependencyNode, err := c.getNodeFor(dependencyType)
			if err != nil {
				return err
			}

			if dependencyNode.IsConnectedWith(node) {
				return fmt.Errorf("circular dependency found: %v",
					[]*graph.Node[resolver.DependencyBuilder[any]]{
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

	return nil
}
