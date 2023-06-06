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

func (subject *EventSubject) AddSubscriber(sub EventSubscriber) {
	subject.Observers.Store(sub, struct{}{})
}

func (subject *EventSubject) RemoveSubscriber(sub EventSubscriber) {
	subject.Observers.Delete(sub)
}

func (es *EventSubject) NotifySubscriber(event Event) {
	es.Observers.Range(func(key interface{}, value interface{}) bool {
		fmt.Println(es.ID, event.Data, key.(EventSubscriber).User)
		if key == nil {
			fmt.Printf("could not find matching Subscriber with id %v to event id: %v", event.TargetID, event.SubjectID)
			return false
		}
		es := key.(EventSubscriber)

		if es.User.ID == event.TargetID {
			// found matching subscriber
			es.NotifyCallback(event)
			return false
		}
		return true
	})
	// fmt.Printf("could not find matching Subscriber with id %v to event id: %v\n", event.TargetID, event.SubjectID)
}
