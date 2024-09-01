package parking

import (
	"encoding/json"
	"fmt"
)

type LoginPayload struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

func NewLoginPayload(email string, password string) *LoginPayload {
	return &LoginPayload{
		Email:             email,
		Password:          password,
		ReturnSecureToken: true,
	}
}

func (l *LoginPayload) ToJson() ([]byte, error) {
	data, err := json.Marshal(l)
	if err != nil {
		return nil, fmt.Errorf("error marshalling login payload: %s", err)
	}
	return data, nil
}
