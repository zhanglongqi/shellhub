package sessionmngr

import (
	"context"

	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/pkg/api/paginator"
	"github.com/shellhub-io/shellhub/pkg/models"
)

type Service interface {
	ListSessions(ctx context.Context, pagination paginator.Query) ([]models.Session, int, error)
	GetSession(ctx context.Context, uid models.UID) (*models.Session, error)
	CreateSession(ctx context.Context, session models.Session) (*models.Session, error)
	DeactivateSession(ctx context.Context, uid models.UID) error
	SetSessionAuthenticated(ctx context.Context, uid models.UID, authenticated bool) error
	RecordSession(ctx context.Context, uid models.UID, recordString string, width int, height int) error
	GetRecord(ctx context.Context, uid models.UID) ([]models.RecordedSession, int, error)
}

type service struct {
	store store.Store
}

func NewService(store store.Store) Service {
	return &service{store}
}

func (s *service) ListSessions(ctx context.Context, pagination paginator.Query) ([]models.Session, int, error) {
	return s.store.ListSessions(ctx, pagination)
}

func (s *service) GetSession(ctx context.Context, uid models.UID) (*models.Session, error) {
	return s.store.GetSession(ctx, uid)
}

func (s *service) CreateSession(ctx context.Context, session models.Session) (*models.Session, error) {
	return s.store.CreateSession(ctx, session)
}

func (s *service) DeactivateSession(ctx context.Context, uid models.UID) error {
	return s.store.DeactivateSession(ctx, uid)
}

func (s *service) SetSessionAuthenticated(ctx context.Context, uid models.UID, authenticated bool) error {
	return s.store.SetSessionAuthenticated(ctx, uid, authenticated)
}

func (s *service) RecordSession(ctx context.Context, uid models.UID, recordString string, width, height int) error {
	return s.store.RecordSession(ctx, uid, recordString, width, height)
}

func (s *service) GetRecord(ctx context.Context, uid models.UID) ([]models.RecordedSession, int, error) {
	return s.store.GetRecord(ctx, uid)
}
