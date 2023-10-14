package envelope

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

const (
	envTag = "envelope"
)

type property = string

const (
	required property = "required"
)

type Field struct {
	Key      string
	Required bool
}

func decodeStruct(refType reflect.Type) (map[string]any, error) {
	fields := map[string]any{}

	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		tagValue := field.Tag.Get(envTag)

		switch field.Type.Kind() {
		case reflect.Struct:
			decoded, err := decodeStruct(field.Type)
			if err != nil {
				return nil, err
			}
			fields[field.Name] = decoded
		default:
			fieldProps := getFieldProperties(tagValue)
			value, ok := os.LookupEnv(fieldProps.Key)
			if !ok && fieldProps.Required {
				return nil, fmt.Errorf(`missing a required field "%s"`, fieldProps.Key)
			}
			fields[field.Name] = value
		}
	}
	return fields, nil
}

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

func getFieldProperties(tagValue string) *Field {
	field := new(Field)
	values := strings.Split(tagValue, ",")

	field.Key = values[0]
	for i := 1; i < len(values); i++ {
		switch values[i] {
		case required:
			field.Required = true
		}
	}
	return field
}
