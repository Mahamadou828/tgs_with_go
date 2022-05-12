package stripe

import "github.com/google/uuid"

func CreateUser(email, phoneNumber, name string) (string, error) {
	return uuid.New().String(), nil
}
