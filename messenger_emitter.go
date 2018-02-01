package messenger

import (
	"errors"
	"os"
	"reflect"
	"time"

	"github.com/BrandonRomano/wrecker"
	"github.com/fsm/emitable"
)

type messageData struct {
	Recipient    messageRecipient `json:"recipient"`
	SenderAction string           `json:"sender_action,omitempty"`
	Message      *sendMessageData `json:"message"`
}

type messageRecipient struct {
	ID string `json:"id"`
}

type sendMessageData struct {
	Text         string       `json:"text,omitempty"`
	Attachment   *attachment  `json:"attachment,omitempty"`
	QuickReplies []quickReply `json:"quick_replies,omitempty"`
}

type attachment struct {
	Type    string  `json:"type"`
	Payload payload `json:"payload"`
}

type payload struct {
	URL string `json:"url"`
}

type facebookEmitter struct {
	UUID string
}

type quickReply struct {
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
	Payload     string `json:"payload"`
}

const typingTime = 1000

func (f *facebookEmitter) Emit(input interface{}) error {
	switch v := input.(type) {
	case string:
		f.Emit(emitable.Typing{Enabled: true})
		f.Emit(emitable.Sleep{LengthMillis: typingTime})
		sendMessage(&messageData{
			Recipient: messageRecipient{
				ID: f.UUID,
			},
			Message: &sendMessageData{
				Text: v,
			},
		})
		return nil

	case emitable.Audio:
		f.Emit(emitable.Typing{Enabled: true})
		f.Emit(emitable.Sleep{LengthMillis: typingTime})
		sendMessage(&messageData{
			Recipient: messageRecipient{
				ID: f.UUID,
			},
			Message: &sendMessageData{
				Attachment: &attachment{
					Type: "audio",
					Payload: payload{
						URL: v.URL,
					},
				},
			},
		})
		return nil

	case emitable.File:
		f.Emit(emitable.Typing{Enabled: true})
		f.Emit(emitable.Sleep{LengthMillis: typingTime})
		sendMessage(&messageData{
			Recipient: messageRecipient{
				ID: f.UUID,
			},
			Message: &sendMessageData{
				Attachment: &attachment{
					Type: "file",
					Payload: payload{
						URL: v.URL,
					},
				},
			},
		})
		return nil

	case emitable.Image:
		f.Emit(emitable.Typing{Enabled: true})
		f.Emit(emitable.Sleep{LengthMillis: typingTime})
		sendMessage(&messageData{
			Recipient: messageRecipient{
				ID: f.UUID,
			},
			Message: &sendMessageData{
				Attachment: &attachment{
					Type: "image",
					Payload: payload{
						URL: v.URL,
					},
				},
			},
		})
		return nil

	case emitable.Video:
		f.Emit(emitable.Typing{Enabled: true})
		f.Emit(emitable.Sleep{LengthMillis: typingTime})
		sendMessage(&messageData{
			Recipient: messageRecipient{
				ID: f.UUID,
			},
			Message: &sendMessageData{
				Attachment: &attachment{
					Type: "video",
					Payload: payload{
						URL: v.URL,
					},
				},
			},
		})
		return nil

	case emitable.QuickReply:
		f.Emit(emitable.Typing{Enabled: true})
		f.Emit(emitable.Sleep{LengthMillis: typingTime})
		replies := make([]quickReply, 0)
		for _, reply := range v.Replies {
			replies = append(replies, quickReply{
				ContentType: "text",
				Title:       reply,
				Payload:     reply,
			})
		}
		sendMessage(&messageData{
			Recipient: messageRecipient{
				ID: f.UUID,
			},
			Message: &sendMessageData{
				Text:         v.Message,
				QuickReplies: replies,
			},
		})
		return nil

	case emitable.Typing:
		action := "typing_off"
		if v.Enabled {
			action = "typing_on"
		}
		sendMessage(&messageData{
			Recipient: messageRecipient{
				ID: f.UUID,
			},
			SenderAction: action,
		})
		return nil

	case emitable.Sleep:
		time.Sleep(time.Millisecond * time.Duration(v.LengthMillis))
		return nil
	}

	return errors.New("FacebookEmitter cannot handle " + reflect.TypeOf(input).String())
}

func sendMessage(m *messageData) {
	client := wrecker.New("https://graph.facebook.com/v2.6")
	client.Post("/me/messages").
		URLParam("access_token", os.Getenv("MESSENGER_ACCESS_TOKEN")).
		Body(m).
		Execute()
}
