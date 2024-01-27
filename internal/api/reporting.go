package api

import (
	"gamelight/internal/storage"
	"gamelight/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportingHandler struct {
	Engine *gin.Engine
	Store  *storage.RedisStorage
}

func NewReportingHandler(engine *gin.Engine, store *storage.RedisStorage) *ReportingHandler {
	return &ReportingHandler{
		Engine: engine,
		Store:  store,
	}
}

func (handler *ReportingHandler) RegisterRoutes() {
	handler.Engine.GET("/message/list", handler.handleListMessages)
}

type listMessagesQuery struct {
	Sender   string `form:"sender" binding:"required"`
	Receiver string `form:"receiver" binding:"required"`
}

func (handler *ReportingHandler) handleListMessages(ctx *gin.Context) {
	var query listMessagesQuery

	if err := ctx.BindQuery(&query); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &types.Response{Error: err.Error()})

		return
	}

	messages, err := handler.Store.List(ctx, query.Sender, query.Receiver)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &types.Response{Error: err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, &types.Response{Data: messages})
}
