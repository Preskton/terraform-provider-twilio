package twilio

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

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

	//destType := reflect.TypeOf(dest)

	log.Debug("Got a good struct")

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
			log.Debugf("Terraform field %s is a set", terraformFieldName)

			nestedSet := dest.Get(terraformFieldName).(*schema.Set)

			if !structs.IsStruct(sourceValue) {
				return fmt.Errorf("Terraform field %s is a Set, but field %s is not a struct", terraformFieldName, sourceField.Name())
			}

			// TODO Refactor
			mappedValues := make(map[string]interface{})

			mappedValues["lol"] = 1

			nestedSet.Add(mappedValues)
		} else if tfschema.Type == schema.TypeList {
			log.Debugf("List type found, iterating over %s", terraformFieldName)

			// TODO Handle list types
		} else {
			log.Debugf("Value type found, setting %s to %s", terraformFieldName, sourceValue)
			dest.Set(terraformFieldName, sourceValue)
		}
	}

	return nil
}
