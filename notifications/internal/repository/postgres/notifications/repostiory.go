package notifications

import (
	"context"
	"fmt"
	"route256/libs/logger"
	tx "route256/libs/txmanager/postgres"
	"route256/notifications/internal/domain"
	"route256/notifications/internal/repository/schema"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
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

var (
	EpochStartTime = time.Time{}
)

func (r *Repository) ListNotifications(
	ctx context.Context,
	userId int64,
	tsFrom time.Time,
	tsTill time.Time,
) ([]domain.Notification, error) {
	query := sq.
		Select(ColumnRecordId, ColumnUserId, ColumnMessage, ColumnCreatedAt).
		From(TableName).
		Where(sq.Eq{ColumnUserId: userId})

	if tsFrom != EpochStartTime {
		query = query.Where(sq.GtOrEq{ColumnCreatedAt: tsFrom})
	}

	if tsTill != EpochStartTime {
		query = query.Where(sq.LtOrEq{ColumnCreatedAt: tsTill})
	}

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query ListNotifcations: %s", err)
	}

	logger.Debugf("SQL: %s", queryRaw)
	logger.Debugf("SQL: %+v", queryArgs)

	var result []schema.Notification
	err = pgxscan.Select(ctx, r.GetQuerier(ctx), &result, queryRaw, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("exec query ListNotifcations: %s", err)
	}

	return ConvNotificationSchemaDomain(result), nil
}

func (r *Repository) SaveNotification(ctx context.Context, recordId int64, userId int64, message string) error {
	query := sq.
		Insert(TableName).
		Columns(ColumnRecordId, ColumnUserId, ColumnMessage, ColumnCreatedAt).
		Values(recordId, userId, message, time.Now())

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build query SaveNotification: %s", err)
	}

	logger.Debugf("SQL: %s", queryRaw)
	logger.Debugf("SQL: %+v", queryArgs)

	_, err = r.GetQuerier(ctx).Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return fmt.Errorf("exec query SaveNotification: %s", err)
	}

	return nil
}
