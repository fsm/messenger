package messenger

import (
	"errors"
	"os"
	"reflect"
	"time"

	"github.com/BrandonRomano/wrecker"
	"github.com/fsm/emitable"
)

type MessageData struct {
	Recipient    MessageRecipient `json:"recipient"`
	SenderAction string           `json:"sender_action,omitempty"`
	Message      *SendMessageData `json:"message"`
}

type MessageRecipient struct {
	ID string `json:"id"`
}

type SendMessageData struct {
	Text         string       `json:"text,omitempty"`
	Attachment   *Attachment  `json:"attachment,omitempty"`
	QuickReplies []QuickReply `json:"quick_replies,omitempty"`
}

type Attachment struct {
	Type    string  `json:"type"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	URL string `json:"url"`
}

type FacebookEmitter struct {
	UUID string
}

type QuickReply struct {
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
	Payload     string `json:"payload"`
}

const typingTime = 1000

func (f *FacebookEmitter) Emit(input interface{}) error {
	switch v := input.(type) {
	case string:
		f.Emit(emitable.Typing{Enabled: true})
		f.Emit(emitable.Sleep{LengthMillis: typingTime})
		SendMessage(&MessageData{
			Recipient: MessageRecipient{
				ID: f.UUID,
			},
			Message: &SendMessageData{
				Text: v,
			},
		})
		return nil

	case emitable.Audio:
		f.Emit(emitable.Typing{Enabled: true})
		f.Emit(emitable.Sleep{LengthMillis: typingTime})
		SendMessage(&MessageData{
			Recipient: MessageRecipient{
				ID: f.UUID,
			},
			Message: &SendMessageData{
				Attachment: &Attachment{
					Type: "audio",
					Payload: Payload{
						URL: v.URL,
					},
				},
			},
		})
		return nil

	case emitable.File:
		f.Emit(emitable.Typing{Enabled: true})
		f.Emit(emitable.Sleep{LengthMillis: typingTime})
		SendMessage(&MessageData{
			Recipient: MessageRecipient{
				ID: f.UUID,
			},
			Message: &SendMessageData{
				Attachment: &Attachment{
					Type: "file",
					Payload: Payload{
						URL: v.URL,
					},
				},
			},
		})
		return nil

	case emitable.Image:
		f.Emit(emitable.Typing{Enabled: true})
		f.Emit(emitable.Sleep{LengthMillis: typingTime})
		SendMessage(&MessageData{
			Recipient: MessageRecipient{
				ID: f.UUID,
			},
			Message: &SendMessageData{
				Attachment: &Attachment{
					Type: "image",
					Payload: Payload{
						URL: v.URL,
					},
				},
			},
		})
		return nil

	case emitable.Video:
		f.Emit(emitable.Typing{Enabled: true})
		f.Emit(emitable.Sleep{LengthMillis: typingTime})
		SendMessage(&MessageData{
			Recipient: MessageRecipient{
				ID: f.UUID,
			},
			Message: &SendMessageData{
				Attachment: &Attachment{
					Type: "video",
					Payload: Payload{
						URL: v.URL,
					},
				},
			},
		})
		return nil

	case emitable.QuickReply:
		f.Emit(emitable.Typing{Enabled: true})
		f.Emit(emitable.Sleep{LengthMillis: typingTime})
		replies := make([]QuickReply, 0)
		for _, reply := range v.Replies {
			replies = append(replies, QuickReply{
				ContentType: "text",
				Title:       reply,
				Payload:     reply,
			})
		}
		SendMessage(&MessageData{
			Recipient: MessageRecipient{
				ID: f.UUID,
			},
			Message: &SendMessageData{
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
		SendMessage(&MessageData{
			Recipient: MessageRecipient{
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

func SendMessage(m *MessageData) {
	client := wrecker.New("https://graph.facebook.com/v2.6")
	client.Post("/me/messages").
		URLParam("access_token", os.Getenv("MESSENGER_ACCESS_TOKEN")).
		Body(m).
		Execute()
}
