/*
Package pivotaltracker implements a wrapper around Pivotal Tracker's API.
*/
package pt

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Story represents a Pivotal Tracker story
type Story struct {
	Id   int    `xml:"id"`
	Name string `xml:"name"`
	Url  string `xml:"url"`
}

// PivotalTracker holds state information about the API.
type PivotalTracker struct {
  ApiKey string
}


// Calls Pivotal Tracker and finds a story for the given story_id.
func (pt PivotalTracker)FindStory(storyId int) (Story, bool) {
	findStory := fmt.Sprintf("stories/%d", storyId)

	response, err := pt.callPivotalTracker(findStory)
	if err != nil {
		fmt.Println(err)
	}

	var story Story
	xml.Unmarshal(response, &story)

	found := (story.Id != 0)
	return story, found
}

// Sends a command to Pivotal Tracker and returns XML representation of the
// response.
func (pt PivotalTracker)callPivotalTracker(command string) (response []byte, err error) {
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
