package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const jsonContentType = "application/json"

func TestLeague(t *testing.T) {
	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		store := StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(&store)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)

		assertLeague(t, got, wantedLeague)

		assertContentype(t, response, jsonContentType)
	})
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func getLeagueFromResponse(t *testing.T, body io.Reader) (league []Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of player, '%v'", body, err)
	}

	return
}

func assertContentype(t *testing.T, response *httptest.ResponseRecorder, want string)  {
	t.Helper()
	if response.Result().Header.Get("content-type") != "application/json" {
		t.Errorf("response did not have content-type of application/json, got %v", response.Result().Header)
	}
}

func assertLeague(t *testing.T, got, want []Player)  {
	t.Helper()

	if !reflect.DeepEqual(got, want){
		t.Errorf("got %v want %v", got, want)
	}
}