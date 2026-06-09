package service_test

import (
	"context"
	"subServ/internal/domain"
	"subServ/internal/mocks"
	"subServ/internal/service"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupService(t *testing.T) (*service.SubscriptionService, *mocks.SubscriptionRepository) {
	mockRepo := mocks.NewSubscriptionRepository(t)
	svc := service.NewSubscriptionService(mockRepo, &mockLogger{})
	return svc, mockRepo
}

type mockLogger struct{}

func (m *mockLogger) Info(msg string, args ...any)  {}
func (m *mockLogger) Error(msg string, args ...any) {}
func (m *mockLogger) Warn(msg string, args ...any)  {}
func (m *mockLogger) Debug(msg string, args ...any) {}

func TestCreate_Success(t *testing.T) {
	svc, mockRepo := setupService(t)

	sub := domain.Subscription{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(sub, nil)

	result, err := svc.Create(context.Background(), sub)

	require.NoError(t, err)
	assert.Equal(t, sub.ServiceName, result.ServiceName)
	assert.Equal(t, sub.Price, result.Price)
}

func TestCreate_InvalidPrice(t *testing.T) {
	svc, _ := setupService(t)

	sub := domain.Subscription{
		ServiceName: "Yandex Plus",
		Price:       0,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}

	_, err := svc.Create(context.Background(), sub)

	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestCreate_EmptyServiceName(t *testing.T) {
	svc, _ := setupService(t)

	sub := domain.Subscription{
		ServiceName: "",
		Price:       400,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}

	_, err := svc.Create(context.Background(), sub)

	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestCreate_EndDateBeforeStartDate(t *testing.T) {
	svc, _ := setupService(t)

	start := time.Now()
	end := start.AddDate(0, -1, 0)

	sub := domain.Subscription{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New(),
		StartDate:   start,
		EndDate:     &end,
	}

	_, err := svc.Create(context.Background(), sub)

	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestGetByID_NotFound(t *testing.T) {
	svc, mockRepo := setupService(t)

	id := uuid.New()
	mockRepo.On("GetByID", mock.Anything, id).Return(domain.Subscription{}, domain.ErrNotFound)

	_, err := svc.GetByID(context.Background(), id)

	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestDelete_Success(t *testing.T) {
	svc, mockRepo := setupService(t)

	id := uuid.New()
	mockRepo.On("Delete", mock.Anything, id).Return(nil)

	err := svc.Delete(context.Background(), id)

	require.NoError(t, err)
}

func TestDelete_NotFound(t *testing.T) {
	svc, mockRepo := setupService(t)

	id := uuid.New()
	mockRepo.On("Delete", mock.Anything, id).Return(domain.ErrNotFound)

	err := svc.Delete(context.Background(), id)

	require.ErrorIs(t, err, domain.ErrNotFound)
}
