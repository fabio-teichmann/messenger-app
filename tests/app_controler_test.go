package tests

import (
	"context"
	"messenger-app/models"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestAcceptEventMsgSent(t *testing.T) {
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

	case e := <-ac.NewUser.Queue:
		t.Errorf("unexpected msg in channel, want nil, received %v", e)

	case e := <-ac.UserLogIn.Queue:
		t.Errorf("unexpected msg in channel, want nil, received %v", e)

	case e := <-ac.UserLogOut.Queue:
		t.Errorf("unexpected msg in channel, want nil, received %v", e)

	case e := <-ac.UserOnl.Queue:
		t.Errorf("unexpected msg in channel, want nil, received %v", e)
	}
}

func TestAcceptEventMsgRcvd(t *testing.T) {
	// create AppControler
	ac := models.InitializeAppControler(&mongo.Client{})
	// create event
	event := models.Event{SubjectID: models.MSG_RECEIVED}

	ac.AcceptEvent(context.TODO(), &event)

	select {
	case e := <-ac.MsgSent.Queue:
		t.Errorf("unexpected msg in channel, want nil, received %v", e)
	case e := <-ac.MsgRcvd.Queue:
		// test if messages are the same
		if !reflect.DeepEqual(e, event) {
			t.Errorf("wrong event in channel, want %v, received %v", event, e)
		}

	case e := <-ac.NewUser.Queue:
		t.Errorf("unexpected msg in channel, want nil, received %v", e)

	case e := <-ac.UserLogIn.Queue:
		t.Errorf("unexpected msg in channel, want nil, received %v", e)

	case e := <-ac.UserLogOut.Queue:
		t.Errorf("unexpected msg in channel, want nil, received %v", e)

	case e := <-ac.UserOnl.Queue:
		t.Errorf("unexpected msg in channel, want nil, received %v", e)
	}
}

func TestAcceptEventNewUser(t *testing.T) {
	// create AppControler
	ac := models.InitializeAppControler(&mongo.Client{})
	// create event
	event := models.Event{SubjectID: models.NEW_USER}

	ac.AcceptEvent(context.TODO(), &event)

	select {
	case e := <-ac.MsgSent.Queue:
		t.Errorf("unexpected msg in channel, want nil, received %v", e)

	case e := <-ac.MsgRcvd.Queue:
		t.Errorf("unexpected msg in channel, want nil, received %v", e)

	case e := <-ac.NewUser.Queue:
		// test if messages are the same
		if !reflect.DeepEqual(e, event) {
			t.Errorf("wrong event in channel, want %v, received %v", event, e)
		}

	case e := <-ac.UserLogIn.Queue:
		t.Errorf("unexpected msg in channel, want nil, received %v", e)

	case e := <-ac.UserLogOut.Queue:
		t.Errorf("unexpected msg in channel, want nil, received %v", e)

	case e := <-ac.UserOnl.Queue:
		t.Errorf("unexpected msg in channel, want nil, received %v", e)
	}
}
