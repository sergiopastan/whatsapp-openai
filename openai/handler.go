package openai

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"

	"github.com/sergiopastan/whatsapp-openai/whatsapp"
)

type Handler struct {
	openaiClient *openai.Client
	wspClient    *whatsapp.Client
}

func NewHandler(apiKey string, wspClient *whatsapp.Client) *Handler {
	return &Handler{
		openaiClient: openai.NewClient(apiKey),
		wspClient:    wspClient,
	}
}

func (h *Handler) getResponse(ctx context.Context, msg string) string {
	log.Infof("getResponse msg: %v", msg)
	resp, err := h.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: msg,
				},
			},
		},
	)

	if err != nil {
		log.Infof("chatcompletion error: %v", err)
		return "ups, intenta de nuevo"
	}

	ans := resp.Choices[0].Message.Content
	log.Infof("openai message: %v", ans)
	return ans
}

func (h *Handler) Execute(ctx context.Context, evt *events.Message) {
	msg := evt.Message.GetConversation()
	if evt.Info.IsFromMe || evt.Info.IsGroup {
		log.Infof("msg is from me or from group %v", evt)
		return
	}
	_, err := h.wspClient.SendMessage(ctx, types.NewJID(evt.Info.Chat.User, evt.Info.Chat.Server), &waProto.Message{
		Conversation: proto.String(h.getResponse(ctx, msg)),
	})
	if err != nil {
		log.Infof("error sending the message %v", evt)
	}
}
