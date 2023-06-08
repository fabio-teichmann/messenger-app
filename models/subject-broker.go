package models

import (
	"fmt"
)

type SubjectBroker interface {
	// A Subject is created using a SubjectID.
	// A Subject listens to events on a Subject's queue
	// and facilitates different actions:
	//
	// - accepting events to the Subject's queue
	// - notify Subscribers about relevant events on
	//   this Subject's queue
	AcceptEvent(Event)
	ReadEvents()
}

type EventSubjectBroker struct {
	EventSubject *EventSubject
	Queue        chan Event
	ControlChan  chan int
}

func NewEventSubjectBroker(es *EventSubject) EventSubjectBroker {
	return EventSubjectBroker{
		EventSubject: es,
		Queue:        make(chan Event),
		ControlChan:  make(chan int, 2),
	}
}

func (esb *EventSubjectBroker) AcceptEvent(event *Event) {
	if event == nil {
		fmt.Println("event is undefined")
		return
	}
	esb.Queue <- *event
}

func (esb *EventSubjectBroker) ReadEvents() {

	for {
		select {
		case msg := <-esb.ControlChan:
			switch msg {
			case DoExit:
				fmt.Printf("exit read events for subject %v\n", esb.EventSubject.ID)
				esb.ControlChan <- ExitOK
				return
			}
		case event := <-esb.Queue:
			// notify
			esb.EventSubject.NotifySubscriber(&event)
			// fmt.Println(msg)
		}
	}
}
