package pagination

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	ObjectTime = "time"
	ObjectCost = "cost"
)

var (
	ErrInvalidParam = errors.New("invalid params")
)

type Cursor struct {
	TransactionID uint64 `json:"transaction_id"`
	FieldValue    string `json:"field_value"`
	Forward       bool   `json:"forward"`
}

type Report struct {
	TransactionID uint64    `db:"transaction_id" json:"transaction_id"`
	OperationTime time.Time `db:"operation_time" json:"operation_time"`
	Description   string    `db:"description" json:"description"`
	Cost          float32   `db:"cost" json:"cost"`
}

func MakeCursorHash(objectSort string, report Report, forward bool) (string, error) {
	var cursor = Cursor{
		TransactionID: report.TransactionID,
		Forward:       forward,
	}

	if objectSort == ObjectCost {
		cursor.FieldValue = fmt.Sprintf("%.2f", report.Cost)
	} else {
		cursor.FieldValue = report.OperationTime.String()
	}

	return CodeBase64Cursor(cursor)
}

func CodeBase64Cursor(cursor interface{}) (string, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	err := json.NewEncoder(encoder).Encode(cursor)
	if err != nil {
		return "", err
	}
	encoder.Close()
	return buf.String(), nil
}

func DecodeBase64Cursor(cursor interface{}, codeCursor string) error {
	return json.NewDecoder(base64.NewDecoder(base64.StdEncoding, strings.NewReader(codeCursor))).Decode(cursor)
}

