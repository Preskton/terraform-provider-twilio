package mapper

import (
	"bytes"
	"fmt"

	"github.com/fatih/structs"
	"github.com/hashicorp/terraform/helper/hashcode"
)

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
