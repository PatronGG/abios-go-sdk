package abios

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	. "github.com/PatronGG/abios-go-sdk/structs"
	"github.com/gobuffalo/uuid"
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
}

// client holds the oauth string returned from Authenticate as well as this sessions
// requestHandler.
type client struct {
	username string
	password string
	oauth    AccessTokenStruct
	handler  *requestHandler
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

// Games queries the /games endpoint and returns a GameStructPaginated.
func (a *client) Games(params Parameters) (GameStructPaginated, *ErrorStruct) {
	if params == nil {
		params = make(Parameters)
	}
	params.Set("access_token", a.oauth.AccessToken)
	result := <-a.handler.addRequest(games, params)

	dec := json.NewDecoder(bytes.NewBuffer(result.body))
	if 200 <= result.statuscode && result.statuscode < 300 {
		target := GameStructPaginated{}
		dec.Decode(&target)
		return target, nil
	} else {
		target := ErrorStruct{}
		dec.Decode(target)
		return GameStructPaginated{}, &target
	}

	return GameStructPaginated{}, &ErrorStruct{}
}

// Series queries the /series endpoint and returns a SeriesStructPaginated.
func (a *client) Series(params Parameters) (SeriesStructPaginated, *ErrorStruct) {
	if params == nil {
		params = make(Parameters)
	}
	params.Set("access_token", a.oauth.AccessToken)
	result := <-a.handler.addRequest(series, params)

	dec := json.NewDecoder(bytes.NewBuffer(result.body))
	if 200 <= result.statuscode && result.statuscode < 300 {
		target := SeriesStructPaginated{}
		dec.Decode(&target)
		return target, nil
	} else {
		target := ErrorStruct{}
		dec.Decode(&target)
		return SeriesStructPaginated{}, &target
	}

	return SeriesStructPaginated{}, &ErrorStruct{}
}

// SeriesById queries the /series/:id endpoint and returns a SeriesStruct.
func (a *client) SeriesById(id int, params Parameters) (SeriesStruct, *ErrorStruct) {
	sId := strconv.Itoa(id)
	if params == nil {
		params = make(Parameters)
	}
	params.Set("access_token", a.oauth.AccessToken)
	result := <-a.handler.addRequest(seriesById+sId, params)

	dec := json.NewDecoder(bytes.NewBuffer(result.body))
	if 200 <= result.statuscode && result.statuscode < 300 {
		target := SeriesStruct{}
		dec.Decode(&target)
		return target, nil
	} else {
		target := ErrorStruct{}
		dec.Decode(&target)
		return SeriesStruct{}, &target
	}

	return SeriesStruct{}, &ErrorStruct{}
}

// MatchesById queries the /matches/:id endpoint and returns a MatchStruct.
func (a *client) MatchesById(id int, params Parameters) (MatchStruct, *ErrorStruct) {
	sId := strconv.Itoa(id)
	if params == nil {
		params = make(Parameters)
	}
	params.Set("access_token", a.oauth.AccessToken)
	result := <-a.handler.addRequest(matches+sId, params)

	dec := json.NewDecoder(bytes.NewBuffer(result.body))
	if 200 <= result.statuscode && result.statuscode < 300 {
		target := MatchStruct{}
		dec.Decode(&target)
		return target, nil
	} else {
		target := ErrorStruct{}
		dec.Decode(&target)
		return MatchStruct{}, &target
	}

	return MatchStruct{}, &ErrorStruct{}
}

// Tournaments queries the /tournaments endpoint and returns a list of TournamentStructPaginated.
func (a *client) Tournaments(params Parameters) (TournamentStructPaginated, *ErrorStruct) {
	if params == nil {
		params = make(Parameters)
	}
	params.Set("access_token", a.oauth.AccessToken)
	result := <-a.handler.addRequest(tournaments, params)

	dec := json.NewDecoder(bytes.NewBuffer(result.body))
	if 200 <= result.statuscode && result.statuscode < 300 {
		target := TournamentStructPaginated{}
		dec.Decode(&target)
		return target, nil
	} else {
		target := ErrorStruct{}
		dec.Decode(&target)
		return TournamentStructPaginated{}, &target
	}

	return TournamentStructPaginated{}, &ErrorStruct{}
}

// TournamentsById queries the /tournaments/:id endpoint and return a TournamentStruct.
func (a *client) TournamentsById(id int, params Parameters) (TournamentStruct, *ErrorStruct) {
	sId := strconv.Itoa(id)
	if params == nil {
		params = make(Parameters)
	}
	params.Set("access_token", a.oauth.AccessToken)
	result := <-a.handler.addRequest(tournamentsById+sId, params)

	dec := json.NewDecoder(bytes.NewBuffer(result.body))
	if 200 <= result.statuscode && result.statuscode < 300 {
		target := TournamentStruct{}
		dec.Decode(&target)
		return target, nil
	} else {
		target := ErrorStruct{}
		dec.Decode(&target)
		return TournamentStruct{}, &target
	}

	return TournamentStruct{}, &ErrorStruct{}
}

// SubstagesById queries the /substages/:id endpoint and returns a SubstageStruct.
func (a *client) SubstagesById(id int, params Parameters) (SubstageStruct, *ErrorStruct) {
	sId := strconv.Itoa(id)
	if params == nil {
		params = make(Parameters)
	}
	params.Set("access_token", a.oauth.AccessToken)
	result := <-a.handler.addRequest(substages+sId, params)

	dec := json.NewDecoder(bytes.NewBuffer(result.body))
	if 200 <= result.statuscode && result.statuscode < 300 {
		target := SubstageStruct{}
		dec.Decode(&target)
		return target, nil
	} else {
		target := ErrorStruct{}
		dec.Decode(&target)
		return SubstageStruct{}, &target
	}

	return SubstageStruct{}, &ErrorStruct{}
}

// Teams queries the /teams endpoint and returns a TeamsStructPaginated.
func (a *client) Teams(params Parameters) (TeamStructPaginated, *ErrorStruct) {
	if params == nil {
		params = make(Parameters)
	}
	params.Set("access_token", a.oauth.AccessToken)
	result := <-a.handler.addRequest(teams, params)

	dec := json.NewDecoder(bytes.NewBuffer(result.body))
	if 200 <= result.statuscode && result.statuscode < 300 {
		target := TeamStructPaginated{}
		dec.Decode(&target)
		return target, nil
	} else {
		target := ErrorStruct{}
		dec.Decode(&target)
		return TeamStructPaginated{}, &target
	}

	return TeamStructPaginated{}, &ErrorStruct{}
}

// TeamsById queues the /teams/:id endpoint and return a TeamStruct.
func (a *client) TeamsById(id int, params Parameters) (TeamStruct, *ErrorStruct) {
	sId := strconv.Itoa(id)
	if params == nil {
		params = make(Parameters)
	}
	params.Set("access_token", a.oauth.AccessToken)
	result := <-a.handler.addRequest(teamsById+sId, params)

	dec := json.NewDecoder(bytes.NewBuffer(result.body))
	if 200 <= result.statuscode && result.statuscode < 300 {
		target := TeamStruct{}
		dec.Decode(&target)
		return target, nil
	} else {
		target := ErrorStruct{}
		dec.Decode(&target)
		return TeamStruct{}, &target
	}

	return TeamStruct{}, &ErrorStruct{}
}

// Players queries the /players endpoint and returns PlayerStructPaginated.
func (a *client) Players(params Parameters) (PlayerStructPaginated, *ErrorStruct) {
	if params == nil {
		params = make(Parameters)
	}
	params.Set("access_token", a.oauth.AccessToken)
	result := <-a.handler.addRequest(players, params)

	dec := json.NewDecoder(bytes.NewBuffer(result.body))
	if 200 <= result.statuscode && result.statuscode < 300 {
		target := PlayerStructPaginated{}
		dec.Decode(&target)
		return target, nil
	} else {
		target := ErrorStruct{}
		dec.Decode(&target)
		return PlayerStructPaginated{}, &target
	}

	return PlayerStructPaginated{}, &ErrorStruct{}
}

// PlayersById queries the /players/:id endpoint and returns a PlayerStruct.
func (a *client) PlayersById(id int, params Parameters) (PlayerStruct, *ErrorStruct) {
	sId := strconv.Itoa(id)
	if params == nil {
		params = make(Parameters)
	}
	params.Set("access_token", a.oauth.AccessToken)
	result := <-a.handler.addRequest(playersById+sId, params)

	dec := json.NewDecoder(bytes.NewBuffer(result.body))
	if 200 <= result.statuscode && result.statuscode < 300 {
		target := PlayerStruct{}
		dec.Decode(&target)
		return target, nil
	} else {
		target := ErrorStruct{}
		dec.Decode(&target)
		return PlayerStruct{}, &target
	}

	return PlayerStruct{}, &ErrorStruct{}
}

// RostersById queries the /rosters/:id endpoint and returns a RosterStruct.
func (a *client) RostersById(id int, params Parameters) (RosterStruct, *ErrorStruct) {
	sId := strconv.Itoa(id)
	if params == nil {
		params = make(Parameters)
	}
	params.Set("access_token", a.oauth.AccessToken)
	result := <-a.handler.addRequest(rosters+sId, params)

	dec := json.NewDecoder(bytes.NewBuffer(result.body))
	if 200 <= result.statuscode && result.statuscode < 300 {
		target := RosterStruct{}
		dec.Decode(&target)
		return target, nil
	} else {
		target := ErrorStruct{}
		dec.Decode(&target)
		return RosterStruct{}, &target
	}

	return RosterStruct{}, &ErrorStruct{}
}

// Search queries the /search endpoint with the given query and returns a list of
// SearchResultStruct.
func (a *client) Search(query string, params Parameters) ([]SearchResultStruct, *ErrorStruct) {
	params.Set("access_token", a.oauth.AccessToken)
	params.Add("q", query)
	result := <-a.handler.addRequest(search, params)

	dec := json.NewDecoder(bytes.NewBuffer(result.body))
	if 200 <= result.statuscode && result.statuscode < 300 {
		target := []SearchResultStruct{}
		dec.Decode(&target)
		return target, nil
	} else {
		target := ErrorStruct{}
		dec.Decode(&target)
		return []SearchResultStruct{}, &target
	}

	return []SearchResultStruct{}, &ErrorStruct{}
}

// Incidents queries the /incidents endpoint and returns an IncidentStructPaginated.
func (a *client) Incidents(params Parameters) (IncidentStructPaginated, *ErrorStruct) {
	if params == nil {
		params = make(Parameters)
	}
	params.Set("access_token", a.oauth.AccessToken)
	result := <-a.handler.addRequest(incidents, params)

	dec := json.NewDecoder(bytes.NewBuffer(result.body))
	if 200 <= result.statuscode && result.statuscode < 300 {
		target := IncidentStructPaginated{}
		dec.Decode(&target)
		return target, nil
	} else {
		target := ErrorStruct{}
		dec.Decode(&target)
		return IncidentStructPaginated{}, &target
	}

	return IncidentStructPaginated{}, &ErrorStruct{}
}

// IncidentBySeriesId queries the /incidents/:series_id endpoint and returns a
// SeriesIncidentsStruct.
func (a *client) IncidentsBySeriesId(id int) (SeriesIncidentsStruct, *ErrorStruct) {
	sId := strconv.Itoa(id)
	params := make(Parameters)
	params.Set("access_token", a.oauth.AccessToken)
	result := <-a.handler.addRequest(incidentsBySeries+sId, params)

	dec := json.NewDecoder(bytes.NewBuffer(result.body))
	if 200 <= result.statuscode && result.statuscode < 300 {
		target := SeriesIncidentsStruct{}
		dec.Decode(&target)
		return target, nil
	} else {
		target := ErrorStruct{}
		dec.Decode(&target)
		return SeriesIncidentsStruct{}, &target
	}

	return SeriesIncidentsStruct{}, &ErrorStruct{}
}

// Organisations queries the /organisations endpoint and returns a OrganisationStructPaginated
func (a *client) Organisations(params Parameters) (OrganisationStructPaginated, *ErrorStruct) {
	if params == nil {
		params = make(Parameters)
	}
	params.Set("access_token", a.oauth.AccessToken)
	result := <-a.handler.addRequest(organisations, params)

	dec := json.NewDecoder(bytes.NewBuffer(result.body))
	if 200 <= result.statuscode && result.statuscode < 300 {
		target := OrganisationStructPaginated{}
		dec.Decode(&target)
		return target, nil
	} else {
		target := ErrorStruct{}
		dec.Decode(&target)
		return OrganisationStructPaginated{}, &target
	}

	return OrganisationStructPaginated{}, &ErrorStruct{}
}

// OrganisationsById queues the /organisations/:id endpoint and return a OrganisationStruct.
func (a *client) OrganisationsById(id int, params Parameters) (OrganisationStruct, *ErrorStruct) {
	sId := strconv.Itoa(id)
	if params == nil {
		params = make(Parameters)
	}
	params.Set("access_token", a.oauth.AccessToken)
	result := <-a.handler.addRequest(organisationsById+sId, params)

	dec := json.NewDecoder(bytes.NewBuffer(result.body))
	if 200 <= result.statuscode && result.statuscode < 300 {
		target := OrganisationStruct{}
		dec.Decode(&target)
		return target, nil
	} else {
		target := ErrorStruct{}
		dec.Decode(&target)
		return OrganisationStruct{}, &target
	}

	return OrganisationStruct{}, &ErrorStruct{}
}

func (a *client) CreateSubscription(sub Subscription) (uuid.UUID, error) {
	params := make(Parameters)
	params.Set("access_token", a.oauth.AccessToken)

	u, err := url.Parse(subscriptions)
	if err != nil {
		return uuid.Nil, err
	}
	u.RawQuery = params.encode()

	subStr, err := json.Marshal(sub)
	if err != nil {
		return uuid.Nil, err
	}
	res, err := http.Post(u.String(), "application/json", bytes.NewBuffer(subStr))
	defer res.Body.Close()

	if err != nil {
		dec := json.NewDecoder(res.Body)
		if res.StatusCode == http.StatusOK {
			s := Message{}
			err = dec.Decode(&s)
			return s.UUID, nil
		} else if res.StatusCode == http.StatusUnprocessableEntity {
			var existingID uuid.UUID

			if res.Header.Get("Location") != "" {
				existingID, err = uuid.FromString(res.Header.Get("Location"))
				if err != nil {
					return uuid.Nil, err
				}

				return existingID, nil
			}
		}
		return uuid.Nil, fmt.Errorf("Unexpected status code %v", res.StatusCode)
	}

	return uuid.Nil, err
}

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
