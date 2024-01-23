package controller

import (
	"bikefest/pkg/model"
	"github.com/gin-gonic/gin"
)

type EventController struct {
	eventService model.EventService
	asynqService model.AsynqNotificationService
}

func NewEventController(eventService model.EventService, asynqService model.AsynqNotificationService) *EventController {
	return &EventController{
		eventService: eventService,
		asynqService: asynqService,
	}
}

// GetAllEvent godoc
// @Summary Get all events
// @Description Retrieves a list of all events with pagination
// @Tags Event
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page for pagination"
// @Success 200 {object} model.EventListResponse "List of events"
// @Failure 500 {object} model.Response "Internal Server Error"
// @Router /events [get]
func (ctrl *EventController) GetAllEvent(c *gin.Context) {
	page, limit := RetrievePagination(c)
	events, err := ctrl.eventService.FindAll(c, int64(page), int64(limit))
	if err != nil {
		c.AbortWithStatusJSON(500, model.Response{
			Msg: err.Error(),
		})
	}

	c.JSON(200, model.EventListResponse{
		Data: events,
	})
}

// UpdateEvent godoc
// @Summary Update an event
// @Description Updates an event by ID with new details
// @Tags Event
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Event ID"
// @Param event body model.CreateEventRequest true "Event Update Information"
// @Success 200 {object} model.EventResponse "Event successfully updated"
// @Failure 400 {object} model.Response "Bad Request - Invalid input"
// @Failure 500 {object} model.Response "Internal Server Error"
// @Router /events/{id} [put]
func (ctrl *EventController) UpdateEvent(c *gin.Context) {
	// TODO: only allow admin to update event

	id := c.Param("id")
	identity, _ := RetrieveIdentity(c, true)
	if identity.UserID != "admin" {
		c.AbortWithStatusJSON(403, model.Response{
			Msg: "permission denied",
		})
		return
	}
	_ = identity.UserID
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
		ID:             request.ID,
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

	c.JSON(200, model.EventResponse{
		Data: event,
	})
}
