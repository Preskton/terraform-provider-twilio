package mapper

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/fatih/structs"

	log "github.com/sirupsen/logrus"
)

// TerraformIDFieldName is a special field name that can be used to identify the value that should be storedin/retrieved from
// the ID field in Terraform.
const TerraformIDFieldName = "id"

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
		if terraformFieldName == TerraformIDFieldName {
			dest.SetId(fmt.Sprintf("%s", sourceValue))
		} else {
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
	}
	return nil
}
