package schema

import "github.com/invopop/jsonschema"

func Reflect[T any]() ([]byte, error) {
	var t T

	reflector := new(jsonschema.Reflector)
	reflector.Anonymous = true
	schema := reflector.Reflect(t)

	return schema.MarshalJSON()
}
