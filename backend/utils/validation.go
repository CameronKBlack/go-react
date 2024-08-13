package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unicode"
)

func NameValidCheck(name string) bool {
	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) {
			return false
		}
	}

	return true
}

func CheckRequiredFields(rawJSON []byte, requiredFields []string) bool {
    var requestBody map[string]interface{}
    if err := json.Unmarshal(rawJSON, &requestBody); err != nil {
        fmt.Println("Error unmarshalling JSON:", err)
        return false
    }

    for _, field := range requiredFields {
        if _, exists := requestBody[field]; !exists || requestBody[field] == "" {
            fmt.Printf("Missing or empty required field: %s\n", field)
            return false
        }
    }
    return true
}

func GetRequiredFields(typ reflect.Type) []string {
	var fields []string
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			fields = append(fields, jsonTag)
		}
	}
	return fields
}
