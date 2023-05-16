package storage

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"kursach/internal/app/models"
	balance_service "kursach/internal/app/service/balance"
)

const (
	BalanceTable = "user_balance"
	UserID       = "user_id"
	Balance      = "balance"
)

func (p *PostgresStorage) UpsetUserBalance(ctx context.Context, balance models.UserBalance) error {
	onConflict := fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %s = %s.%s + $2", UserID, Balance, BalanceTable, Balance)
	query, args, err := sq.
		Insert(BalanceTable).
		Columns(UserID, Balance).
		Values(
			balance.UserID,
			balance.Balance,
		).
		Suffix(onConflict).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, query, args...)
	return err

}

func (p *PostgresStorage) GetUserBalance(ctx context.Context, userID uint64) (float32, error) {
	var balance float32
	query, args, err := sq.
		Select(Balance).
		From(BalanceTable).
		Where(sq.Eq{UserID: userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return balance, err
	}

	err = p.db.GetContext(ctx, &balance, query, args...)
	return balance, err
}

func (p *PostgresStorage) RemoveMoney(ctx context.Context, userID uint64, cost float32) error {

	query, args, err := sq.
		Update(BalanceTable).
		Set(Balance, sq.Expr(fmt.Sprintf("%s-%f", Balance, cost))).
		Where(sq.Eq{UserID: userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, query, args...)
	return err
}

func (p *PostgresStorage) TransferMoney(ctx context.Context, senderID uint64, recipientID uint64, cost float32) error {
	tx, err := p.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	if err != nil {
		return err
	}

	defer tx.Rollback()

	balance, err := p.GetUserBalance(ctx, senderID)
	if err != nil {
		return err
	}
	if balance < cost {
		return balance_service.ErrInsufficientBalance
	}

	err = p.RemoveMoney(ctx, senderID, cost)
	if err != nil {
		return err
	}

	err = p.UpsetUserBalance(ctx, models.UserBalance{UserID: recipientID, Balance: cost})
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *PostgresStorage) WithdrawMoney(ctx context.Context, userID uint64, cost float32) error {
	tx, err := p.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	if err != nil {
		return err
	}

	defer tx.Rollback()

	balance, err := p.GetUserBalance(ctx, userID)
	if err != nil {
		return err
	}
	if balance < cost {
		return balance_service.ErrInsufficientBalance
	}

	err = p.RemoveMoney(ctx, userID, cost)
	if err != nil {
		return err
	}

	return tx.Commit()
}
