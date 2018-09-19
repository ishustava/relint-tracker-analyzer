package trackerclient

import (
	"time"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

const trackerURL = "https://www.pivotaltracker.com/services/v5"

type TrackerClient interface {
	GetCompletedStoriesFor(duration time.Duration) (Stories, error)
	GetStoryTransitions(storyID int) (StoryTransitions, error)
}

type trackerClient struct {
	apiToken string
	projectID string
}

func New(apiToken string, projectID string) TrackerClient {
	return trackerClient{
		apiToken: apiToken,
		projectID: projectID,
	}
}

func (tc trackerClient) GetCompletedStoriesFor(duration time.Duration) (Stories, error) {
	var stories Stories
	offset := 0

	for pageOfStories, err := tc.getAPageOfStories(duration, offset); len(pageOfStories) > 0; pageOfStories, err = tc.getAPageOfStories(duration, offset) {
		if err != nil {
			return nil, err
		}
		stories = append(stories, pageOfStories...)
		offset += len(pageOfStories)
	}

	return stories, nil
}

func (tc trackerClient) GetStoryTransitions(storyID int) (StoryTransitions, error) {
	requestURL := fmt.Sprintf("%s/projects/%s/stories/%d/transitions", trackerURL, tc.projectID, storyID)

	resp, err := tc.doGetRequest(requestURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("%s - %s", resp.Status, string(body))
	}

	var transitions StoryTransitions

	err = json.NewDecoder(resp.Body).Decode(&transitions)

	return transitions, err
}

func (tc trackerClient) getAPageOfStories(duration time.Duration, offset int) (Stories, error) {
	acceptedAfter := time.Now().Add(-duration).UnixNano()/1000000
	requestURL := fmt.Sprintf("%s/projects/%s/stories?accepted_after=%d&with_state=accepted&&offset=%d", trackerURL, tc.projectID, acceptedAfter, offset)

	resp, err := tc.doGetRequest(requestURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("%s - %s", resp.Status, string(body))
	}

	var stories Stories

	err = json.NewDecoder(resp.Body).Decode(&stories)

	return stories, err
}

func (tc trackerClient) doGetRequest(requestURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{"X-TrackerToken": []string{tc.apiToken}}

	resp, err := http.DefaultClient.Do(req)

	return resp, err
}