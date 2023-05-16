package storage

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"kursach/internal/app/models"
	"kursach/internal/app/service/balance_holder"
)

const (
	BalanceHolderTable = "balance_holder"
	OrderID            = "order_id"
	ServiceID          = "service_id"
	ServiceName        = "service_name"
)

func (p *PostgresStorage) ReserveMoney(ctx context.Context, operation models.HolderOperation) error {
	tx, err := p.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	if err != nil {
		return err
	}

	defer tx.Rollback()

	balance, err := p.GetUserBalance(ctx, operation.UserID)
	if err != nil {
		return err
	}
	if balance < operation.Cost {
		return balance_holder.ErrInsufficientBalance
	}

	if err := p.RemoveMoney(ctx, operation.UserID, operation.Cost); err != nil {
		return err
	}

	query, args, err := sq.
		Insert(BalanceHolderTable).
		Columns(OrderID, ServiceID, UserID, ServiceName, Cost).
		Values(
			operation.OrderID,
			operation.ServiceID,
			operation.UserID,
			operation.ServiceName,
			operation.Cost,
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *PostgresStorage) GetHolderOperationByOrderID(ctx context.Context, orderID uint64) (models.HolderOperation, error) {
	var operation models.HolderOperation

	query, args, err := sq.
		Select(OrderID, ServiceID, UserID, ServiceName, Cost).
		From(BalanceHolderTable).
		Where(sq.Eq{OrderID: orderID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return operation, err
	}

	err = p.db.GetContext(ctx, &operation, query, args...)

	return operation, err
}

func (p *PostgresStorage) DelHolderOperationByOrderID(ctx context.Context, orderID uint64) error {

	query, args, err := sq.
		Delete(BalanceHolderTable).
		Where(sq.Eq{OrderID: orderID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = p.db.Exec(query, args...)

	return err
}

func (p *PostgresStorage) RevenueConfirmation(ctx context.Context, orderID uint64) error {
	err := p.DelHolderOperationByOrderID(ctx, orderID)
	return err
}

func (p *PostgresStorage) ReturnMoney(ctx context.Context, orderID uint64) error {
	tx, err := p.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	operation, err := p.GetHolderOperationByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	if err := p.DelHolderOperationByOrderID(ctx, orderID); err != nil {
		return err
	}

	if err := p.UpsetUserBalance(ctx, models.UserBalance{UserID: operation.UserID, Balance: operation.Cost}); err != nil {
		return err
	}

	return tx.Commit()
}
