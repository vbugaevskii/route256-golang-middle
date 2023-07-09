package notificationsoutbox

import (
	"context"
	"encoding/json"
	"fmt"
	"route256/libs/logger"
	tx "route256/libs/txmanager/postgres"
	"route256/loms/internal/converter"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	tx.Manager
}

func NewNotificationsOutboxRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{tx.Manager{Pool: pool}}
}

const (
	TableName = "notifications_outbox"

	ColumnRecordId  = "record_id"
	ColumnKey       = "key"
	ColumnValue     = "value"
	ColumnState     = "state"
	ColumnCreatedAt = "created_at"
)

func (r *Repository) CreateNotification(ctx context.Context, orderId int64, status domain.StatusType) (int64, error) {
	key := fmt.Sprint(orderId)
	val, err := json.Marshal(domain.Notification{
		OrderId: orderId,
		Status:  status,
	})
	if err != nil {
		return 0, err
	}

	query := sq.
		Insert(TableName).
		Columns(ColumnKey, ColumnValue, ColumnState, ColumnCreatedAt).
		Values(key, string(val), string(schema.StateWaiting), time.Now()).
		Suffix(fmt.Sprintf("RETURNING %s", ColumnRecordId))

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query CreateNotification: %s", err)
	}

	logger.Debugf("SQL: %s", queryRaw)
	logger.Debugf("SQL: %+v", queryArgs)

	row := r.GetQuerier(ctx).QueryRow(ctx, queryRaw, queryArgs...)

	var recordId int64
	if err := row.Scan(&recordId); err != nil {
		return 0, fmt.Errorf("exec query CreateNotification: %s", err)
	}

	return recordId, nil
}

func (r *Repository) SetNotificationDelivered(ctx context.Context, recordId int64) error {
	query := sq.
		Update(TableName).
		Set(ColumnState, string(schema.StateDelivered)).
		Where(sq.Eq{ColumnRecordId: recordId})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build query SetNotificationDelivered: %s", err)
	}

	// logger.Debugf("SQL: %s", queryRaw)
	// logger.Debugf("SQL: %+v", queryArgs)

	_, err = r.GetQuerier(ctx).Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return fmt.Errorf("exec query SetNotificationDelivered: %s", err)
	}

	return nil
}

func (r *Repository) ListNotificationsWaiting(ctx context.Context) ([]domain.Notification, error) {
	query := sq.
		Select(ColumnRecordId, ColumnKey, ColumnValue, ColumnState, ColumnCreatedAt).
		From(TableName).
		Where(sq.Eq{ColumnState: schema.StateWaiting}).
		OrderBy(ColumnCreatedAt)

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query ListNotificationsWaiting: %s", err)
	}

	// logger.Debugf("SQL: %s", queryRaw)
	// logger.Debugf("SQL: %+v", queryArgs)

	var notesSchema []schema.Notification
	if err := pgxscan.Select(ctx, r.GetQuerier(ctx), &notesSchema, queryRaw, queryArgs...); err != nil {
		return nil, err
	}

	notesDomain, err := converter.ConvNotificationsOutboxSchemaDomain(notesSchema)
	if err != nil {
		return nil, err
	}

	return notesDomain, nil
}

func (r *Repository) DeleteNotificationsDelivered(ctx context.Context) error {
	query := sq.
		Delete(TableName).
		Where(sq.Eq{ColumnState: string(schema.StateDelivered)})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build query DeleteNotificationsDelivered: %s", err)
	}

	// logger.Debugf("SQL: %s", queryRaw)
	// logger.Debugf("SQL: %+v", queryArgs)

	_, err = r.GetQuerier(ctx).Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return fmt.Errorf("exec query DeleteNotificationsDelivered: %s", err)
	}

	return nil
}
