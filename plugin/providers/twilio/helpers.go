package twilio

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/fatih/structs"

	log "github.com/sirupsen/logrus"
)

// MarshalToTerraform takes a source struct (`src`), a destination Terraform *ResourceData (`dest`), and a Terraform schema map[string]*Schema
// and then marshals the source data into the destination data given a `terraform` tag present on the fields in the source struct.
func MarshalToTerraform(src interface{}, dest *schema.ResourceData, sm map[string]*schema.Schema) error {
	if src == nil || !structs.IsStruct(src) {
		return fmt.Errorf("src cannot be nil and must be a struct")
	}

	if dest == nil {
		return fmt.Errorf("dest cannot be null")
	}

	if sm == nil {
		return fmt.Errorf("sm cannot be null")
	}

	mv, err := MapStructByTag(src, "terraform")

	if err != nil {
		return fmt.Errorf("Failed to map values: %s", err)
	}

	for terraformFieldName, sourceValue := range mv {
		switch sm[terraformFieldName].Type {
		case schema.TypeSet:
			nestedSet := schema.NewSet(SimpleHashcode, nil)

			if !structs.IsStruct(sourceValue) {
				return fmt.Errorf("Terraform field `%s` is a Set, but target value is not a struct", terraformFieldName)
			}

			mappedValue, err := MapStructByTag(sourceValue, "terraform")
			nestedSet.Add(mappedValue)

			if err != nil {
				log.Errorf("Unable to marshal %s to terraform struct map", terraformFieldName)
			}

			err = dest.Set(fmt.Sprintf("%s", terraformFieldName), nestedSet)

			if err != nil {
				return fmt.Errorf("Setting `%s` failed: %s", terraformFieldName, err)
			}
		default:
			dest.Set(terraformFieldName, sourceValue)
		}
	}

	return nil
}

// MapStructByTag takes a struct and target tag name present on fields in that struct,
// then converts it into a map[string]interface{}. The target tag should be of the format `myTag:"destinationFieldName"`,
// where `destinationFieldName` is a valid map[string] key.
// TODO Nested structs - currently only handle one level.
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

// SimpleHashcode calculates a simple integer hashcode by iterating over all the fields/keys in a map,
// concating the values in buffer, then calculating the hashcode of that buffer.
func SimpleHashcode(v interface{}) int {
	var buf bytes.Buffer

	if structs.IsStruct(v) {
		for _, value := range structs.Map(v) {
			if value != nil {
				buf.WriteString(value.(string))
			} else {
				buf.WriteString("nil")
			}
		}
	} else {
		switch v.(type) {
		case map[string]interface{}:
			for _, value := range v.(map[string]interface{}) {
				if value != nil {
					buf.WriteString(fmt.Sprintf("%v", value))
				} else {
					buf.WriteString("nil")
				}
			}
		default:
			return -1
		}
	}

	result := hashcode.String(buf.String())

	return result
}
