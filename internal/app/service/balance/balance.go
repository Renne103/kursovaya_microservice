package balance

import (
	"context"
	"errors"
	"kursach/internal/app/models"
	"log"
)

var (
	ErrInvalidCost         = errors.New("invalid balance")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type Storage interface {
	UpsetUserBalance(context.Context, models.UserBalance) error
	GetUserBalance(context.Context, uint64) (float32, error)
	TransferMoney(context.Context, uint64, uint64, float32) error
	WithdrawMoney(context.Context, uint64, float32) error
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) UpsetUserBalance(ctx context.Context, userBalance models.UserBalance) error {
	if userBalance.Balance < 0 {
		return ErrInvalidCost
	}
	err := s.storage.UpsetUserBalance(ctx, userBalance)
	if err != nil {
		log.Printf("UpsetUserBalance: can't upset user balance: %v", err)
	}
	return err
}

func (s *Service) GetUserBalance(ctx context.Context, userID uint64) (float32, error) {
	return s.storage.GetUserBalance(ctx, userID)
}

func (s *Service) TransferMoney(ctx context.Context, senderID uint64, recipientID uint64, cost float32) error {
	if cost < 0 {
		return ErrInvalidCost
	}

	return s.storage.TransferMoney(ctx, senderID, recipientID, cost)
}

func (s *Service) WithdrawMoney(ctx context.Context, userID uint64, cost float32) error {
	if cost < 0 {
		return ErrInvalidCost
	}

	return s.storage.WithdrawMoney(ctx, userID, cost)
}
