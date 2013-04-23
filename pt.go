/*
Package pt implements a wrapper around Pivotal Tracker's API.
*/
package pt

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Story represents a Pivotal Tracker story
type Story struct {
	Id           int    `xml:"id"`
	Name         string `xml:"name"`
	Url          string `xml:"url"`
	CurrentState string `xml:"current_state"`
}

func (s Story) State() string {
	return strings.Title(s.CurrentState)
}

// PivotalTracker holds state information about the API.
type PivotalTracker struct {
	ApiKey string
}

// Calls Pivotal Tracker and finds a story for the given story_id.
func (pt PivotalTracker) FindStory(storyId string) (story Story, err error) {
	findStory := fmt.Sprintf("stories/%s", storyId)

	response, err := pt.callPivotalTracker(findStory)
	if err != nil {
		return
	}

	xml.Unmarshal(response, &story)

	if story.Id == 0 {
		err = errors.New("No Story found for " + storyId + ".")
	}

	return
}

// Sends a command to Pivotal Tracker and returns XML representation of the
// response.
func (pt PivotalTracker) callPivotalTracker(command string) (response []byte, err error) {
	client := new(http.Client)

	url := "https://www.pivotaltracker.com/services/v4/" + command
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	request.Header.Add("X-TrackerToken", pt.ApiKey)

	resp, err := client.Do(request)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)

	return
}
