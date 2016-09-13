package ansible

import (
	"bytes"
	"testing"
)

func TestEventStreamSingleEvent(t *testing.T) {
	in := bytes.NewBufferString(`{"eventType":"PLAY_START", "eventData": {"name":"somePlay"}}`)
	es := EventStream(in)

	for e := range es {
		switch event := e.(type) {
		default:
			t.Error("Invalid event type received")
		case *PlayStartEvent:
			if event.Name != "somePlay" {
				t.Errorf("Expected play name %q, but got %q", "somePlay", event.Name)
			}
		}
	}

}

func TestEventStreamMultipleEvents(t *testing.T) {
	in := bytes.NewBufferString(`{"eventType":"PLAY_START", "eventData": {"name":"somePlay"}}
{"eventType":"PLAY_START", "eventData": {"name":"somePlay"}}
{"eventType":"PLAY_START", "eventData": {"name":"somePlay"}}
`)

	es := EventStream(in)

	i := 0
	for e := range es {
		switch e.(type) {
		default:
			t.Error("invalid event type received")
		case *PlayStartEvent:
			i++
		}
	}

	if i != 3 {
		t.Errorf("invalid number of events received")
	}
}
