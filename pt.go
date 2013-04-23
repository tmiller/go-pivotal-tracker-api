/*
Package pivotaltracker implements a wrapper around Pivotal Tracker's API.
*/
package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	if story, ok := FindStory(/* AN INTEGER */); ok {
		fmt.Printf("[#%d] \n\n%s\n%s\n", story.Id, story.Name, story.Url)
	}
}

// Story represents a Pivotal Tracker story
type Story struct {
	Id   int    `xml:"id"`
	Name string `xml:"name"`
	Url  string `xml:"url"`
}

// Calls Pivotal Tracker and finds a story for the given story_id.
func FindStory(storyId int) (Story, bool) {
	findStory := fmt.Sprintf("stories/%d", storyId)

	response, err := callPivotalTracker(findStory)
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
func callPivotalTracker(command string) (response []byte, err error) {
	client := new(http.Client)

	url := "https://www.pivotaltracker.com/services/v4/" + command
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	apiKey := "PIVOTAL API KEY"
	request.Header.Add("X-TrackerToken", apiKey)

	resp, err := client.Do(request)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	return
}
