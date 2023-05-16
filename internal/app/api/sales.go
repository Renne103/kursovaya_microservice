package api

import (
	"github.com/gin-gonic/gin"
	"kursach/internal/app/models"
	"kursach/internal/app/service/balance_holder"
	"net/http"
)

type Reserve struct {
	OrderID     uint64  `json:"order_id"`
	ServiceID   uint64  `json:"service_id"`
	UserID      uint64  `json:"user_id"`
	ServiceName string  `json:"service_name"`
	Cost        float32 `json:"cost"`
}

// ReserveMoney godoc
// @Summary Резервирование средств
// @Schemes
// @Description Резервирование средств с основного баланса на отдельном счете
// @Tags sales
// @Accept json
// @Param data body Reserve true "Входные параметры"
// @Success 200
// @Router /sales/reserve_money [post]
func (api *Api) ReserveMoney(ctx *gin.Context) {
	var reserve Reserve

	if err := ctx.BindJSON(&reserve); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := api.balanceHolder.ReserveMoney(ctx, models.HolderOperation{
		OrderID:     reserve.OrderID,
		ServiceID:   reserve.ServiceID,
		UserID:      reserve.UserID,
		ServiceName: reserve.ServiceName,
		Cost:        reserve.Cost,
	})

	if err == balance_holder.ErrInvalidCost {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err == balance_holder.ErrInsufficientBalance {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
}

type Confirmation struct {
	OrderID     uint64  `db:"order_id" json:"order_id"`
	ServiceID   uint64  `db:"service_id" json:"service_id"`
	UserID      uint64  `db:"user_id" json:"user_id"`
	ServiceName string  `json:"service_name"`
	Cost        float32 `db:"cost" json:"cost"`
}

// RevenueConfirmation godoc
// @Summary Метод признания выручки
// @Schemes
// @Description Списывает из резерва деньги, добавляет данные в отчет для бухгалтерии
// @Tags sales
// @Accept json
// @Param data body Confirmation true "Входные параметры"
// @Success 200
// @Router /sales/revenue_confirmation [post]
func (api *Api) RevenueConfirmation(ctx *gin.Context) {
	var confirmation Confirmation

	if err := ctx.BindJSON(&confirmation); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := api.balanceHolder.RevenueConfirmation(ctx, confirmation.OrderID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	api.userTransaction.MakeServiceTransaction(ctx, confirmation.UserID, confirmation.ServiceName, confirmation.Cost)
	api.accounting.AddServiceTransaction(ctx, models.HolderOperation{
		OrderID:     confirmation.OrderID,
		ServiceID:   confirmation.ServiceID,
		UserID:      confirmation.UserID,
		ServiceName: confirmation.ServiceName,
		Cost:        confirmation.Cost,
	})

	ctx.AbortWithStatus(http.StatusOK)
}

type Refund struct {
	OrderID uint64 `db:"order_id" json:"order_id"`
}

// ReturnMoney godoc
// @Summary Возвращение средств
// @Schemes
// @Description Разрезервирование средств и возвращение их на баланс пользователя при отмене операции
// @Tags sales
// @Accept json
// @Param data body Refund true "Входные параметры"
// @Success 200
// @Router /sales/return_money [post]
func (api *Api) ReturnMoney(ctx *gin.Context) {
	var refund Refund

	if err := ctx.BindJSON(&refund); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := api.balanceHolder.ReturnMoney(ctx, refund.OrderID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
}
