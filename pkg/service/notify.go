package service

import (
	"bikefest/pkg/bootstrap"
	"bikefest/pkg/model"
	"encoding/json"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeEventReminder = "reminder"
)

// Task payload for any event notification related tasks.
type eventNotificationPayload struct {
	UserID  string
	EventID string
}

type AsynqServiceImpl struct {
	client *asynq.Client
	env    *bootstrap.Env
}

func newEventNotification(user_id, event_id string) (*asynq.Task, error) {
	payload, err := json.Marshal(eventNotificationPayload{UserID: user_id, EventID: event_id})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEventReminder, payload), nil
}

func (as *AsynqServiceImpl) EnqueueEvent(user_id, event_id, event_start_time string) {
	t, err := newEventNotification(user_id, event_id)
func (as *AsynqServiceImpl) EnqueueEventNotification(userID, eventID, eventStartTime string) {
	t, err := newEventNotification(userID, eventID)
	if err != nil {
		log.Fatal(err)
	}

	location, _ := time.LoadLocation(as.env.Server.TimeZone)
	//TODO: currently we only set the process time 30 minutes before the event start time
	processTime, _ := time.ParseInLocation(model.EventTimeLayout, eventStartTime, location)
	processTime = processTime.Add(-time.Minute * 30)

	info, err := as.client.Enqueue(t, asynq.ProcessAt(processTime), asynq.TaskID(userID+eventID))

	if err != nil {
		log.Fatal(err)
	}
	log.Printf(" [*] Successfully enqueued task: %+v\nThe task should be executed at %s", info, processTime.String())
}

func NewAsynqService(client *asynq.Client, env *bootstrap.Env) model.AsynqNotificationService {
	return &AsynqServiceImpl{
		client: client,
		env:    env,
	}
}
