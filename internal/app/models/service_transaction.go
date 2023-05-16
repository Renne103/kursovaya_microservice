package models

import (
	"fmt"
	"time"
)

type ServiceTransaction struct {
	OrderID       uint64    `db:"order_id" json:"order_id"`
	ServiceID     uint64    `db:"service_id" json:"service_id"`
	UserID        uint64    `db:"user_id" json:"user_id"`
	OperationTime time.Time `db:"operation_time" json:"operation_time"`
	ServiceName   string    `db:"service_name" json:"service_name"`
	Cost          float32   `db:"cost" json:"cost"`
}

func NewServiceTransaction(operation HolderOperation) ServiceTransaction {
	return ServiceTransaction{
		OrderID:       operation.OrderID,
		ServiceID:     operation.ServiceID,
		UserID:        operation.UserID,
		OperationTime: time.Now(),
		ServiceName:   operation.ServiceName,
		Cost:          operation.Cost,
	}
}

type AccountingReport struct {
	ServiceName string  `db:"service_name"`
	TotalCost   float32 `db:"total_cost"`
}

func (rep AccountingReport) ToCSV() []string {
	return []string{rep.ServiceName, fmt.Sprintf("%f", rep.TotalCost)}
}
