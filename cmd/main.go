package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"kursach/internal/app/api"
	"kursach/internal/app/config"
	"kursach/internal/app/service/accounting"
	"kursach/internal/app/service/balance"
	"kursach/internal/app/service/balance_holder"
	"kursach/internal/app/service/pagination"
	"kursach/internal/app/service/user_transaction"
	"kursach/internal/app/storage"
)

func balanceStorage(storage *storage.PostgresStorage) balance.Storage {
	return storage
}

func userTransactionStorage(storage *storage.PostgresStorage) user_transaction.Storage {
	return storage
}

func balanceHolderStorage(storage *storage.PostgresStorage) balance_holder.Storage {
	return storage
}

func serviceHistoryStorage(storage *storage.PostgresStorage) accounting.Storage {
	return storage
}

func paginationStorage(storage *storage.PostgresStorage) pagination.Storage {
	return storage
}

// main godoc
func main() {
	fx.New(
		fx.Provide(
			context.Background,
			config.NewConfig,
			api.NewApi,
			gin.Default,
			storage.NewPostgresStorage,
			balanceStorage,
			balance.NewService,
			userTransactionStorage,
			user_transaction.NewService,
			balanceHolderStorage,
			balance_holder.NewService,
			serviceHistoryStorage,
			accounting.NewService,
			paginationStorage,
			pagination.NewService,
		),
		fx.Invoke(api.StartHook),
	).Run()
}
