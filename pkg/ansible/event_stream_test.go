package ansible

import (
	"bytes"
	"testing"
)

func TestEventStreamSingleEvent(t *testing.T) {
	in := bytes.NewBufferString(`{"eventType":"PLAY_START", "eventData": {"name":"somePlay"}}`)
	es := EventStream(in)

	gotEvent := false
	for e := range es {
		switch event := e.(type) {
		default:
			gotEvent = true
			t.Error("Invalid event type received")
		case *PlayStartEvent:
			gotEvent = true
			if event.Name != "somePlay" {
				t.Errorf("Expected play name %q, but got %q", "somePlay", event.Name)
			}
		}
	}
	if !gotEvent {
		t.Errorf("Did not get the event")
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
		t.Errorf("invalid number of events received. expected 3, got %d", i)
	}
}

func TestEventStreamNoEvents(t *testing.T) {
	es := EventStream(bytes.NewBufferString(""))
	for e := range es {
		t.Errorf("got an unexpected event: %v", e)
	}
}

func TestEventStreamBadEventIsIgnored(t *testing.T) {
	in := bytes.NewBufferString(`{"eventType":"PLAY_START", "eventData": {"name":"somePlay"}}
{"eventType":"BAD_EVENT", "eventData": {"name":"somePlay"}}
someBadStuffHere...
{"eventType":"PLAY_START", "eventData": {"name":"somePlay"}}
`)
	es := EventStream(in)
	expectedGoodEvents := 2
	gotEvents := 0
	for _ = range es {
		gotEvents++
	}
	if gotEvents != expectedGoodEvents {
		t.Errorf("got %d events, but expected %d", gotEvents, expectedGoodEvents)
	}
}
