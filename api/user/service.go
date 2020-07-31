package user

import (
	"context"
	"errors"

	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/pkg/models"
)

var ErrUnauthorized = errors.New("unauthorized")

type Service interface {
	UpdateDataUser(ctx context.Context, updateUser models.User, tenant string) error
}

type service struct {
	store store.Store
}

func NewService(store store.Store) Service {
	return &service{store}
}

func (s *service) UpdateDataUser(ctx context.Context, updateUser models.User, tenant string) error {
	user, _ := s.store.GetUserByTenant(ctx, tenant)
	if user != nil {
		return s.store.UpdateUser(ctx, updateUser, tenant)
	}

	return ErrUnauthorized
}
