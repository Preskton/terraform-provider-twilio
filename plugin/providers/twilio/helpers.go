package twilio

import (
	"errors"
	"fmt"
	"reflect"
)

type TestData struct {
	SomeID     int    `terraform:"key=some_id,type=int"`
	SomeString string `terraform:"key=some_string,type=string`
}

func MarshalToTerraform(src interface{}) error {
	t := reflect.TypeOf(src)

	if t == nil {
		return errors.New("Source object cannot be nil")
	}

	for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
		field := t.Field(fieldIndex)

		tag := field.Tag

		if tfspec, exists := tag.Lookup("terraform"); exists {
			fmt.Printf("asdf %s", tfspec)
		}
	}

	return nil
}
