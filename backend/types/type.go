package types

import "backend/models"

type GetUserByName struct {
	Message string       `json:"message"`
	Success bool         `json:"success,omitempty"`
	User    *models.User `json:"user,omitempty"`
}
