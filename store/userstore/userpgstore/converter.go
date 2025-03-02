package userpgstore

import (
	"github.com/TakumaKurosawa/sqlc-common-transaction/model"
	"github.com/TakumaKurosawa/sqlc-common-transaction/pkg/db"
)

// Converts from db.User to model.User
func toModelUser(dbUser db.User) model.User {
	return model.User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}
}

// Converts from db.User slice to model.User slice
func toModelUserList(dbUsers []db.User) []model.User {
	users := make([]model.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = toModelUser(dbUser)
	}
	return users
}
