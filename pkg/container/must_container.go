package container

type MustContainer struct {
	container *Container
}

func (c *MustContainer) TokenSingleton(dependencies map[string]any) *MustContainer {
	err := c.container.TokenSingleton(dependencies)
	if err != nil {
		panic(err)
	}

	return c
}

func (c *MustContainer) Singleton(resolvers ...any) *MustContainer {
	err := c.container.Singleton(resolvers...)
	if err != nil {
		panic(err)
	}

	return c
}

func (c *MustContainer) Token(dependencies map[string]any) *MustContainer {
	err := c.container.Token(dependencies)
	if err != nil {
		panic(err)
	}

	return c
}

func (c *MustContainer) Dependencies(resolvers ...any) *MustContainer  {
	err := c.container.Transient(resolvers...)
	if err != nil {
		panic(err)
	}

	return c
}


