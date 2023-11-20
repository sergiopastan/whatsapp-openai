package whatsapp

import (
	"context"

	"go.mau.fi/whatsmeow/types/events"
)

type eventHandler interface {
	Execute(context.Context, *events.Message)
}

func MessageReceiptHandler(ctx context.Context, handler eventHandler) func(interface{}) {
	return func(rawEvt interface{}) {
		switch evt := rawEvt.(type) {
		case *events.Message:
			handler.Execute(ctx, evt)
		}
	}
}
