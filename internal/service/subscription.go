package service

import (
	"context"
	"subServ/internal/domain"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo domain.SubscriptionRepository
	log  domain.Logger
}

func NewSubscriptionService(repo domain.SubscriptionRepository, log domain.Logger) *SubscriptionService {
	return &SubscriptionService{repo: repo, log: log}
}

func (s *SubscriptionService) Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	if sub.Price <= 0 || sub.ServiceName == "" {
		return domain.Subscription{}, domain.ErrInvalidInput
	}
	if sub.EndDate != nil && sub.EndDate.Before(sub.StartDate) {
		return domain.Subscription{}, domain.ErrInvalidInput
	}

	result, err := s.repo.Create(ctx, sub)
	if err != nil {
		s.log.Error("failed to create subscription", "err", err)
		return domain.Subscription{}, err
	}

	s.log.Info("subscription created", "id", result.ID)
	return result, nil
}

func (s *SubscriptionService) GetByID(ctx context.Context, id uuid.UUID) (domain.Subscription, error) {
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get subscription", "id", id, "err", err)
		return domain.Subscription{}, err
	}
	return result, nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.log.Error("failed to delete subscription", "id", id, "err", err)
		return err
	}

	s.log.Info("subscription deleted", "id", id)
	return nil
}

func (s *SubscriptionService) Update(ctx context.Context, sub domain.Subscription) error {
	if sub.Price <= 0 || sub.ServiceName == "" {
		return domain.ErrInvalidInput
	}
	if sub.EndDate != nil && sub.EndDate.Before(sub.StartDate) {
		return domain.ErrInvalidInput
	}

	if err := s.repo.Update(ctx, sub); err != nil {
		s.log.Error("failed to update subscription", "id", sub.ID, "err", err)
		return err
	}

	s.log.Info("subscription updated", "id", sub.ID)
	return nil
}

func (s *SubscriptionService) List(ctx context.Context, filter domain.SubscriptionFilter) (domain.SubscriptionList, error) {
	result, err := s.repo.List(ctx, filter)
	if err != nil {
		s.log.Error("failed to list subscriptions", "err", err)
		return domain.SubscriptionList{}, err
	}
	return result, nil
}

func (s *SubscriptionService) TotalCost(ctx context.Context, filter domain.TotalCostFilter) (int, error) {
	result, err := s.repo.TotalCost(ctx, filter)
	if err != nil {
		s.log.Error("failed to calculate total cost", "err", err)
		return 0, err
	}
	return result, nil
}
