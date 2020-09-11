package user

import (
	"context"
	"errors"

	"github.com/shellhub-io/shellhub/api/store"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrConflict = errors.New("conflict")

type Service interface {
	UpdateDataUser(ctx context.Context, username, email, currentPassword, newPassword, tenant string) ([]InvalidField, error)
}

type service struct {
	store store.Store
}

const (
	conflictName = "This username already exists"
	conflictEmail = "This email already exists"
)

type InvalidField struct {
	Name string
	Message string
	Kind string
}

func NewService(store store.Store) Service {
	return &service{store}
}

func (s *service) UpdateDataUser(ctx context.Context, username, email, currentPassword, newPassword, tenant string) ([]InvalidField, error) {
	var invalidFields []InvalidField
	user, err := s.store.GetUserByTenant(ctx, tenant)

	if err != nil {
		return invalidFields, err
	}
	if newPassword != "" && user.Password != currentPassword {
		return invalidFields, ErrUnauthorized
	}

	var checkName, checkEmail bool

	user, err = s.store.GetUserByUsername(ctx, username)
	if err == nil && user.TenantID != tenant {
		checkName = true
		invalidFields = append(invalidFields, InvalidField{"username", conflictName, "conflict"})
	}
	user, err = s.store.GetUserByEmail(ctx, email)
	if err == nil && user.TenantID != tenant {
		checkEmail = true
		invalidFields = append(invalidFields, InvalidField{"email", conflictEmail, "conflict"})
	}
	if checkName || checkEmail {
		return invalidFields, ErrConflict
	}
	return invalidFields, s.store.UpdateUser(ctx, username, email, currentPassword, newPassword, tenant)
}
