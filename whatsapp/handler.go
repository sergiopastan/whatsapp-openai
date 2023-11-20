package whatsapp

import (
	"context"

	"go.mau.fi/whatsmeow/types/events"
)

type eventHandler interface {
	Execute(context.Context, *events.Message)
}

func MessageReceiptHandler(handler eventHandler) func(interface{}) {
	return func(rawEvt interface{}) {
		ctx := context.Background()
		switch evt := rawEvt.(type) {
		case *events.Message:
			handler.Execute(ctx, evt)
		}
	}
}
