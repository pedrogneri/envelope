package envelope

import (
	"strings"
)

const propertySep = ","
const keyValueSep = ":"

type property = string

const (
	required     property = "required"
	defaultValue property = "default"
)

type tagProps struct {
	key          string
	defaultValue string
	isRequired   bool
}

func getTagProperties(tagValue string) *tagProps {
	field := new(tagProps)
	values := strings.Split(tagValue, propertySep)

	field.key = values[0]
	for i := 1; i < len(values); i++ {
		kv := strings.Split(values[i], keyValueSep)
		isKeyValueProperty := len(kv) > 1

		if isKeyValueProperty {
			switch kv[0] {
			case defaultValue:
				field.defaultValue = kv[1]
			}
			continue
		}

		switch values[i] {
		case required:
			field.isRequired = true
		}
	}
	return field
}
