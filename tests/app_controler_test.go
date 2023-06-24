package tests

import (
	"context"
	"messenger-app/models"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestAcceptEvent(t *testing.T) {
	// create AppControler
	ac := models.InitializeAppControler(&mongo.Client{})
	// create event
	event := models.Event{SubjectID: models.MSG_SENT}

	ac.AcceptEvent(context.TODO(), &event)

	select {
	case e := <-ac.MsgSent.Queue:
		// test if messages are the same
		if !reflect.DeepEqual(e, event) {
			t.Errorf("wrong event in channel, want %v, received %v", event, e)
		}
	case e := <-ac.MsgRcvd.Queue:
		t.Errorf("unexpected msg in channel, want nil, received %v", e)
	}

}
