package balance_holder

import (
	"context"
	"errors"
	"kursach/internal/app/models"
)

var (
	ErrInvalidCost         = errors.New("invalid cost")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type Storage interface {
	ReserveMoney(context.Context, models.HolderOperation) error
	RevenueConfirmation(context.Context, uint64) error
	ReturnMoney(context.Context, uint64) error
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) ReserveMoney(ctx context.Context, operation models.HolderOperation) error {
	if operation.Cost < 0 {
		return ErrInvalidCost
	}

	return s.storage.ReserveMoney(ctx, operation)
}

func (s *Service) RevenueConfirmation(ctx context.Context, orderID uint64) error {
	return s.storage.RevenueConfirmation(ctx, orderID)
}

func (s *Service) ReturnMoney(ctx context.Context, orderID uint64) error {
	return s.storage.ReturnMoney(ctx, orderID)
}
