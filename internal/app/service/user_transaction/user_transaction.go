package user_transaction

import (
	"context"
	"fmt"
	"kursach/internal/app/models"
	"log"
)

const (
	UpBalance     = "Пополнение баланса"
	WithdrawMoney = "Снятие денег"
)

type Storage interface {
	MakeUserTransaction(context.Context, models.UserTransaction) error
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) UpBalanceTransaction(ctx context.Context, userID uint64, cost float32) {
	err := s.storage.MakeUserTransaction(ctx, models.NewUserTransaction(userID, UpBalance, cost))
	if err != nil {
		log.Print(err)
	}
}

func (s *Service) MakeServiceTransaction(ctx context.Context, userID uint64, serviceName string, cost float32) {
	err := s.storage.MakeUserTransaction(ctx, models.NewUserTransaction(userID, serviceName, cost))
	if err != nil {
		log.Print(err)
	}
}

func (s *Service) MakeTransferTransaction(ctx context.Context, senderID uint64, recipientName string, cost float32) {
	description := fmt.Sprintf("Перевод пользователю %s", recipientName)
	err := s.storage.MakeUserTransaction(ctx, models.NewUserTransaction(senderID, description, cost))
	if err != nil {
		log.Print(err)
	}
}

func (s *Service) MakeReceiveTransaction(ctx context.Context, recipientID uint64, senderName string, cost float32) {
	description := fmt.Sprintf("Перевод от пользователя %s", senderName)
	err := s.storage.MakeUserTransaction(ctx, models.NewUserTransaction(recipientID, description, cost))
	if err != nil {
		log.Print(err)
	}
}

func (s *Service) MakeWithdrawMoneyTransaction(ctx context.Context, userID uint64, cost float32) {
	err := s.storage.MakeUserTransaction(ctx, models.NewUserTransaction(userID, WithdrawMoney, cost))
	if err != nil {
		log.Print(err)
	}
}
