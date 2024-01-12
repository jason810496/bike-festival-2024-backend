package controller

import (
	"bikefest/pkg/model"
	"github.com/gin-gonic/gin"
)

type EventController struct {
	eventService model.EventService
}

func NewEventController(eventService model.EventService) *EventController {
	return &EventController{eventService: eventService}
}

func (ctrl *EventController) GetAllEvent(c *gin.Context) {
	page, limit := RetrievePagination(c)
	events, err := ctrl.eventService.FindAll(c, page, limit)
	if err != nil {
		c.AbortWithStatusJSON(500, model.Response{
			Msg: err.Error(),
		})
	}

	c.JSON(200, model.Response{
		Data: events,
	})
}

func (ctrl *EventController) GetEventByID(c *gin.Context) {
	id := c.Param("id")
	event, err := ctrl.eventService.FindByID(c, id)
	if err != nil {
		c.AbortWithStatusJSON(500, model.Response{
			Msg: err.Error(),
		})
	}

	c.JSON(200, model.Response{
		Data: event,
	})
}

func (ctrl *EventController) CreateEvent(c *gin.Context) {
	var request model.CreateEventRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(400, model.Response{
			Msg: err.Error(),
		})
		return
	}
	newEvent := &model.Event{
		UserID:         request.UserID,
		EventID:        request.EventID,
		EventTimeStart: request.EventTimeStart,
		EventTimeEnd:   request.EventTimeEnd,
		EventDetail:    request.EventDetail,
	}
	err := ctrl.eventService.Store(c, newEvent)
	if err != nil {
		c.AbortWithStatusJSON(500, model.Response{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(200, model.Response{
		Data: newEvent,
	})
}

func (ctrl *EventController) UpdateEvent(c *gin.Context) {
	id := c.Param("id")
	var request model.CreateEventRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(400, model.Response{
			Msg: err.Error(),
		})
		return
	}
	event, err := ctrl.eventService.FindByID(c, id)
	if err != nil {
		c.AbortWithStatusJSON(500, model.Response{
			Msg: err.Error(),
		})
		return
	}
	updatedEvent := &model.Event{
		UserID:         request.UserID,
		EventID:        request.EventID,
		EventTimeStart: request.EventTimeStart,
		EventTimeEnd:   request.EventTimeEnd,
		EventDetail:    request.EventDetail,
	}
	_, err = ctrl.eventService.Update(c, updatedEvent)
	if err != nil {
		c.AbortWithStatusJSON(500, model.Response{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(200, model.Response{
		Data: event,
	})
}

func (ctrl *EventController) DeleteEvent(c *gin.Context) {
	userID := c.GetString("user_id")
	eventID := c.Param("event_id")
	_, err := ctrl.eventService.DeleteByUser(c, userID, eventID)
	if err != nil {
		c.AbortWithStatusJSON(500, model.Response{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(200, model.Response{
		Msg: "delete success",
	})
}
