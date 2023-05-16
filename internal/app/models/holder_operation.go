package models

type HolderOperation struct {
	OrderID     uint64  `db:"order_id" json:"order_id"`
	ServiceID   uint64  `db:"service_id" json:"service_id"`
	UserID      uint64  `db:"user_id" json:"user_id"`
	ServiceName string  `db:"service_name" json:"service_name"`
	Cost        float32 `db:"cost" json:"cost"`
}
