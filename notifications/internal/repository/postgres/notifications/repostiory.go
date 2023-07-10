package notifications

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/libs/tracing"
	"route256/notifications/internal/domain"
	"route256/notifications/internal/repository/schema"

	tx "route256/libs/txmanager/postgres"

	sq "github.com/Masterminds/squirrel"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
)

type Repository struct {
	tx.Manager
}

func NewNotificationsRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{tx.Manager{Pool: pool}}
}

const (
	TableName = "notifications"

	ColumnRecordId = "record_id"
	ColumnUserId   = "user_id"
	ColumnMessage  = "message"

	ColumnCreatedAt = "created_at"
)

func (r *Repository) ListNotifications(ctx context.Context, userId int64) ([]domain.Notification, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "loms/orders/list_order")
	defer span.Finish()

	query := sq.
		Select(ColumnRecordId, ColumnUserId, ColumnMessage, ColumnCreatedAt).
		From(TableName).
		Where(sq.Eq{ColumnUserId: userId})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		err = tracing.MarkSpanWithError(ctx, fmt.Errorf("build query ListOrder: %s", err))
		return nil, err
	}

	logger.Debugf("SQL: %s", queryRaw)
	logger.Debugf("SQL: %+v", queryArgs)

	var result []schema.Notification
	err = pgxscan.Select(ctx, r.GetQuerier(ctx), &result, queryRaw, queryArgs...)
	if err != nil {
		err = tracing.MarkSpanWithError(ctx, fmt.Errorf("exec query ListStocks: %s", err))
		return nil, err
	}

	return ConvNotificationSchemaDomain(result), nil
}
