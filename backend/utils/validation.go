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
    var users []map[string]interface{}

    if err := json.Unmarshal(rawJSON, &users); err != nil {
        var user map[string]interface{}
        if err := json.Unmarshal(rawJSON, &user); err != nil {
            fmt.Println("Error unmarshalling JSON:", err)
            return false
        }
        users = append(users, user)
    }

    for _, user := range users {
        for _, field := range requiredFields {
            if value, exists := user[field]; !exists || value == "" {
                fmt.Printf("Missing or empty required field: %s\n", field)
                return false
            }
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
