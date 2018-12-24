package twilio

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/fatih/structs"

	log "github.com/sirupsen/logrus"
)

type SimpleFieldMappingInstruction struct {
	SourceField      string
	SourceValue      interface{}
	DestinationField string
	DestinationValue interface{}
}

func MarshalToTerraform(src interface{}, dest *schema.ResourceData, sm map[string]*schema.Schema) error {
	return marshalToTerraform(src, dest, sm)
}

func marshalToTerraform(src interface{}, dest *schema.ResourceData, sm map[string]*schema.Schema) error {
	if src == nil {
		return fmt.Errorf("Source cannot be null")
	}

	if dest == nil {
		return fmt.Errorf("Destination cannot be null")
	}

	if sm == nil {
		return fmt.Errorf("Schema cannot be null")
	}

	t := reflect.TypeOf(src)

	if t == nil || !structs.IsStruct(src) {
		return errors.New("Source cannot be nil and must be a struct")
	}

	for _, sourceField := range structs.Fields(src) {
		tag := sourceField.Tag("terraform")

		if tag == "" {
			continue
		}

		options := strings.Split(tag, ",")
		terraformFieldName := options[0]

		sourceValue := sourceField.Value()
		tfschema := sm[terraformFieldName]

		// TODO probably refactor this out

		if tfschema.Type == schema.TypeSet {
			//nestedSet := dest.Get(terraformFieldName).(*schema.Set)
			nestedSet := schema.NewSet(SimpleStructHashcode, nil)

			if !structs.IsStruct(sourceValue) {
				return fmt.Errorf("Terraform field `%s` is a Set, but field `%s` is not a struct", terraformFieldName, sourceField.Name())
			}

			mappedValue, err := MapStructByTag(sourceValue, "terraform")
			nestedSet.Add(mappedValue)

			if err != nil {
				log.Errorf("Unable to marshal %s to terraform struct map", sourceField.Name())
			}

			err = dest.Set(fmt.Sprintf("%s", terraformFieldName), nestedSet)

			if err != nil {
				return fmt.Errorf("Setting `%s` failed: %s", terraformFieldName, err)
			}
		} else if tfschema.Type == schema.TypeList {
			log.Warnf("schema.TypeList not yet implemented")
			// TODO Handle list types
		} else {
			dest.Set(terraformFieldName, sourceValue)
		}
	}

	return nil
}

func MapStructByTag(src interface{}, tagName string) (map[string]interface{}, error) {
	if src == nil || !structs.IsStruct(src) {
		return nil, errors.New("Source cannot be nil and must be a struct")
	}

	result := make(map[string]interface{})

	for _, sourceField := range structs.Fields(src) {
		tag := sourceField.Tag(tagName)

		if tag == "" {
			log.Debugf("Field %s doesn't have tag %s, skipping", sourceField.Name(), tagName)
			continue
		}

		options := strings.Split(tag, ",")

		if len(options) < 1 {
			log.Debugf("Field %s doesn't have any options, skipping", sourceField.Name())
			continue
		}

		destinationFieldName := options[0]
		sourceValue := sourceField.Value()

		result[destinationFieldName] = sourceValue
	}

	return result, nil
}

type MappedField struct {
	SourceFieldName      string
	DestinationFieldName string
	Value                interface{}
}

func SimpleStructHashcode(v interface{}) int {
	var buf bytes.Buffer

	if structs.IsStruct(v) {
		for _, value := range structs.Map(v) {
			if value != nil {
				buf.WriteString(value.(string))
			} else {
				buf.WriteString("nil")
			}
		}
	}

	result := hashcode.String(buf.String())

	log.Debugf("Calculated hashcode = %v", result)

	return result
}
