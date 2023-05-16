package models

import (
	"github.com/google/uuid"
	"time"
)

type UserTransaction struct {
	TransactionID uint64    `db:"transaction_id" json:"transaction_id"`
	UserID        uint64    `db:"user_id" json:"user_id"`
	OperationTime time.Time `db:"operation_time" json:"operation_time"`
	Description   string    `db:"description" json:"description"`
	Cost          float32   `db:"cost" json:"cost"`
}

func NewUserTransaction(userID uint64, description string, cost float32) UserTransaction {
	return UserTransaction{
		TransactionID: uint64(uuid.New().ID()),
		UserID:        userID,
		OperationTime: time.Now(),
		Description:   description,
		Cost:          cost,
	}
}
