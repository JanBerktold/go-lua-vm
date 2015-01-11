package lua

type Environment struct {
	m map[string]interface{}
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