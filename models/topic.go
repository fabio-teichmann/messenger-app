package models

// import (
// 	"fmt"
// 	"sync"
// )

// type Topic interface {
// 	// A topic is created using a SubjectID.
// 	// A topic listens to events on a Subject's queue
// 	// and facilitates different actions:
// 	//
// 	// - accepting events to the Subject's queue
// 	// - notify Subscribers about relevant events on
// 	//   this Subject's queue
// 	ReadEvents(chan ControlMsg) // from User centred queue

// 	AcceptEvent(Event)

// 	NotifySubscriber(Event)

// 	UpdateChat(User, Event) // update messages in Chat.History
// }

// type EventTopic struct {
// 	SubjectID int
// 	Queue     chan Event
// 	Observers sync.Map
// }

// func NewEventTopic(subjectID int) EventTopic {
// 	return EventTopic{
// 		SubjectID: subjectID,
// 		Queue:     make(chan Event),
// 	}
// }

// func (et *EventTopic) AcceptEvent2(event *Event) {
// 	if event == nil {
// 		fmt.Println("event is undefined")
// 		return
// 	}
// 	et.Queue <- *event
// }

// func (et *EventTopic) NotifySubscriber(event Event) {
// 	et.Observers.Range(func(key interface{}, value interface{}) bool {
// 		if key == nil {
// 			fmt.Printf("could not find matching Subscriber with id %v to event id: %v", event.Target.ID, event.SubjectID)
// 			return false
// 		}
// 		es := key.(EventSubscriber)

// 		if es.User.ID == event.Target.ID {
// 			// found matching subscriber
// 			es.NotifyCallback(&event)
// 			return false
// 		}
// 		return true
// 	})
// }

// func (et *EventTopic) ReadEvents2(control chan ControlMsg) {

// 	for {
// 		select {
// 		case msg := <-control:
// 			switch msg {
// 			case DoExit:
// 				fmt.Printf("exit read events for subject %v\n", et.SubjectID)
// 				control <- ExitOK
// 				return
// 			}
// 		case event := <-et.Queue:
// 			// notify
// 			et.NotifySubscriber(event)
// 			// fmt.Println(msg)
// 		}
// 	}
// }
