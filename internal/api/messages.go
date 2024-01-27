package api

import (
	"gamelight/internal/eventbus"
	"gamelight/internal/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	Engine   *gin.Engine
	Producer *eventbus.AMQPProducer
}

func NewMessageHandler(engine *gin.Engine, producer *eventbus.AMQPProducer) *MessageHandler {
	return &MessageHandler{
		Engine:   engine,
		Producer: producer,
	}
}

func (handler *MessageHandler) RegisterRoutes() {
	handler.Engine.POST("message", handler.handlePostMessage)
}

func (handler *MessageHandler) handlePostMessage(ctx *gin.Context) {
	var msg types.Message

	if err := ctx.ShouldBindJSON(&msg); err != nil {

		ctx.AbortWithStatusJSON(http.StatusBadRequest, &types.Response{Error: err.Error()})
		return
	}

	msg.Time = time.Now()

	err := handler.Producer.Produce(ctx, &msg)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadGateway, &types.Response{Error: err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, &types.Response{Error: "Created", Data: msg})
}
