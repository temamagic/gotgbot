package ext_test

import (
	"encoding/json"
	"fmt"
	"sort"
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

func TestDispatcher_TwoHandlersOneGroup(t *testing.T) {
	updateChan := make(chan json.RawMessage)
	d := ext.NewDispatcher(updateChan, nil)

	var orderOfEvents []int
	// Add to handlers to group 0 (default). First should run, second should NOT.
	d.AddHandler(handlers.NewMessage(message.All, func(b *gotgbot.Bot, ctx *ext.Context) error {
		orderOfEvents = append(orderOfEvents, 0)
		return nil
	}))
	d.AddHandler(handlers.NewMessage(message.All, func(b *gotgbot.Bot, ctx *ext.Context) error {
		orderOfEvents = append(orderOfEvents, 0)
		fmt.Println("should not execute the second handler in group 0")
		t.Fail()
		return nil
	}))

	d.AddHandlerToGroup(handlers.NewMessage(message.All, func(b *gotgbot.Bot, ctx *ext.Context) error {
		orderOfEvents = append(orderOfEvents, 1)
		return nil
	}), 1)

	d.ProcessUpdate(nil, &gotgbot.Update{
		Message: &gotgbot.Message{Text: "test text"},
	}, nil)

	// ensure events handled in order
	if !sort.IntsAreSorted(orderOfEvents) || len(orderOfEvents) != 2 {
		// only one item should be triggered
		fmt.Println("order of events is not sorted, or was not 2")
		t.Fail()
	}
}
