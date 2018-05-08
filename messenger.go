package messenger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/fsm/fsm"
	targetutil "github.com/fsm/target-util"
)

const platform = "facebook-messenger"

// SetupWebhook adds support for the Messenger Platform's webhook verification
// to your webhook. This is required to ensure your webhook is authentic and working.
//
//  This must be a GET request, and have the same URL as the POST request.
//
// https://developers.facebook.com/docs/messenger-platform/getting-started/webhook-setup
func SetupWebhook(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	mode := queryParams.Get("hub.mode")
	challenge := queryParams.Get("hub.challenge")
	verifyToken := queryParams.Get("hub.verify_token")

	if mode == "subscribe" && verifyToken == os.Getenv("MESSENGER_VERIFY_TOKEN") {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

// GetMessageReceivedWebhook is the webhook that facebook posts to when a message is
// received from a user.
//
// This must be a POST request, and have the same URL as the GET request.
//
// https://developers.facebook.com/docs/messenger-platform/getting-started/webhook-setup
func GetMessageReceivedWebhook(stateMachine fsm.StateMachine, store fsm.Store) func(http.ResponseWriter, *http.Request) {
	// Build Statemap
	stateMap := targetutil.GetStateMap(stateMachine)

	// Return HandlerFunc
	return func(w http.ResponseWriter, r *http.Request) {
		// Get body
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()

		// Parse body into struct
		cb := new(messageReceivedCallback)
		json.Unmarshal([]byte(body), cb)

		// For each entry
		for _, i := range cb.Entry {
			// Iterate over each messaging event
			for _, messagingEvent := range i.MessagingEvents {
				// Perform a Step
				go targetutil.Step(
					platform,
					messagingEvent.Sender.ID,
					messagingEvent.Message.Text,
					store,
					&facebookEmitter{
						UUID: messagingEvent.Sender.ID,
					},
					stateMap,
				)
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
