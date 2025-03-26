package utils

import (
	"errors"

	"github.com/xeipuuv/gojsonschema"
)

func ValidateJSONSchema(schema interface{}, body []byte) error {
	schemaLoader := gojsonschema.NewGoLoader(schema)
	bodyLoader := gojsonschema.NewBytesLoader(body)
	result, err := gojsonschema.Validate(schemaLoader, bodyLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return errors.New(result.Errors()[0].String())
	}
	return nil
}
