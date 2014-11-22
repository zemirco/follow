package follow

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// create new document
func document(id string) {
	url := fmt.Sprintf("%s%s/%s", Url, Database, id)
	var data = []byte(`{"key": "value"}`)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
}

func TestBefore(t *testing.T) {
	// create new database
	url := fmt.Sprintf("%s%s", Url, Database)
	req, err := http.NewRequest("PUT", url, nil)
	check(err)
	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()

	// create document in database
	document("john")
}

func TestChanges(t *testing.T) {
	params := QueryParameters{}
	changes, err := Changes(params)
	if err != nil {
		t.Fatal("changes error")
	}
	if changes.LastSeq != 1 {
		t.Fatal("changes LastSeq error")
	}
	if len(changes.Results) != 1 {
		t.Fatal("changes Results error")
	}
}

func TestContinuous(t *testing.T) {
	params := QueryParameters{}
	changes, errors := Continuous(params)

	// add document in separate goroutine
	go func() {
		// wait for for loop to start listening
		time.Sleep(100 * time.Millisecond)
		document("steve")
	}()

	// start listening for changes in main goroutine
loop:
	for {
		select {
		// read from channel and test if it has been closed
		case change, ok := <-changes:
			if ok == false {
				t.Fatal("continuous error")
			}
			if change.Seq == 1 && change.Id != "john" {
				t.Fatal("continuous error")
			}
			if change.Seq == 2 && change.Id != "steve" {
				t.Fatal("continuous error")
			}
		case err := <-errors:
			t.Fatal(err)
		case <-time.After(1 * time.Second):
			break loop
		}
	}

}

func TestAfter(t *testing.T) {
	// delete database
	url := fmt.Sprintf("%s%s", Url, Database)
	req, err := http.NewRequest("DELETE", url, nil)
	check(err)
	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
}
