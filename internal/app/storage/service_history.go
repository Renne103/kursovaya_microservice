package storage

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"kursach/internal/app/models"
)

const (
	ServiceHistoryTable = "service_history"
	SubHistory          = "sub_history"
	TotalCost           = "total_cost"
)

func (p *PostgresStorage) AddServiceTransaction(ctx context.Context, transaction models.ServiceTransaction) error {
	query, args, err := sq.
		Insert(ServiceHistoryTable).
		Columns(OrderID, ServiceID, UserID, OperationTime, ServiceName, Cost).
		Values(
			transaction.OrderID,
			transaction.ServiceID,
			transaction.UserID,
			transaction.OperationTime,
			transaction.ServiceName,
			transaction.Cost,
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, query, args...)
	return err
}

func (p *PostgresStorage) GetReport(ctx context.Context, date string) ([]models.AccountingReport, error) {
	var reportRows []models.AccountingReport

	query := fmt.Sprintf(`select service_name, sum(cost) as total_cost 
from service_history 
where cast(operation_time as varchar(10)) like '%s-%%' 
group by service_name`, date)

	err := p.db.SelectContext(ctx, &reportRows, query)
	return reportRows, err
}
