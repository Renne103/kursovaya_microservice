package api

import (
	"context"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"kursach/docs"
	"kursach/internal/app/config"
	"kursach/internal/app/service/accounting"
	"kursach/internal/app/service/balance"
	"kursach/internal/app/service/balance_holder"
	"kursach/internal/app/service/pagination"
	"kursach/internal/app/service/user_transaction"
)

const (
	Title    = "Balance handler API"
	BasePath = "/api/v1/"
)

type Api struct {
	router           *gin.Engine
	balance          *balance.Service
	userTransaction  *user_transaction.Service
	balanceHolder    *balance_holder.Service
	accounting       *accounting.Service
	cursorPagination *pagination.Service
}

func (api *Api) Run() {
	cfg := config.Load()
	docs.SwaggerInfo.BasePath = BasePath
	docs.SwaggerInfo.Title = Title

	api.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DefaultModelsExpandDepth(-1)))

	api.router.Use(CORSMiddleware())
	api.router.Run(cfg.Api.GetAddr())
}

func NewApi(
	router *gin.Engine,
	balance *balance.Service,
	userTransaction *user_transaction.Service,
	balanceHolder *balance_holder.Service,
	accounting *accounting.Service,
	cursorPagination *pagination.Service,
) *Api {
	svc := &Api{
		router:           router,
		balance:          balance,
		userTransaction:  userTransaction,
		balanceHolder:    balanceHolder,
		accounting:       accounting,
		cursorPagination: cursorPagination,
	}
	svc.registerRoutes()
	return svc
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, accept, origin, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func StartHook(lifecycle fx.Lifecycle, api *Api) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go api.Run()
				return nil
			},
		})
}

func (api *Api) registerRoutes() {
	base := api.router.Group(BasePath)

	balance := base.Group("/balance")
	balance.POST("/up_balance/", api.UpBalance)
	balance.POST("/get_balance/", api.GetBalance)
	balance.POST("/transfer_money", api.TransferMoney)
	balance.POST("/withdraw_money", api.WithdrawMoney)

	sales := base.Group("/sales")
	sales.POST("/reserve_money", api.ReserveMoney)
	sales.POST("/revenue_confirmation", api.RevenueConfirmation)
	sales.POST("/return_money", api.ReturnMoney)

	accounting := base.Group("/accounting")
	accounting.POST("/report_link", api.GetReportLink)
	accounting.GET("/report/:file_name", api.SendReport)

	pagination := base.Group("/user_transaction")
	pagination.GET("/time_sort", api.SortByTime)
	pagination.GET("/cost_sort", api.SortByCost)

}
