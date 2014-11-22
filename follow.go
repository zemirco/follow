// Package follow implements CouchDB _changes API.
// For more info see http://docs.couchdb.org/en/latest/api/database/changes.html.
//
// 	follow.Url = "http://127.0.0.1:5984/"
// 	follow.Database = "_users"
//
// 	params := follow.QueryParameters{
// 		Limit: 10,
// 	}
//
// 	changes, err := follow.Changes(params)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(changes)
package follow

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"net/http"
)

var (
	Url      = "http://127.0.0.1:5984/"
	Database = "test"
)

type QueryParameters struct {
	DocIds          []string `url:"doc_ids,omitempty"`
	Conflicts       bool     `url:"conflicts,omitempty"`
	Descending      bool     `url:"descending,omitempty"`
	Feed            string   `url:"feed,omitempty"`
	Filter          string   `url:"filter,omitempty"`
	Heartbeat       int64    `url:"heartbeat,omitempty"`
	IncludeDocs     bool     `url:"include_docs,omitempty"`
	Attachments     bool     `url:"attachments,omitempty"`
	AttEncodingInfo bool     `url:"att_encoding_info,omitempty"`
	LastEventId     int64    `url:"last-event-id,omitempty"`
	Limit           int      `url:"limit,omitempty"`
	Since           int      `url:"since,omitempty"`
	Style           string   `url:"style,omitempty"`
	Timeout         int64    `url:"timeout,omitempty"`
	View            string   `url:"view,omitempty"`
}

type Response struct {
	LastSeq int      `json:"last_seq,omitempty"`
	Results []Result `json:"results,omitempty"`
}

type Result struct {
	Changes []Rev  `json:"changes,omitempty"`
	Id      string `json:"id,omitempty"`
	Seq     int    `json:"seq,omitempty"`
	Deleted bool   `json:"deleted,omitempty"`
}

type Rev struct {
	Rev string `json:"rev,omitempty"`
}

// Changes returns all changes immediately.
func Changes(params QueryParameters) (*Response, error) {
	q, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s%s/_changes?%s", Url, Database, q.Encode())
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	response := &Response{}
	return response, json.NewDecoder(resp.Body).Decode(&response)
}

// Continuous creates a continuous feed that stays open
// and connected to the database.
func Continuous(params QueryParameters) (<-chan *Result, <-chan error) {

	// use continuous feed
	params.Feed = "continuous"
	q, err := query.Values(params)
	if err != nil {
		return nil, nil
	}

	// create channels
	result := make(chan *Result)
	e := make(chan error)

	url := fmt.Sprintf("%s%s/_changes?%s", Url, Database, q.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		e <- err
		return nil, e
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		e <- err
		return nil, e
	}
	dec := json.NewDecoder(res.Body)

	// start goroutine
	go func() {
		defer close(result)
		defer res.Body.Close()
		for {
			var r Result
			if err := dec.Decode(&r); err != nil {
				e <- err
				return
			}
			result <- &r
		}
	}()

	return result, e
}
