package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"subServ/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	query := `
	INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(
		&sub.ID,
		&sub.CreatedAt,
	)
	if err != nil {
		return domain.Subscription{}, err
	}
	return sub, nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Subscription, error) {
	query := `
		SELECT * FROM subscriptions WHERE id = $1
	`
	var result domain.Subscription

	err := r.db.QueryRow(ctx, query, id).Scan(
		&result.ID,
		&result.ServiceName,
		&result.Price,
		&result.UserID,
		&result.StartDate,
		&result.EndDate,
		&result.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subscription{}, domain.ErrNotFound
		}
		return domain.Subscription{}, err
	}
	return result, nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, sub domain.Subscription) error {
	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5
		WHERE id = $6
		`
	result, err := r.db.Exec(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *SubscriptionRepository) List(ctx context.Context, filter domain.SubscriptionFilter) (domain.SubscriptionList, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at FROM subscriptions`

	var conditions []string
	var args []any
	argNum := 1

	if filter.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argNum))
		args = append(args, filter.UserID)
		argNum++
	}

	if filter.ServiceName != nil {
		conditions = append(conditions, fmt.Sprintf("service_name = $%d", argNum))
		args = append(args, filter.ServiceName)
		argNum++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, filter.Limit, (filter.Page-1)*filter.Limit)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return domain.SubscriptionList{}, err
	}
	defer rows.Close()

	var items []domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
			&sub.CreatedAt,
		)
		if err != nil {
			return domain.SubscriptionList{}, err
		}
		items = append(items, sub)
	}

	if err := rows.Err(); err != nil {
		return domain.SubscriptionList{}, err
	}

	countQuery := `SELECT COUNT(*) FROM subscriptions`
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err = r.db.QueryRow(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return domain.SubscriptionList{}, err
	}

	return domain.SubscriptionList{
		Data: items,
		Meta: domain.PaginationMeta{
			Page:  filter.Page,
			Limit: filter.Limit,
			Total: total,
		},
	}, nil
}

func (r *SubscriptionRepository) TotalCost(ctx context.Context, filter domain.TotalCostFilter) (int, error) {
	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions`

	var conditions []string
	var args []any
	argNum := 1

	if filter.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argNum))
		args = append(args, filter.UserID)
		argNum++
	}

	if filter.ServiceName != nil {
		conditions = append(conditions, fmt.Sprintf("service_name = $%d", argNum))
		args = append(args, filter.ServiceName)
		argNum++
	}

	if filter.From != nil {
		conditions = append(conditions, fmt.Sprintf("start_date >= $%d", argNum))
		args = append(args, filter.From)
		argNum++
	}

	if filter.To != nil {
		conditions = append(conditions, fmt.Sprintf("start_date <= $%d", argNum))
		args = append(args, filter.To)
		argNum++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err := r.db.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