type Storage interface {
	DownTimeFirstRequest(context.Context, uint64) ([]Report, error)
	DownTimeNext(context.Context, uint64, uint64, string) ([]Report, error)
	DownTimePrev(context.Context, uint64, uint64, string) ([]Report, error)
	UpTimeFirstRequest(context.Context, uint64) ([]Report, error)
	UpTimeNext(context.Context, uint64, uint64, string) ([]Report, error)
	UpTimePrev(context.Context, uint64, uint64, string) ([]Report, error)
	DownCostFirstRequest(context.Context, uint64) ([]Report, error)
	DownCostNext(context.Context, uint64, uint64, float32) ([]Report, error)
	DownCostPrev(context.Context, uint64, uint64, float32) ([]Report, error)
	UpCostFirstRequest(context.Context, uint64) ([]Report, error)
	UpCostNext(context.Context, uint64, uint64, float32) ([]Report, error)
	UpCostPrev(context.Context, uint64, uint64, float32) ([]Report, error)
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func validateURLParameters(increase, userId string) (uint64, error) {
	if increase != "true" && increase != "false" {
		return 0, ErrInvalidParam
	}

	userID, err := strconv.Atoi(userId)
	if err != nil {
		return 0, ErrInvalidParam
	}
	return uint64(userID), err
}

func ChooseCursor(pReport, nReport Report, cursorStr string, reportLen int, objectSort string) (string, string, error) {
	var prev string
	var next string
	if cursorStr == "" && reportLen == 5 {
		hash, err := MakeCursorHash(objectSort, nReport, true)
		if err != nil {
			return "", "", err
		}
		next = hash
		return prev, next, err
	}

	if cursorStr == "" && reportLen < 5 {
		return "", "", nil
	}

	var cursor Cursor
	_ = DecodeBase64Cursor(&cursor, cursorStr)

	if cursor.Forward && reportLen < 5 {
		hash, err := MakeCursorHash(objectSort, pReport, false)
		if err != nil {
			return prev, next, err
		}
		prev = hash
		return prev, next, err
	}

	if !cursor.Forward && reportLen < 5 {
		hash, err := MakeCursorHash(objectSort, nReport, true)
		if err != nil {
			return prev, next, err
		}
		next = hash
		return prev, next, err
	}

	hash, err := MakeCursorHash(objectSort, nReport, true)
	if err != nil {
		return prev, next, err
	}
	next = hash

	hash, err = MakeCursorHash(objectSort, pReport, false)
	if err != nil {
		return prev, next, err
	}
	prev = hash

	return prev, next, nil
}

func (s *Service) PaginationCost(ctx context.Context, increase, userId, cursorStr string) ([]Report, error) {
	userID, err := validateURLParameters(increase, userId)
	if err != nil {
		return nil, err
	}

	if cursorStr == "" {
		if increase == "false" {
			return s.storage.DownCostFirstRequest(ctx, userID)
		}
		if increase == "true" {
			return s.storage.UpCostFirstRequest(ctx, userID)
		}
	}

	var cursor Cursor
	err = DecodeBase64Cursor(&cursor, cursorStr)
	if err != nil {
		return nil, ErrInvalidParam
	}

	cost64, err := strconv.ParseFloat(cursor.FieldValue, 32)
	if err != nil {
		return nil, ErrInvalidParam
	}
	cost := float32(cost64)

	if increase == "false" {
		if cursor.Forward {
			return s.storage.DownCostNext(ctx, cursor.TransactionID, userID, cost)
		}
		return s.storage.DownCostPrev(ctx, cursor.TransactionID, userID, cost)
	} else {
		if cursor.Forward {
			return s.storage.UpCostNext(ctx, cursor.TransactionID, userID, cost)
		}
		return s.storage.UpCostPrev(ctx, cursor.TransactionID, userID, cost)
	}
}

func (s *Service) PaginationTime(ctx context.Context, increase, userId, cursorStr string) ([]Report, error) {
	userID, err := validateURLParameters(increase, userId)
	if err != nil {
		return nil, err
	}

	if cursorStr == "" {
		if increase == "false" {
			return s.storage.DownTimeFirstRequest(ctx, userID)
		}
		return s.storage.UpTimeFirstRequest(ctx, userID)
	}

	var cursor Cursor
	err = DecodeBase64Cursor(&cursor, cursorStr)
	if err != nil {
		return nil, ErrInvalidParam
	}

	if increase == "false" {
		if cursor.Forward {
			return s.storage.DownTimeNext(ctx, cursor.TransactionID, userID, cursor.FieldValue)
		}
		return s.storage.DownTimePrev(ctx, cursor.TransactionID, userID, cursor.FieldValue)
	} else {
		if cursor.Forward {
			return s.storage.UpTimeNext(ctx, cursor.TransactionID, userID, cursor.FieldValue)
		}
		return s.storage.UpTimePrev(ctx, cursor.TransactionID, userID, cursor.FieldValue)
	}
}

func (s *Service) DownTimeFirstRequest(ctx context.Context, userID uint64) ([]Report, error) {
	return s.storage.DownTimeFirstRequest(ctx, userID)
}

func (s *Service) DownTimeNext(ctx context.Context, transactionID uint64, userID uint64, atTime string) ([]Report, error) {
	return s.storage.DownTimeNext(ctx, transactionID, userID, atTime)
}

func (s *Service) DownTimePrev(ctx context.Context, transactionID uint64, userID uint64, atTime string) ([]Report, error) {
	return s.storage.DownTimePrev(ctx, transactionID, userID, atTime)
}

func (s *Service) UpTimeFirstRequest(ctx context.Context, userID uint64) ([]Report, error) {
	return s.storage.UpTimeFirstRequest(ctx, userID)
}

func (s *Service) UpTimeNext(ctx context.Context, transactionID uint64, userID uint64, atTime string) ([]Report, error) {
	return s.storage.UpTimeNext(ctx, transactionID, userID, atTime)
}

func (s *Service) UpTimePrev(ctx context.Context, transactionID uint64, userID uint64, atTime string) ([]Report, error) {
	return s.storage.UpTimePrev(ctx, transactionID, userID, atTime)
}

func (s *Service) DownCostFirstRequest(ctx context.Context, userID uint64) ([]Report, error) {
	return s.storage.DownCostFirstRequest(ctx, userID)
}

func (s *Service) DownCostNext(ctx context.Context, transactionID uint64, userID uint64, cost float32) ([]Report, error) {
	return s.storage.DownCostNext(ctx, transactionID, userID, cost)
}

func (s *Service) DownCostPrev(ctx context.Context, transactionID uint64, userID uint64, cost float32) ([]Report, error) {
	return s.storage.DownCostPrev(ctx, transactionID, userID, cost)
}

func (s *Service) UpCostFirstRequest(ctx context.Context, userID uint64) ([]Report, error) {
	return s.storage.DownCostFirstRequest(ctx, userID)
}

func (s *Service) UpCostNext(ctx context.Context, transactionID uint64, userID uint64, cost float32) ([]Report, error) {
	return s.storage.UpCostNext(ctx, transactionID, userID, cost)
}

func (s *Service) UpCostPrev(ctx context.Context, transactionID uint64, userID uint64, cost float32) ([]Report, error) {
	return s.storage.UpCostPrev(ctx, transactionID, userID, cost)
}
