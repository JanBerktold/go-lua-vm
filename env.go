package lua

type Environment struct {
	m      map[string]interface{}
	parent *Environment
}

func (env *Environment) SearchValue(key string) interface{} {
	for k, v := range env.m {
		if k == key {
			return v
		}
	}
	if env.parent != nil {
		return env.parent.SearchValue(key)
	}
	return nil
}

func (env *Environment) SetValue(key string, value interface{}) {
	env.m[key] = value
}

func NewEnv(parent *Environment) *Environment {
	return &Environment{
		make(map[string]interface{}),
		parent,
	}
}
