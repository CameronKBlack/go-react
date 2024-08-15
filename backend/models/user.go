package models

type User struct {
	Username     string `json:"username" bson:"username"`
	FirstName    string `json:"first_name" bson:"first_name"`
	LastName     string `json:"last_name" bson:"last_name"`
	DateOfBirth  string `json:"date_of_birth" bson:"date_of_birth"`
	EmailAddress string `json:"email_address" bson:"email_address"`
	Password     string `json:"password" bson:"password"`
}
