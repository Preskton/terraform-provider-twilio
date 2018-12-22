package twilio

import (
	"errors"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/fatih/structs"

	log "github.com/sirupsen/logrus"
)

func MarshalToTerraform(src interface{}, dest *schema.ResourceData) error {
	t := reflect.TypeOf(src)

	if t == nil && structs.IsStruct(src) {
		return errors.New("Source object cannot be nil and must be a struct")
	}

	log.Debug("Got a good struct")

	for _, field := range structs.Fields(src) {
		tag := field.Tag("terraform")

		if tag == "" {
			continue
		}

		options := strings.Split(tag, ",")
		tfname := options[0]

		value := field.Value()

		log.Debugf("Setting %s to %s", tfname, value)

		dest.Set(tfname, value)
	}

	return nil
}
