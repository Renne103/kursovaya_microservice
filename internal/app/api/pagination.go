package api

import (
	"github.com/gin-gonic/gin"
	"kursach/internal/app/service/pagination"
	"net/http"
)

type Page struct {
	List       []pagination.Report `json:"list"`
	NextCursor string              `json:"next_cursor,omitempty"`
	PrevCursor string              `json:"prev_cursor,omitempty"`
}

type SortQuery struct {
	ObjectSort string `uri:"object_sort"`
	Increase   bool   `uri:"increase"`
	UserID     uint64 `uri:"user_id"`
	Cursor     string `uri:"cursor"`
}

// SortByTime godoc
// @Summary Сортировка транзакций по времени с пагинацией
// @Schemes
// @Description Сортирует транзакции пользователя по времени и с использованием пагинации
// @Tags pagination
// @Success 200
// @Param increase query string  true  "сортировка по возрастанию — true, по убыванию — false"
// @Param user_id query integer  true  "id пользователя"
// @Param cursor query string false "курсор пагинации"
// @Produce json
// @Success 200 {object} Page
// @Router /user_transaction/time_sort [get]
func (api *Api) SortByTime(ctx *gin.Context) {
	increase := ctx.Query("increase")
	userID := ctx.Query("user_id")
	cursor := ctx.Query("cursor")

	report, err := api.cursorPagination.PaginationTime(ctx, increase, userID, cursor)
	if err == pagination.ErrInvalidParam {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	page := Page{List: report}
	lenReport := len(report)
	prev, next, err := pagination.ChooseCursor(page.List[0], page.List[lenReport-1], cursor, lenReport, "time")
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if prev != "" {
		page.PrevCursor = prev
	}
	if next != "" {
		page.NextCursor = next
	}

	ctx.JSON(http.StatusOK, page)
}

// SortByCost godoc
// @Summary Сортировка транзакций по стоимости c пагинацией
// @Schemes
// @Description Сортирует транзакции пользователя по стоимости и с использованием пагинации
// @Tags pagination
// @Success 200
// @Param increase query string  true  "сортировка по возрастанию — true, по убыванию — false"
// @Param user_id query integer  true  "id пользователя"
// @Param cursor query string false "курсор пагинации"
// @Produce json
// @Success 200 {object} Page
// @Router /user_transaction/cost_sort [get]
func (api *Api) SortByCost(ctx *gin.Context) {
	increase := ctx.Query("increase")
	userID := ctx.Query("user_id")
	cursor := ctx.Query("cursor")

	report, err := api.cursorPagination.PaginationCost(ctx, increase, userID, cursor)
	if err == pagination.ErrInvalidParam {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	page := Page{List: report}
	lenReport := len(report)
	prev, next, err := pagination.ChooseCursor(page.List[0], page.List[lenReport-1], cursor, lenReport, "cost")
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if prev != "" {
		page.PrevCursor = prev
	}
	if next != "" {
		page.NextCursor = next
	}

	ctx.JSON(http.StatusOK, page)
}
