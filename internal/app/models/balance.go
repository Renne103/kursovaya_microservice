package models

type UserBalance struct {
	UserID  uint64  `db:"user_id" json:"user_id"`
	Balance float32 `db:"balance" json:"balance"`
}
