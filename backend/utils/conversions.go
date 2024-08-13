package utils

import (
	"encoding/json"
	"fmt"
	"go-react/backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Using this instead of MarshalJSON() due to issues when adding to DB where it exluded password
func ConvertUserFormat(users []models.User) []map[string]interface{} {
	var userResponses []map[string]interface{}

	for _, u := range users {
		userResponse := map[string]interface{}{
			"username": u.Username,
			"first_name": u.FirstName,
			"last_name": u.LastName,
			"date_of_birth": u.DateOfBirth,
			"email_address": u.EmailAddress,
		}
		userResponses = append(userResponses, userResponse)
	}
	return userResponses
}

func ConvertFromBSONToUserSlice(bs []primitive.M) ([]models.User, error) {
	var userList []models.User
	for _, doc := range bs {
        userJSON, err := json.Marshal(doc)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal user data: %w", err)
        }

        var usr models.User
        if err := json.Unmarshal(userJSON, &usr); err != nil {
            return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
        }
		userList = append(userList, usr)
	}
	return userList, nil
}