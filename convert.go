package envelope

import (
	"fmt"
	"reflect"
	"strconv"
)

func convert(typeKind reflect.Kind, value string) (converted any, err error) {
	switch typeKind {
	case reflect.String:
		converted = value
	case reflect.Int:
		converted, err = strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
	case reflect.Float64:
		converted, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
	case reflect.Bool:
		converted, err = strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported type")
	}
	return
}
