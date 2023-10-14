package envelope

import "strings"

const sep = ","

type property = string

const (
	required property = "required"
)

type Field struct {
	Key      string
	Required bool
}

func getFieldProperties(tagValue string) *Field {
	field := new(Field)
	values := strings.Split(tagValue, sep)

	field.Key = values[0]
	for i := 1; i < len(values); i++ {
		switch values[i] {
		case required:
			field.Required = true
		}
	}
	return field
}
