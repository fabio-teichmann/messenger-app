package tests

import (
	"messenger-app/models"
	"testing"
)

func TestAddDeleteChat(t *testing.T) {
	sub1 := models.NewEventSubscriberByName("test1")
	sub2 := models.NewEventSubscriberByName("test2")

	if len(sub1.Chats) != 0 {
		t.Fatalf("unexpected chat registered; got %v, want %v\n", sub1.Chats, nil)
	}

	sub1.AddChat(sub1)
	if len(sub1.Chats) != 0 {
		t.Fatalf("chat with owner ID registered; got %v, want %v\n", sub1.ID, nil)
	}

	sub1.AddChat(sub2)
	if len(sub1.Chats) == 0 {
		t.Fatalf("chat not registered; got %v, want %v\n", sub1.Chats, nil)
	}

	sub1.AddChat(sub2)
	if len(sub1.Chats) != 1 {
		t.Fatalf("chat with identical ID registered; got %v, want %v\n", sub1.Chats, nil)
	}

	sub1.DeleteChat(sub2)
	if len(sub1.Chats) != 0 {
		t.Fatalf("chat not removed; got %v, want %v", sub1.Chats, nil)
	}
}
