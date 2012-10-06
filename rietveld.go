package main

import (
	// "os"
	// "io"
	"log"
	"net/http"
	"encoding/json"
	"fmt"
	"errors"
	"time"
)

type Issue struct {
	Id          uint `json:"issue"`
	Owner       string
	OwnerEmail  string `json:"owner_email"`
	Reviewers   []string
	Cc          []string
	Subject     string
	Description string
	BaseUrl     string `json:"base_url"`
	PatchsetIds []uint `json:"patchsets"`
	Private     bool
	Closed      bool
	Created     string
	Modified    string
}

func (i Issue) String() string {
	return fmt.Sprintf(
		"[%d] %s\nBase URL: %s\nOwner: %s (%s)\nReviewers: %s\n" +
		"Private: %t\nClosed: %t\nUpdated: %s\n\n",
		i.Id, i.Subject, i.BaseUrl, i.Owner, i.OwnerEmail, i.Reviewers,
		i.Private, i.Closed, i.Modified)
}

type Response struct {
	Cursor string
	Issues []Issue `json:"results"`
}

func (r Response) String() string {
	return fmt.Sprintf("Issues count: %d, Cursor: %s", len(r.Issues), r.Cursor)
}

func Search() (r *Response, e error) {
	url := "https://codereview.appspot.com/search?format=json&limit=10"	
	resp, err := http.Get(url)
	if err != nil {
		e = err
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		e = errors.New(resp.Status)
		return
	}

	r = new(Response)
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		e = err
		return
	}

	return
}

func main() {
	resp, err := Search()
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(resp)
	for _, issue := range resp.Issues {
		fmt.Println(issue)
	}
}
