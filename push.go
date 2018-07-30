package abios

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gobuffalo/uuid"
	"github.com/gorilla/websocket"

	. "github.com/PatronGG/abios-go-sdk/structs"
)

func (a *client) ListSubscriptions() ([]Subscription, error) {
	subs := []Subscription{}

	params := make(Parameters)
	params.Set("access_token", a.oauth.AccessToken)

	u, err := url.Parse(subscriptions)
	if err != nil {
		return subs, err
	}
	u.RawQuery = params.encode()
	res, err := http.Get(u.String())
	defer res.Body.Close()

	dec := json.NewDecoder(res.Body)
	if res.StatusCode == http.StatusOK {
		dec.Decode(&subs)
		return subs, nil
	}

	return subs, fmt.Errorf("Unexpected status code %v", res.StatusCode)
}

func (a *client) DeleteSubscription(id uuid.UUID) error {
	params := make(Parameters)
	params.Set("access_token", a.oauth.AccessToken)

	u, err := url.Parse(subscriptionsById + id.String())
	if err != nil {
		return err
	}
	u.RawQuery = params.encode()

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("Unexpected status code %v", res.StatusCode)
}

/*
func (a *client) UpdateSubscription(id int, sub Subscription) (Subscription, error) {

}
func (a *client) PushServiceConfig() ([]byte, error) {

}
*/

const timestampMillisFormat = "2006-01-02 15:04:05.000"

func (a *client) PushServiceConnect() error {
	params := make(Parameters)
	params.Set("access_token", a.oauth.AccessToken)

	if a.reconnectToken != uuid.Nil {
		params.Set("reconnect_token", a.reconnectToken.String())
	}
	var dialer *websocket.Dialer
	conn, res, err := dialer.Dial(wsBaseUrl, nil)

	if err == websocket.ErrBadHandshake {
		fmt.Printf("%s [ERROR]: Failed to connect to WS url. Handshake status='%d'\n",
			time.Now().Format(timestampMillisFormat), res.StatusCode)
		return err
	} else if err != nil {
		fmt.Printf("%s [ERROR]: Failed to connect to WS url. Error='%s'\n",
			time.Now().Format(timestampMillisFormat), err.Error())
		return err
	}

	a.wsConn = conn

	return nil
}

func (a *client) PushServiceInit() error {
	if err := a.PushServiceConnect(); err != nil {
		return err
	}

	initMsg, err := a.readInitMessage()
	if err != nil {
		return err
	}
	a.reconnectToken = initMsg.ReconnectToken

	errs := make(chan error, 1)

	go keepAliveLoop(errs)
	go messageReadLoop(errs)

	for _, e := range errors {
		fmt.Println(e)
	}

	return nil
}

func (a *client) handleInitMessage() (InitResponseMessage, error) {
	var m InitResponseMessage

	_, message, err := a.conn.ReadMessage()
	if closeErr, ok := err.(*websocket.CloseError); ok {
		var errMsg string
		switch closeErr.Code {
		case CloseUnknownSubscriptionID:
			errMsg = fmt.Sprintf("Subscription ID '%s' is not registered on server", subscriptionIDOrName)
		case CloseMissingSubscriptionID:
			errMsg = "Missing subscription ID or name in setup request"
		case CloseMaxNumSubscribers:
			errMsg = "The max number of concurrent subscribers for the account has been exceeded"
		case CloseMaxNumSubscriptions:
			errMsg = "The max number of registered subscriptions for the account has been exceeded"
		case CloseInternalError:
			errMsg = "Unknown server error"
		default:
			errMsg = fmt.Sprintf("Server sent unrecognized error code %d", closeErr.Code)
		}

		fmt.Printf("%s [ERROR]: Server closed connection: %s\n",
			time.Now().Format(timestampMillisFormat), errMsg)
		return m, err
	} else if err != nil {
		// Websocket read encountered some other error, we won't try to recover
		fmt.Printf("%s [ERROR]: Failed to read `init' message. Error='%s'\n",
			time.Now().Format(timestampMillisFormat), err.Error())
		return m, err
	}

	json.Unmarshal(initMsg, &m)
	return m, nil
}

// This will read messages from the server and print them to stdout.
// If the websocket is closed it will automatically re-establish the
// connection using the reconnect token to ensure no messages were lost
// during the disconnect.
func (a *client) messageReadLoop(errors <-chan error) {
	// From here on we will start receiving push events that match our
	// subscription filters
	for {
		_, message, err := conn.ReadMessage()

		// If the websocket is closed we need to reconnect
		if closeErr, ok := err.(*websocket.CloseError); ok {
			fmt.Printf("%s [INFO]: Websocket was closed, starting reconnect loop. Reason='%s'\n",
				time.Now().Format(timestampMillisFormat), closeErr.Error())

			// TODO: make sure to generate a new access token as the original one may be too old
			err = a.PushServiceConnect()
			if err != nil {
				errors <- err
				return
			}

			continue
		} else if err != nil {
			// Websocket read encountered some other error, we won't try to recover
			fmt.Printf("%s [ERROR]: Failed to read message. Error='%s'\n",
				time.Now().Format(timestampMillisFormat), err.Error())
			errors <- err
			return
		}

		// Sanity check that the JSON can be marshalled into the correct message
		// format
		var m PushMessage
		err = json.Unmarshal(message, &m)
		if err != nil {
			fmt.Printf("%s [ERROR]: Failed to unmarshal to message struct. Error='%s', Message='%s'\n",
				time.Now().Format(timestampMillisFormat), err.Error(), message)
			errors <- err
			// Ignore message and keep reading from websocket
			continue
		}

		log.Printf("msg: %v\n", message)
	}
}

func (a *client) keepAliveLoop(errors <-chan error) {
	for {
		time.Sleep(time.Second * 30)
		if conn != nil {
			err := a.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(3*time.Second))
			if err != nil {
				fmt.Printf("%s [ERROR]: Failed to send Ping message. Error='%s'\n",
					time.Now().Format(timestampMillisFormat), err.Error())
				errors <- err
				continue
			}
		}
	}
}
