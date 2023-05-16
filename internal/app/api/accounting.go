package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type ReportRequest struct {
	Date string `json:"date"`
}

type ReportLink struct {
	Link string `json:"link"`
}

// GetReportLink godoc
// @Summary Получение ссылки с отчётом за месяц
// @Schemes
// @Description Метод по году и месяцу возвращает ссылку на скачивание отчёта
// @Tags accounting
// @Accept json
// @Param data body ReportRequest true "Входные параметры"
// @Produce json
// @Success 200 {object} ReportLink
// @Router /accounting/report_link [post]
func (api *Api) GetReportLink(ctx *gin.Context) {
	var request ReportRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	link, err := api.accounting.MakeReportLink(ctx, request.Date)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, ReportLink{Link: link})
}

type FileName struct {
	FileName string `uri:"file_name"`
}

func (api *Api) SendReport(ctx *gin.Context) {
	var fileName FileName
	if err := ctx.ShouldBindUri(&fileName); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	fileBytes, _ := ioutil.ReadFile(fmt.Sprintf("%s.csv", fileName.FileName))
	ctx.Writer.Header().Set("content-disposition", fmt.Sprintf("attachment; filename=%s.csv", fileName.FileName))
	ctx.Writer.Header().Set("Content-Type", "application/octet-stream")
	ctx.Writer.Write(fileBytes)
	ctx.AbortWithStatus(http.StatusOK)
}
