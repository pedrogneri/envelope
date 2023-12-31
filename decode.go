package envelope

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

const (
	structFieldTag = "envelope"
)

// Decode decodes environment variables into a struct. The struct fields must be tagged with "envelope".
// The function returns an error aggregator that aggregates all the errors on decoding process.
// The function modifies the struct passed as a pointer.
func Decode[T any](v *T) error {
	refType := reflect.TypeOf(*v)

	decodedEnv, errMsg := decodeEnv(refType)
	if errMsg != "" {
		return errors.New(errMsg)
	}

	parsedStruct, ok := decodedEnv.(T)
	if !ok {
		return errors.New("failed to parse decoded environment to struct")
	}

	*v = parsedStruct

	return nil
}

func decodeEnv(refType reflect.Type) (any, string) {
	errorAggregate := []string{}
	envModelElem := reflect.New(refType).Elem()

	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		tagValue, tagFound := field.Tag.Lookup(structFieldTag)

		elemField := envModelElem.FieldByName(field.Name)

		typeKind := field.Type.Kind()
		if typeKind == reflect.Struct {
			decodedEnv, errMsg := decodeEnv(field.Type)
			if errMsg != "" {
				errorAggregate = append(errorAggregate, errMsg)
				continue
			}

			errMsg = setFieldValue(field, elemField, decodedEnv)
			if errMsg != "" {
				errorAggregate = append(errorAggregate, errMsg)
				continue
			}
		}

		if !tagFound {
			continue
		}

		tagProps := getTagProperties(tagValue)
		value, foundEnv := os.LookupEnv(tagProps.key)
		if !foundEnv {
			if tagProps.isRequired {
				errMsg := fmt.Sprintf(`missing a required field "%s"`, tagProps.key)
				errorAggregate = append(errorAggregate, errMsg)
				continue
			}

			if tagProps.defaultValue == "" {
				continue
			}

			value = tagProps.defaultValue
		}

		convertedValue, err := convert(typeKind, value)
		if err != nil {
			errMsg := fmt.Sprintf(`error converting value from "%s" field into %s`, tagProps.key, typeKind)
			errorAggregate = append(errorAggregate, errMsg)
			continue
		}

		errMsg := setFieldValue(field, elemField, convertedValue)
		if errMsg != "" {
			errorAggregate = append(errorAggregate, errMsg)
		}
	}

	if len(errorAggregate) > 0 {
		return nil, strings.Join(errorAggregate, "; ")
	}

	return envModelElem.Interface(), ""
}

func setFieldValue(refField reflect.StructField, refValue reflect.Value, setValue any) (errMsg string) {
	if !refValue.IsValid() {
		errMsg = fmt.Sprintf(`field "%s" was invalid`, refField.Name)
		return
	}

	if !refValue.CanSet() {
		errMsg = fmt.Sprintf(`field "%s" can't be set`, refField.Name)
		return
	}

	refValue.Set(reflect.ValueOf(setValue))
	return
}
