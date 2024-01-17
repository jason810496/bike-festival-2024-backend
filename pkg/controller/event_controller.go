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
	events, err := ctrl.eventService.FindAll(c, page, limit)
	if err != nil {
		c.AbortWithStatusJSON(500, model.Response{
			Msg: err.Error(),
		})
	}

	c.JSON(200, model.EventListResponse{
		Data: events,
	})
}

// GetEventByID godoc
// @Summary Get event by ID
// @Description Retrieves an event using its unique ID
// @Tags Event
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} model.EventResponse "Event successfully retrieved"
// @Failure 500 {object} model.Response "Internal Server Error"
// @Router /events/{id} [get]
func (ctrl *EventController) GetEventByID(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	event, err := ctrl.eventService.FindByID(c, id)
	if err != nil {
		c.AbortWithStatusJSON(500, model.Response{
			Msg: err.Error(),
		})
	}

	if event.UserID != userID {
		c.AbortWithStatusJSON(403, model.Response{
			Msg: "permission denied",
		})
		return
	}

	c.JSON(200, model.EventResponse{
		Data: event,
	})
}

// GetUserEvent godoc
// @Summary Get User Events
// @Description Retrieves a list of events associated with a user
// @Tags Event
// @Accept json
// @Produce json
// @Security ApiKeyAuth // include this line if the endpoint is protected by an API key or other security mechanism
// @Success 200 {object} model.EventListResponse "List of events associated with the user"
// @Failure 500 {object} model.Response "Internal Server Error"
// @Router /events/user [get] // adjust the path and HTTP method according to your routing
func (ctrl *EventController) GetUserEvent(c *gin.Context) {
	userID := c.GetString("user_id")
	events, err := ctrl.eventService.FindByUserID(c, userID)
	if err != nil {
		c.AbortWithStatusJSON(500, model.Response{
			Msg: err.Error(),
		})
	}

	c.JSON(200, model.EventListResponse{
		Data: events,
	})
}

// SubscribeEvent godoc
// @Summary Subscribe to an event
// @Description Subscribes a user to an event with the provided details
// @Tags Event
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body model.CreateEventRequest true "Event Subscription Request"
// @Success 200 {object} model.EventResponse "Successfully subscribed to the event"
// @Failure 400 {object} model.Response "Bad Request - Invalid input"
// @Failure 500 {object} model.Response "Internal Server Error"
// @Router /events [post]
func (ctrl *EventController) SubscribeEvent(c *gin.Context) {
	userID := c.GetString("user_id")
	var request model.CreateEventRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(400, model.Response{
			Msg: err.Error(),
		})
		return
	}
	newEvent := &model.Event{
		UserID:         userID,
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

	c.JSON(200, model.EventResponse{
		Data: newEvent,
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
	id := c.Param("id")
	userID := c.GetString("user_id")
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
		UserID:         userID,
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

	c.JSON(200, model.EventResponse{
		Data: event,
	})
}

// DeleteEvent godoc
// @Summary Delete event
// @Description Deletes a specific event by its ID for a given user
// @Tags Event
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user_id header string true "User ID"
// @Param event_id path string true "Event ID"
// @Success 200 {object} model.Response "Event successfully deleted"
// @Failure 500 {object} model.Response "Internal Server Error"
// @Router /events/{event_id} [delete]
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
