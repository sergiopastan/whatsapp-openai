package openai

import (
	"context"
	"fmt"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"

	"github.com/sergiopastan/whatsapp-openai/whatsapp"
)

type Handler struct {
	wspClient *whatsapp.Client
}

func NewHandler(wspClient *whatsapp.Client) *Handler {
	return &Handler{
		wspClient: wspClient,
	}
}

func (h Handler) Execute(ctx context.Context, evt *events.Message) {
	msg := evt.Message.GetConversation()
	_, err := h.wspClient.SendMessage(ctx, types.NewJID(evt.Info.Chat.User, evt.Info.Chat.Server), &waProto.Message{
		Conversation: proto.String(msg),
	})
	if err != nil {
		fmt.Println("Error sending the flag message")
	}
}
