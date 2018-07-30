package abios

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gobuffalo/uuid"
)

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
