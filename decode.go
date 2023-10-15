package envelope

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

const (
	structFieldTag = "envelope"
)

func Decode[T any](environmentModel *T) error {
	refType := reflect.TypeOf(*environmentModel)

	decodedMap, errMsg := decodeStruct(refType)
	if errMsg != "" {
		return errors.New(errMsg)
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

func decodeStruct(refType reflect.Type) (map[string]any, string) {
	fields := map[string]any{}
	errorAggregate := []string{}

	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		tagValue, tagFound := field.Tag.Lookup(structFieldTag)

		typeKind := field.Type.Kind()
		if typeKind == reflect.Struct {
			decoded, errMsg := decodeStruct(field.Type)
			if errMsg != "" {
				errorAggregate = append(errorAggregate, errMsg)
				continue
			}

			if field.Anonymous {
				for k, v := range decoded {
					fields[k] = v
				}
			} else {
				fields[field.Name] = decoded
			}

			continue
		}

		if !tagFound {
			continue
		}

		fieldProps := getFieldProperties(tagValue)
		value, foundEnv := os.LookupEnv(fieldProps.Key)
		if !foundEnv && fieldProps.Required {
			errMsg := fmt.Sprintf(`missing a required field "%s"`, fieldProps.Key)
			errorAggregate = append(errorAggregate, errMsg)
			continue
		}

		convertedValue, err := convert(typeKind, value)
		if err != nil {
			errMsg := fmt.Sprintf(`error converting value from "%s" field into %s"`, fieldProps.Key, typeKind)
			errorAggregate = append(errorAggregate, errMsg)
			continue
		}

		fields[field.Name] = convertedValue
	}

	if len(errorAggregate) > 0 {
		return nil, strings.Join(errorAggregate, "; ")
	}

	return fields, ""
}
