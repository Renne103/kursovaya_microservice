package api

import (
	"github.com/gin-gonic/gin"
	"kursach/internal/app/models"
	balance_service "kursach/internal/app/service/balance"
	"net/http"
)

type UserID struct {
	ID uint64 `uri:"user_id"`
}

type BalanceReplenishment struct {
	UserID        uint64  `json:"user_id"`
	Replenishment float32 `json:"replenishment"`
}

type Balance struct {
	Balance float32 `json:"balance"`
}

// UpBalance godoc
// @Summary Пополнение или инициализация баланса
// @Schemes
// @Description Пополняет баланс пользователя или создаёт его при первом пополнении
// @Tags balance
// @Accept json
// @Param data body BalanceReplenishment true "Входные параметры"
// @Success 200
// @Router /balance/up_balance [post]
func (api *Api) UpBalance(ctx *gin.Context) {
	var replenishment BalanceReplenishment

	if err := ctx.BindJSON(&replenishment); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := api.balance.UpsetUserBalance(ctx, models.UserBalance{UserID: replenishment.UserID, Balance: replenishment.Replenishment})
	if err == balance_service.ErrInvalidCost {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	api.userTransaction.UpBalanceTransaction(ctx, replenishment.UserID, replenishment.Replenishment)

	ctx.AbortWithStatus(http.StatusOK)
}

// GetBalance godoc
// @Summary Получение баланса пользователя
// @Schemes
// @Description Возвращает баланс пользователя по его id
// @Tags balance
// @Accept json
// @Param data body UserID true "Входные параметры"
// @Produce json
// @Success 200 {object} Balance
// @Router /balance/get_balance [post]
func (api *Api) GetBalance(ctx *gin.Context) {
	var userID UserID

	if err := ctx.BindJSON(&userID); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	balance, err := api.balance.GetUserBalance(ctx, userID.ID)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, Balance{Balance: balance})

}

type Transfer struct {
	SenderID      uint64  `json:"sender_id"`
	SenderName    string  `json:"sender_name"`
	RecipientID   uint64  `json:"recipient_id"`
	RecipientName string  `json:"recipient_name"`
	Cost          float32 `json:"cost"`
}

// TransferMoney godoc
// @Summary Перевод денег
// @Schemes
// @Description Метод перевода денег от пользователя к пользователю
// @Tags balance
// @Accept json
// @Param data body Transfer true "Входные параметры"
// @Produce json
// @Success 200
// @Router /balance/transfer_money [post]
func (api *Api) TransferMoney(ctx *gin.Context) {
	var transfer Transfer

	if err := ctx.BindJSON(&transfer); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := api.balance.TransferMoney(ctx, transfer.SenderID, transfer.RecipientID, transfer.Cost)
	if err == balance_service.ErrInvalidCost {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err == balance_service.ErrInsufficientBalance {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	api.userTransaction.MakeTransferTransaction(ctx, transfer.SenderID, transfer.RecipientName, transfer.Cost)
	api.userTransaction.MakeReceiveTransaction(ctx, transfer.RecipientID, transfer.SenderName, transfer.Cost)

	ctx.AbortWithStatus(http.StatusOK)
}

type DebitingMoney struct {
	UserID    uint64  `json:"user_id"`
	DebitCost float32 `json:"debit_cost"`
}

// WithdrawMoney godoc
// @Summary Снятие денег со счёта
// @Schemes
// @Description Метод снимает указанное количество средств со счёта
// @Tags balance
// @Accept json
// @Param data body DebitingMoney true "Входные параметры"
// @Produce json
// @Success 200
// @Router /balance/withdraw_money [post]
func (api *Api) WithdrawMoney(ctx *gin.Context) {
	var debit DebitingMoney

	if err := ctx.BindJSON(&debit); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := api.balance.WithdrawMoney(ctx, debit.UserID, debit.DebitCost)
	if err == balance_service.ErrInvalidCost {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err == balance_service.ErrInsufficientBalance {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	api.userTransaction.MakeWithdrawMoneyTransaction(ctx, debit.UserID, debit.DebitCost)

	ctx.AbortWithStatus(http.StatusOK)
}
