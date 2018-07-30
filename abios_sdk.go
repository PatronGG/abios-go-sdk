package abios

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	. "github.com/PatronGG/abios-go-sdk/structs"
	"github.com/gobuffalo/uuid"
	"github.com/gorilla/websocket"
)

// Constant variables that represents endpoints
const (
	baseUrl           = "https://api.abiosgaming.com/v2/"
	errorEndpoint     = baseUrl + "error"
	access_token      = baseUrl + "oauth/access_token"
	games             = baseUrl + "games"
	series            = baseUrl + "series"
	seriesById        = series + "/"
	matches           = baseUrl + "matches/"
	tournaments       = baseUrl + "tournaments"
	tournamentsById   = tournaments + "/"
	substages         = baseUrl + "substages/"
	teams             = baseUrl + "teams"
	teamsById         = teams + "/"
	players           = baseUrl + "players"
	playersById       = players + "/"
	rosters           = baseUrl + "rosters/"
	search            = baseUrl + "search"
	incidents         = baseUrl + "incidents"
	incidentsBySeries = incidents + "/"
	organisations     = baseUrl + "organisations"
	organisationsById = organisations + "/"

	// PUSH API
	wsBaseUrl         = "https://ws.abiosgaming.com/v0/"
	subscriptions     = wsBaseUrl + "subscription"
	subscriptionsById = subscriptions + "/"
	pushConfig        = wsBaseUrl + "/config"
)

// AbiosSdk defines the interface of an implementation of a SDK targeting the Abios endpoints.
type AbiosSdk interface {
	SetRate(second, minute int)
	Games(params Parameters) (GameStructPaginated, *ErrorStruct)
	Series(params Parameters) (SeriesStructPaginated, *ErrorStruct)
	SeriesById(id int, params Parameters) (SeriesStruct, *ErrorStruct)
	MatchesById(id int, params Parameters) (MatchStruct, *ErrorStruct)
	Tournaments(params Parameters) (TournamentStructPaginated, *ErrorStruct)
	TournamentsById(id int, params Parameters) (TournamentStruct, *ErrorStruct)
	SubstagesById(id int, params Parameters) (SubstageStruct, *ErrorStruct)
	Teams(params Parameters) (TeamStructPaginated, *ErrorStruct)
	TeamsById(id int, params Parameters) (TeamStruct, *ErrorStruct)
	Players(params Parameters) (PlayerStructPaginated, *ErrorStruct)
	PlayersById(id int, params Parameters) (PlayerStruct, *ErrorStruct)
	RostersById(id int, params Parameters) (RosterStruct, *ErrorStruct)
	Search(query string, params Parameters) ([]SearchResultStruct, *ErrorStruct)
	Incidents(params Parameters) (IncidentStructPaginated, *ErrorStruct)
	IncidentsBySeriesId(id int) (SeriesIncidentsStruct, *ErrorStruct)
	Organisations(params Parameters) (OrganisationStructPaginated, *ErrorStruct)
	OrganisationsById(id int, params Parameters) (OrganisationStructPaginated, *ErrorStruct)

	// PUSH API
	CreateSubscription(sub Subscription) (uuid.UUID, error)
	ListSubscriptions() ([]Subscription, error)
	// UpdateSubscription(id int, sub Subscription) (Subscription, error)
	// DeleteSubscription(id int) error
	// PushServiceConfig() ([]byte, error)
	PushServiceConnect() error
}

// client holds the oauth string returned from Authenticate as well as this sessions
// requestHandler.
type client struct {
	username       string
	password       string
	oauth          AccessTokenStruct
	handler        *requestHandler
	wsConn         *websocket.Conn
	reconnectToken uuid.UUID
}

// authenticator makes sure the oauth token doesn't expire.
func (a *client) authenticator() {
	for {
		// Wait until token is about to expire
		expires := time.Duration(a.oauth.ExpiresIn) * time.Second
		time.Sleep(expires - time.Minute*9) // Sleep until at most 9 minutes left.

		err := a.authenticate() // try once
		if err == nil {
			continue // It succeded.
		}

		// If we get an error we retry every 30 seconds for 5 minutes before we override
		// the responses.
		retry := time.NewTicker(30 * time.Second)
		fail := time.NewTimer(5 * time.Minute)

		select {
		case <-retry.C:
			err = a.authenticate()
			if err == nil {
				a.handler.override = responseOverride{override: false, data: result{}}
				break
			}
		case <-fail.C:
			a.handler.override = responseOverride{override: true, data: *err}
			break
		}
	}
}

// NewAbios returns a new endpoint-wrapper for api version 2 with given credentials.
func New(username, password string) *client {
	r := newRequestHandler()
	c := &client{username, password, AccessTokenStruct{}, r}
	err := c.authenticate()
	if err != nil {
		c.handler.override = responseOverride{override: true, data: *err}
	}
	go c.authenticator() // Launch authenticator
	return c
}

// SetRate sets the outgoing rate to "second" requests per second and "minute" requests
// per minte. A value less than or equal to 0 means previous
// value is kept. Default values are (5, 300)
func (a *client) SetRate(second, minute int) {
	a.handler.setRate(second, minute)
}

// authenticate queries the /oauth/access_token endpoint with the given credentials and
// stores the returned oauth token. Return nil if the request was successful.
func (a *client) authenticate() *result {
	var payload = []byte(`grant_type=client_credentials&client_id=` + a.username + `&client_secret=` + a.password)

	req, _ := http.NewRequest("POST", access_token, bytes.NewBuffer(payload))
	req.Header = http.Header{"Content-Type": {"application/x-www-form-urlencoded"}}

	statusCode, b := apiCall(req)
	dec := json.NewDecoder(bytes.NewBuffer(b))
	if 200 <= statusCode && statusCode < 300 {
		target := AccessTokenStruct{}
		dec.Decode(&target)
		a.oauth = target
		return nil
	} else {
		return &result{statuscode: statusCode, body: b}
	}

	return nil
}
