package models

import (
	"context"
	"fmt"
	"sync"
)

type Subject interface {
	AddSubscriber(Subscriber)
	RemoveSubscriber(Subscriber)
	NotifySubscriber(Event)
}

type EventSubject struct {
	EventType EventType
	Queue     chan Event
	Observers sync.Map
}

func NewEventSubject(eventType EventType) *EventSubject {
	return &EventSubject{EventType: eventType, Queue: make(chan Event)}
}

func (subject *EventSubject) AddSubscriber(sub *EventSubscriber) {
	// fmt.Println("AddSubscriber:", sub)
	subject.Observers.Store(sub, struct{}{})
}

func (subject *EventSubject) RemoveSubscriber(sub EventSubscriber) {
	subject.Observers.Delete(sub)
}

func (es *EventSubject) NotifySubscriber(ctx context.Context, ac *AppControler, event *Event) {
	es.Observers.Range(func(key interface{}, value interface{}) bool {
		// fmt.Println(es.ID, event.Data, key.(EventSubscriber).User)
		if key == nil {
			fmt.Printf("could not find matching Subscriber %s to event: %v", event.Target.User.Name, event)
			return false
		}
		subscriber := key.(*EventSubscriber)

		if subscriber.User.ID == event.Target.User.ID {
			// found matching subscriber
			subscriber.NotifyCallback(ctx, ac, event)
			return false
		}
		return true
	})
	// fmt.Printf("could not find matching Subscriber with id %v to event id: %v\n", event.TargetID, event.SubjectID)
}

// func (es *EventSubject) AcceptEvent(event *Event) {
// 	if event == nil {
// 		fmt.Println("event is undefined")
// 		return
// 	}
// 	es.Queue <- *event
// }

// func (es *EventSubjectNew) ReadEvents(control chan ControlMsg) {

// 	for {
// 		select {
// 		case msg := <-control:
// 			switch msg {
// 			case DoExit:
// 				fmt.Printf("exit read events for subject %v\n", es.EventType)
// 				control <- ExitOK
// 				return
// 			}
// 		case event := <-es.Queue:
// 			// notify
// 			es.NotifySubscriber(&event)
// 			// fmt.Println(msg)
// 		}
// 	}
// }
