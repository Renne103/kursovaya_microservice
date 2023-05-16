package storage

import (
	"context"
	"fmt"
	"kursach/internal/app/service/pagination"
)

const (
	PageLimit = 5
)

func (p *PostgresStorage) DownTimeFirstRequest(ctx context.Context, userID uint64) ([]pagination.Report, error) {
	var reportList []pagination.Report
	query := fmt.Sprintf(`SELECT transaction_id, operation_time, description, cost FROM transaction_history WHERE user_id = %d ORDER BY operation_time ASC Limit 5`, userID)

	err := p.db.SelectContext(ctx, &reportList, query)
	return reportList, err
}

func (p *PostgresStorage) DownTimeNext(ctx context.Context, transactionID uint64, userID uint64, atTime string) ([]pagination.Report, error) {
	var reportList []pagination.Report
	query := fmt.Sprintf(`SELECT transaction_id, operation_time, description, cost FROM transaction_history WHERE (operation_time, transaction_id) < ('%s' :: timestamp, %d) AND user_id = %d GROUP BY operation_time DESC LIMIT 5`, atTime, transactionID, userID)

	err := p.db.SelectContext(ctx, &reportList, query)
	return reportList, err
}

func (p *PostgresStorage) DownTimePrev(ctx context.Context, transactionID uint64, userID uint64, atTime string) ([]pagination.Report, error) {
	var reportList []pagination.Report
	query := fmt.Sprintf(`SELECT transaction_id, operation_time, description, cost FROM transaction_history WHERE (operation_time, transaction_id) > ('%s' :: timestamp, %d) AND user_id = %d GROUP BY operation_time DESC LIMIT 5`, atTime, transactionID, userID)

	err := p.db.SelectContext(ctx, &reportList, query)
	return reportList, err
}

func (p *PostgresStorage) UpTimeFirstRequest(ctx context.Context, userID uint64) ([]pagination.Report, error) {
	var reportList []pagination.Report
	query := fmt.Sprintf(`SELECT transaction_id, operation_time, description, cost FROM transaction_history WHERE user_id = %d ORDER BY operation_time DESC Limit 5`, userID)

	err := p.db.SelectContext(ctx, &reportList, query)
	return reportList, err
}

func (p *PostgresStorage) UpTimeNext(ctx context.Context, transactionID uint64, userID uint64, atTime string) ([]pagination.Report, error) {
	var reportList []pagination.Report
	query := fmt.Sprintf(`SELECT transaction_id, operation_time, description, cost FROM transaction_history WHERE (operation_time, transaction_id) > ('%s' :: timestamp, %d) and user_id = %d ORDER BY operation_time ASC LIMIT 5`, atTime, transactionID, userID)

	err := p.db.SelectContext(ctx, &reportList, query)
	return reportList, err
}

func (p *PostgresStorage) UpTimePrev(ctx context.Context, transactionID uint64, userID uint64, atTime string) ([]pagination.Report, error) {
	var reportList []pagination.Report
	query := fmt.Sprintf(`SELECT transaction_id, operation_time, description, cost FROM transaction_history WHERE (operation_time, transaction_id) < ('%s' :: timestamp, %d) AND user_id = %d GROUP BY operation_time ASC LIMIT 5`, atTime, transactionID, userID)

	err := p.db.SelectContext(ctx, &reportList, query)
	return reportList, err
}

func (p *PostgresStorage) DownCostFirstRequest(ctx context.Context, userID uint64) ([]pagination.Report, error) {
	var reportList []pagination.Report
	query := fmt.Sprintf(`SELECT transaction_id, operation_time, description, cost FROM transaction_history WHERE user_id = %d ORDER BY cost DESC LIMIT 5;`, userID)

	err := p.db.SelectContext(ctx, &reportList, query)
	return reportList, err
}

func (p *PostgresStorage) DownCostNext(ctx context.Context, transactionID uint64, userID uint64, cost float32) ([]pagination.Report, error) {
	var reportList []pagination.Report
	query := fmt.Sprintf(`SELECT transaction_id, operation_time, description, cost FROM transaction_history WHERE (cost, transaction_id) < (%.2f, %d) and user_id = %d GROUP BY cost DESC LIMIT 5`, cost, transactionID, userID)

	err := p.db.SelectContext(ctx, &reportList, query)
	return reportList, err
}

func (p *PostgresStorage) DownCostPrev(ctx context.Context, transactionID uint64, userID uint64, cost float32) ([]pagination.Report, error) {
	var reportList []pagination.Report
	query := fmt.Sprintf(`SELECT transaction_id, operation_time, description, cost FROM transaction_history WHERE (cost, transaction_id) > (%.2f, %d) and user_id = %d GROUP BY cost DESC LIMIT 5`, cost, transactionID, userID)

	err := p.db.SelectContext(ctx, &reportList, query)
	return reportList, err
}

func (p *PostgresStorage) UpCostFirstRequest(ctx context.Context, userID uint64) ([]pagination.Report, error) {
	var reportList []pagination.Report
	query := fmt.Sprintf(`SELECT transaction_id, operation_time, description, cost FROM transaction_history WHERE user_id = %d ORDER BY cost ASC LIMIT 5;`, userID)

	err := p.db.SelectContext(ctx, &reportList, query)
	return reportList, err
}

func (p *PostgresStorage) UpCostNext(ctx context.Context, transactionID uint64, userID uint64, cost float32) ([]pagination.Report, error) {
	var reportList []pagination.Report
	query := fmt.Sprintf(`SELECT transaction_id, operation_time, description, cost FROM transaction_history WHERE (cost, transaction_id) > (%.2f, %d) and user_id = %d ORDER BY cost ASC LIMIT 5;`, cost, transactionID, userID)

	err := p.db.SelectContext(ctx, &reportList, query)
	return reportList, err
}

func (p *PostgresStorage) UpCostPrev(ctx context.Context, transactionID uint64, userID uint64, cost float32) ([]pagination.Report, error) {
	var reportList []pagination.Report
	query := fmt.Sprintf(`SELECT transaction_id, operation_time, description, cost FROM transaction_history WHERE (cost, transaction_id) < (%.2f, %d) and user_id = %d ORDER BY cost ASC LIMIT 5;`, cost, transactionID, userID)

	err := p.db.SelectContext(ctx, &reportList, query)
	return reportList, err
}
