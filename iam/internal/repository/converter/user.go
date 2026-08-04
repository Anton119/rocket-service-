package converter

import (
	"github.com/Anton119/rocket-service-/iam/internal/model"
	"github.com/Anton119/rocket-service-/iam/internal/repository/record"
)

// UserToRecord переводит доменную модель в запись PostgreSQL.
func UserToRecord(u model.User) record.User {
	return record.User{
		UUID:         u.UUID,
		Login:        u.Login,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

// UserFromRecord переводит запись PostgreSQL в доменную модель.
func UserFromRecord(u record.User) model.User {
	return model.User{
		UUID:         u.UUID,
		Login:        u.Login,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
