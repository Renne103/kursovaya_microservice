package accounting

import (
	"context"
	"encoding/csv"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"kursach/internal/app/models"
	"log"
	"os"
)

type Storage interface {
	AddServiceTransaction(context.Context, models.ServiceTransaction) error
	GetReport(context.Context, string) ([]models.AccountingReport, error)
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) AddServiceTransaction(ctx context.Context, operation models.HolderOperation) {
	err := s.storage.AddServiceTransaction(ctx, models.NewServiceTransaction(operation))
	if err != nil {
		log.Print(err)
	}
}

func (s *Service) GetReport(ctx context.Context, date string) ([]models.AccountingReport, error) {
	return s.storage.GetReport(ctx, date)
}

func (s *Service) MakeReportLink(ctx context.Context, date string) (string, error) {
	link := fmt.Sprintf("http://localhost:8080/api/v1/accounting/report/%s-report", date)
	report, err := s.storage.GetReport(ctx, date)
	if err != nil {
		return link, err
	}
	file, err := os.Create(fmt.Sprintf("%s-report.csv", date))
	if err != nil {
		return link, err
	}
	defer file.Close()
	t := transform.NewWriter(file, unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewEncoder())
	w := csv.NewWriter(t)
	for _, row := range report {
		w.Write(row.ToCSV())
	}
	w.Flush()
	return link, err
}
