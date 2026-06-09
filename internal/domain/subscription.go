package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
	CreatedAt   time.Time
}

type SubscriptionService interface {
	Create(ctx context.Context, sub Subscription) (Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (Subscription, error)
	List(ctx context.Context, filter SubscriptionFilter) (SubscriptionList, error)
	Update(ctx context.Context, sub Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	TotalCost(ctx context.Context, filter TotalCostFilter) (int, error)
}

type SubscriptionList struct {
	Data []Subscription
	Meta PaginationMeta
}

type PaginationMeta struct {
	Page  int
	Limit int
	Total int
}

type SubscriptionFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	Page        int
	Limit       int
}

type TotalCostFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	From        *time.Time
	To          *time.Time
}

type SubscriptionRepository interface {
	Create(ctx context.Context, sub Subscription) (Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (Subscription, error)
	List(ctx context.Context, filter SubscriptionFilter) (SubscriptionList, error)
	Update(ctx context.Context, sub Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	TotalCost(ctx context.Context, filter TotalCostFilter) (int, error)
}
