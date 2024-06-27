package echo

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/alphafast/asmt-fw/libs/domain/noti"
	"github.com/alphafast/asmt-fw/libs/utils/kafka/producer"
	logUtil "github.com/alphafast/asmt-fw/libs/utils/log"
)

type NotiHandlerDeps struct {
	NotiUseCase noti.NotiUseCase
	Producer    *producer.Producer
}

type NotiHandlerOps struct {
	NotiRequestTopicName string
}

type NotificationHandler struct {
	notiUseCase noti.NotiUseCase
	producer    *producer.Producer

	opts NotiHandlerOps
}

func NewRequestNotificationHandler(e *echo.Echo, d NotiHandlerDeps, opt NotiHandlerOps) {
	handler := &NotificationHandler{
		notiUseCase: d.NotiUseCase,
		producer:    d.Producer,
		opts:        opt,
	}

	e.POST("/notification/asynchronous", handler.RequestNotificationAsynchronous)
	e.POST("/notification", handler.RequestNotification)
}

func (nh *NotificationHandler) RequestNotification(c echo.Context) error {
	var req noti.RequestNotificationBySourceEvent
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx, logger := logUtil.SetupZeroLogWithCtx(c.Request().Context())
	notiReqs, err := nh.notiUseCase.BuildNotiRequestBySourceEvent(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	notiResults, err := nh.notiUseCase.Notify(ctx, notiReqs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	for _, res := range notiResults {
		if err := nh.notiUseCase.SaveNotifyResult(ctx, res); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	logger.Info().Msgf("[NotificationHandler.RequestNotification] %d notification(s) processed complete", len(notiResults))

	return c.JSON(http.StatusOK, notiResults)
}

func (nh *NotificationHandler) RequestNotificationAsynchronous(c echo.Context) error {
	var req noti.RequestNotificationBySourceEvent
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx, logger := logUtil.SetupZeroLogWithCtx(c.Request().Context())
	notiReqs, err := nh.notiUseCase.BuildNotiRequestBySourceEvent(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	events := []producer.Message{}
	for _, notiReq := range notiReqs {
		events = append(events, producer.Message{
			Topic:   nh.opts.NotiRequestTopicName,
			Message: notiReq,
		})
	}
	err = nh.producer.Produce(ctx, events...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "[NotificationHandler.RequestNotificationAsynchronous] error occurred").Error())
	}

	logger.Info().Msgf("[NotificationHandler.RequestNotificationAsynchronous] %d notification(s) request accepted", len(notiReqs))

	return c.JSON(http.StatusAccepted, notiReqs)
}
