package ext_test

import (
	"encoding/json"
	"sort"
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

func TestDispatcher(t *testing.T) {
	type testHandler struct {
		group     int
		shouldRun bool
		returnVal error
	}

	for name, testParams := range map[string]struct {
		handlers   []testHandler
		numMatches int
	}{
		"one group two handlers": {
			handlers: []testHandler{
				{
					group:     0,
					shouldRun: true,
					returnVal: nil,
				}, {
					group:     0,
					shouldRun: false, // same group, so doesnt run
					returnVal: nil,
				},
			},
			numMatches: 1,
		},
		"two handlers two groups": {
			handlers: []testHandler{
				{
					group:     0,
					shouldRun: true,
					returnVal: nil,
				}, {
					group:     1,
					shouldRun: true, // second group, so also runs
					returnVal: nil,
				},
			},
			numMatches: 2,
		},
		"end groups": {
			handlers: []testHandler{
				{
					group:     0,
					shouldRun: true,
					returnVal: ext.EndGroups,
				}, {
					group:     1,
					shouldRun: false, // ended, so second group doesnt run
					returnVal: nil,
				},
			},
			numMatches: 1,
		},
		"continue groups": {
			handlers: []testHandler{
				{
					group:     0,
					shouldRun: true,
					returnVal: ext.ContinueGroups,
				}, {
					group:     0,
					shouldRun: true, // continued, so second item in same group runs
					returnVal: nil,
				},
			},
			numMatches: 2,
		},
	} {
		name, testParams := name, testParams

		t.Run(name, func(t *testing.T) {
			d := ext.NewDispatcher(make(chan json.RawMessage), nil)
			var events []int
			for idx, h := range testParams.handlers {
				idx, h := idx, h

				t.Logf("Loading handler %d in group %d", idx, h.group)
				d.AddHandlerToGroup(handlers.NewMessage(message.All, func(b *gotgbot.Bot, ctx *ext.Context) error {
					if !h.shouldRun {
						t.Errorf("handler %d in group %d should not have run", idx, h.group)
						t.FailNow()
					}

					t.Logf("handler %d in group %d has run, as expected", idx, h.group)
					events = append(events, idx)
					return h.returnVal
				}), h.group)
			}

			t.Log("Processing one update...")
			d.ProcessUpdate(nil, &gotgbot.Update{
				Message: &gotgbot.Message{Text: "test text"},
			}, nil)

			// ensure events handled in order
			if !sort.IntsAreSorted(events) {
				t.Errorf("order of events is not sorted: %v", events)
			}
			if len(events) != testParams.numMatches {
				t.Errorf("got %d matches, expected %d ", len(events), testParams.numMatches)
			}
		})
	}
}
