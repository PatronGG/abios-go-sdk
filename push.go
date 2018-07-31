package abios

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gobuffalo/uuid"
	"github.com/gorilla/websocket"

	. "github.com/PatronGG/abios-go-sdk/structs"
)

// Custom status codes sent by the server for the 'close' command.
// The websocket standard (RFC6455) allocates the
// 4000-4999 range to application specific status codes.
const (
	CloseMissingAccessToken    = 4000 // Missing access token in ws setup request
	CloseInvalidAccessToken    = 4001 // Invalid access token in ws setup request
	CloseNotAuthorized         = 4002 // Client account does not have access to the push API
	CloseMaxNumSubscribers     = 4003 // Max number of concurrent subscribers connected for client id
	CloseMaxNumSubscriptions   = 4004 // Max number of registered subscriptions exist for client id
	CloseInvalidReconnectToken = 4005 // Invalid reconnect token in ws setup request
	CloseMissingSubscriptionID = 4006 // Missing subscription id in ws setup request
	CloseUnknownSubscriptionID = 4007 // The supplied subscriber id in ws setup request does not exist in server
	CloseInternalError         = 4500 // Unspecified error due to problem in server
)

/*
func (a *client) UpdateSubscription(id int, sub Subscription) (Subscription, error) {

}
func (a *client) PushServiceConfig() ([]byte, error) {

}
*/

func (a *client) PushServiceConnect(subscriptionID uuid.UUID) error {
	params := make(Parameters)
	params.Set("access_token", a.oauth.AccessToken)
	params.Set("subscription_id", subscriptionID.String())

	if a.reconnectToken != uuid.Nil {
		params.Set("reconnect_token", a.reconnectToken.String())
	}

	u, err := url.Parse(wsBaseUrl)
	if err != nil {
		return err
	}
	u.RawQuery = params.encode()

	log.Printf("[INFO] Dialing to socket '%v'\n", u.String())

	var dialer *websocket.Dialer
	conn, res, err := dialer.Dial(u.String(), nil)

	if err == websocket.ErrBadHandshake {
		log.Printf("[ERROR]: Failed to connect to WS url. Handshake status='%d'\n", res.StatusCode)
		return err
	} else if err != nil {
		log.Printf("[ERROR]: Failed to connect to WS url. Error='%s'\n", err.Error())
		return err
	}

	a.wsConn = conn

	return nil
}

func (a *client) PushServiceInit(subscriptionID uuid.UUID) (chan SeriesMessage, chan error) {
	errors := make(chan error, 1)
	series := make(chan SeriesMessage, 1)

	if err := a.PushServiceConnect(subscriptionID); err != nil {
		errors <- err
		return series, errors
	}

	initMsg, err := a.handleInitMessage(subscriptionID)
	if err != nil {
		errors <- err
		return series, errors
	}
	a.reconnectToken = initMsg.ReconnectToken

	go a.keepAliveLoop(errors)
	go a.messageReadLoop(subscriptionID, series, errors)

	return series, errors
}

func (a *client) handleInitMessage(subscriptionID uuid.UUID) (InitResponseMessage, error) {
	var m InitResponseMessage

	_, message, err := a.wsConn.ReadMessage()
	if closeErr, ok := err.(*websocket.CloseError); ok {
		var errMsg string
		switch closeErr.Code {
		case CloseUnknownSubscriptionID:
			errMsg = fmt.Sprintf("Subscription ID '%s' is not registered on server", subscriptionID)
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

		log.Printf("[ERROR]: Server closed connection: %s\n", errMsg)
		return m, err
	} else if err != nil {
		// Websocket read encountered some other error, we won't try to recover
		log.Printf("[ERROR]: Failed to read `init' message. Error='%s'\n", err.Error())
		return m, err
	}

	json.Unmarshal(message, &m)
	return m, nil
}

func (a *client) messageReadLoop(subscriptionID uuid.UUID, series chan<- SeriesMessage, errors chan<- error) {
	for {
		_, message, err := a.wsConn.ReadMessage()

		if closeErr, ok := err.(*websocket.CloseError); ok {
			log.Printf("[INFO]: Websocket was closed, starting reconnect loop. Reason='%s'\n", closeErr.Error())

			authErr := a.authenticate()
			if authErr != nil {
				errors <- fmt.Errorf("%v", authErr)
				return
			}

			err = a.PushServiceConnect(subscriptionID)
			if err != nil {
				errors <- err
				return
			}

			continue
		} else if err != nil {
			log.Printf("[ERROR]: Failed to read message. Error='%s'\n", err.Error())
			errors <- err
			return
		}

		// sanity check
		var m PushMessage
		err = json.Unmarshal(message, &m)
		if err != nil {
			log.Printf("[ERROR]: Failed to unmarshal to message struct. Error='%s', Message='%s'\n", err.Error(), message)
			errors <- err
			continue
		}

		switch m.Channel {
		case "series":
			var s SeriesMessage
			err = json.Unmarshal(message, &s)
			if err != nil {
				log.Printf("[ERROR]: Failed to unmarshal to message struct. Error='%s', Message='%s'\n", err.Error(), message)
				errors <- err
				continue
			}

			s.Raw = message
			series <- s
		}

	}
}

func (a *client) keepAliveLoop(errors chan<- error) {
	for {
		time.Sleep(time.Second * 30)
		if a.wsConn != nil {
			err := a.wsConn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(3*time.Second))
			if err != nil {
				log.Printf("[ERROR]: Failed to send Ping message. Error='%s'\n", err.Error())
				errors <- err
				continue
			}
		}
	}
}
