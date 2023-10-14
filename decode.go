package envelope

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

const (
	envTag = "envelope"
)

func Decode[T any](environmentModel *T) error {
	refType := reflect.TypeOf(*environmentModel)

	decodedMap, err := decodeStruct(refType)
	if err != nil {
		return err
	}

	marshal, err := json.Marshal(decodedMap)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(marshal, &environmentModel); err != nil {
		return err
	}
	return nil
}

func decodeStruct(refType reflect.Type) (map[string]any, error) {
	fields := map[string]any{}

	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		tagValue := field.Tag.Get(envTag)

		typeKind := field.Type.Kind()
		if typeKind == reflect.Struct {
			decoded, err := decodeStruct(field.Type)
			if err != nil {
				return nil, err
			}
			fields[field.Name] = decoded
			continue
		}

		fieldProps := getFieldProperties(tagValue)
		value, ok := os.LookupEnv(fieldProps.Key)
		if !ok && fieldProps.Required {
			return nil, fmt.Errorf(`missing a required field "%s"`, fieldProps.Key)
		}

		convertedValue, err := convert(typeKind, value)
		if err != nil {
			return nil, err
		}

		fields[field.Name] = convertedValue
	}
	return fields, nil
}
