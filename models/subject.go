package models

import "fmt"

type Subscriber interface {
	NotifyCallback(Event)
}

type Subject interface {
	AddSubscriber(Subscriber)
	RemoveSubscriber(Subscriber)
	NotifySubscriber(Event)
}

func NewEventSubject(id int) *EventSubject {
	return &EventSubject{ID: id}
}

func (subject *EventSubject) AddSubscriber(sub *EventSubscriber) {
	// fmt.Println("AddSubscriber:", sub)
	subject.Observers.Store(sub, struct{}{})
}

func (subject *EventSubject) RemoveSubscriber(sub EventSubscriber) {
	subject.Observers.Delete(sub)
}

func (es *EventSubject) NotifySubscriber(event Event) {
	es.Observers.Range(func(key interface{}, value interface{}) bool {
		// fmt.Println(es.ID, event.Data, key.(EventSubscriber).User)
		if key == nil {
			fmt.Printf("could not find matching Subscriber %s to event: %s", event.Target, event)
			return false
		}
		es := key.(*EventSubscriber)

		if es.User.ID == event.Target.ID {
			// found matching subscriber
			es.NotifyCallback(event)
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

// func (es *EventSubject) ReadEvents(control chan ControlMsg) {

// 	for {
// 		select {
// 		case msg := <-control:
// 			switch msg {
// 			case DoExit:
// 				fmt.Printf("exit read events for subject %v\n", es.ID)
// 				control <- ExitOK
// 				return
// 			}
// 		case event := <-es.Queue:
// 			// notify
// 			es.NotifySubscriber(event)
// 			// fmt.Println(msg)
// 		}
// 	}
// }
