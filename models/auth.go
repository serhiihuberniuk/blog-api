package models

type LoginPayload struct {
	Email    string `bson:"email, omitempty"`
	Password string `bson:"password, omitempty"`
}
