package mapper

import (
	"errors"
	"strings"

	"github.com/fatih/structs"
)

// MapStructByTag takes a struct and target tag name present on fields in that struct,
// then converts it into a map[string]interface{}. The target tag should be of the format `myTag:"destinationFieldName"`,
// where `destinationFieldName` is a valid map[string] key.
func MapStructByTag(src interface{}, tagName string) (map[string]interface{}, error) {
	if src == nil || !structs.IsStruct(src) {
		return nil, errors.New("Source cannot be nil and must be a struct")
	}

	result := make(map[string]interface{})

	for _, sourceField := range structs.Fields(src) {
		tag := sourceField.Tag(tagName)
		if tag == "" {
			continue
		}

		options := strings.Split(tag, ",")
		if len(options) < 1 {
			continue
		}

		destinationFieldName := options[0]
		sourceValue := sourceField.Value()

		result[destinationFieldName] = sourceValue
	}

	return result, nil
}
