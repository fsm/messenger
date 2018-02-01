<a href="https://github.com/fsm"><p align="center"><img src="https://user-images.githubusercontent.com/2105067/35464215-a014d512-02a9-11e8-8913-63a066f6064e.png" alt="FSM" width="350px" align="center;"/></p></a>

# Messenger

Messenger is a Facebook Messenger target for [fsm](https://github.com/fsm/fsm).

## Environment Variables

When using this target, you must set two environment variables:

```
MESSENGER_VERIFY_TOKEN=""
MESSENGER_ACCESS_TOKEN=""
```

## Getting Started

> Note: The environment variables above are assumed to be set in this example code:

```go
package main

import (
	"net/http"

	"github.com/fsm/fsm"
	"github.com/fsm/messenger"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := &httprouter.Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
    }

	router.HandlerFunc(http.MethodGet, "/facebook", messenger.SetupWebhook)
    router.HandlerFunc(http.MethodPost, "/facebook", messenger.GetMessageReceivedWebhook(getStateMachine(), getStore()))

	http.ListenAndServe(":5000", router)
}

func getStateMachine() fsm.StateMachine {
	// ...
}

func getStore() fsm.Store {
	// ...
}
```

# License

[MIT](LICENSE.md)
