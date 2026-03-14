package main

import (
	"net/http"
	"rest-api-go-gin/internal/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *application) createEvent(c *gin.Context) {
	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := app.GetUserFromContext(c)
	event.OwnerId = user.Id

	err := app.models.Events.Insert(&event)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// getEvents return all events
//
// @Summary Returns all events
// @Description Returns all events
// @Tags Events
// @Accept json
// @Produce json
// @Success 200 {object} []database.Event
// @Router /api/v1/events [get]
func (app *application) getAllEvents(c *gin.Context) {
	events, err := app.models.Events.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (app *application) getEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive event"})
		return
	}

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, event)

}

func (app *application) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	user := app.GetUserFromContext(c)
	event, err := app.models.Events.Get(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this event"})
		return
	}

	updatedEvent := &database.Event{}

	if err := c.ShouldBindJSON(updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEvent.Id = id

	if err := app.models.Events.Update(updatedEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, updatedEvent)
}

func (app *application) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	user := app.GetUserFromContext(c)
	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this event"})
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (app *application) addAttendeeToEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	event, err := app.models.Events.Get(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	userToAdd, err := app.models.Users.Get(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive user"})
		return
	}
	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user := app.GetUserFromContext(c)
	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to add attendee to this event"})
		return
	}

	ExistingAttendee, err := app.models.Attendees.GetByEventAndAttendee(eventId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive attendee"})
		return
	}

	if ExistingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Attendee already exists"})
		return
	}

	attendee := database.Attendee{
		EventId: eventId,
		UserId:  userId,
	}

	_, err = app.models.Attendees.Insert(&attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attendee"})
		return
	}

	c.JSON(http.StatusCreated, attendee)
}

func (app *application) getAttendeesForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
		return
	}

	users, err := app.models.Attendees.GetAttendeesByEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Faild to retrive attendees for event"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (app *application) deleteAttendeeFromEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	event, err := app.models.Events.Get(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	user := app.GetUserFromContext(c)
	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete attendee from this event"})
		return
	}

	err = app.models.Attendees.Delete(userId, eventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (app *application) getEventsByAttendee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
		return
	}

	events, err := app.models.Attendees.GetEventsByAttendee(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get events"})
		return
	}

	c.JSON(http.StatusOK, events)
}
