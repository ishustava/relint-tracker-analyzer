package main

import (
	"github.com/ishustava/relint-tracker-analyzer/trackerclient"
	"os"
	"time"
	"fmt"
	"strings"
	"github.com/ishustava/relint-tracker-analyzer/worktime"
	"github.com/ishustava/relint-tracker-analyzer/stats"
)

const ThreeWeeks = 504 * time.Hour
const SixWeeks = 2 * ThreeWeeks
const NineWeeks = 3 * ThreeWeeks

func main() {
	apiToken := os.Getenv("TRACKER_API_TOKEN")
	projectID := os.Getenv("PROJECT_ID")

	if apiToken == "" || projectID == "" {
		fmt.Println("TRACKER_API_TOKEN and PROJECT_ID are required")
		os.Exit(1)
	}

	trackerClient := trackerclient.New(apiToken, projectID)

	printSummaryFor(ThreeWeeks, 3, trackerClient)
	//printSummaryFor(SixWeeks, 6, trackerClient)
	//printSummaryFor(NineWeeks, 9, trackerClient)
}

func printSummaryFor(duration time.Duration, numWeeks int, trackerClient trackerclient.TrackerClient) {
	stories, err := trackerClient.GetCompletedStoriesFor(duration)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	stories, err = setStoryCycleTime(trackerClient, stories)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	brokenBuilds := getStoriesWithLabel(stories, "broken build")
	pullRequests := getStoriesWithLabel(stories, "github-pull-request")
	issues := getStoriesWithLabel(stories, "github-issue")
	originalFeatures := getFeaturesWithoutLabels(stories, "broken build", "gcp-502s", "github-pull-request", "github-issue")
	fmt.Println(issues)
	fmt.Println(originalFeatures)

	fmt.Printf("LAST %d WEEKS\t\t\tNumber\tPercentage\t\tTotal worktime\t\t\tAve worktime per story\n", numWeeks)
	fmt.Println("-----------------------------------------------------------------------------------------------------------------------------")
	fmt.Printf("All stories: \t\t\t%d\t%0.2f%%\t\t\t%s\t\t\t%s\n", len(stories), stats.Percent(len(stories), len(stories)), totalStoryCycleTime(stories), aveStoryCycleTime(stories))
	fmt.Printf("Broken builds: \t\t\t%d\t%0.2f%%\t\t\t%s\t\t\t%s\n", len(brokenBuilds), stats.Percent(len(brokenBuilds), len(stories)), totalStoryCycleTime(brokenBuilds), aveStoryCycleTime(brokenBuilds))
	fmt.Printf("Pull requests: \t\t\t%d\t%0.2f%%\t\t\t%s\t\t\t%s\n", len(pullRequests), stats.Percent(len(pullRequests), len(stories)), totalStoryCycleTime(pullRequests), aveStoryCycleTime(pullRequests))
	fmt.Printf("GitHub issues: \t\t\t%d\t%0.2f%%\t\t\t%s\t\t\t%s\n", len(issues), stats.Percent(len(issues), len(stories)), totalStoryCycleTime(issues), aveStoryCycleTime(issues))
	fmt.Printf("Original features: \t\t%d\t%0.2f%%\t\t\t%s\t\t\t%s\n", len(originalFeatures), stats.Percent(len(originalFeatures), len(stories)), totalStoryCycleTime(originalFeatures), aveStoryCycleTime(originalFeatures))
	fmt.Println()
}

func totalStoryCycleTime(stories trackerclient.Stories) time.Duration {
	var total time.Duration

	for _, s := range stories {
		total += s.CycleTime
	}

	return total
}

func aveStoryCycleTime(stories trackerclient.Stories) time.Duration {
	return totalStoryCycleTime(stories)/time.Duration(len(stories))
}
func getFeaturesWithoutLabels(stories trackerclient.Stories, labels ...string) trackerclient.Stories {
	var result trackerclient.Stories

	for _, story := range stories {
		if !story.HasALabelFrom(labels) && story.StoryType == "feature" && !strings.Contains(story.Name, "Triumphant Herald") {
			result = append(result, story)
		}
	}

	return result
}

func getStoriesWithLabel(allStories trackerclient.Stories, label string) trackerclient.Stories {
	var result trackerclient.Stories

	for _, story := range allStories {
		if story.HasLabel(label) {
			result = append(result, story)
		}
	}
	return result
}

func setStoryCycleTime(client trackerclient.TrackerClient, allStories trackerclient.Stories) (trackerclient.Stories, error) {
	for i, story := range allStories {
		transitions, err := client.GetStoryTransitions(story.ID)
		if err != nil {
			return nil, err
		}

		cycleTime, err := computeCycleTimeForStory(story, transitions)
		if err != nil {
			return nil, err
		}

		allStories[i].CycleTime = cycleTime
	}
	return allStories, nil
}

func computeCycleTimeForStory(story trackerclient.Story, storyTransitions trackerclient.StoryTransitions) (time.Duration, error) {
	var startTime time.Time
	var cycleTime time.Duration

	started := false

	for _, tr := range storyTransitions {
		transitionTime, err := time.Parse(time.RFC3339, tr.OccurredAt)
		if err != nil {
			return 0, err
		}
		choreFinished := tr.State == "accepted" && story.StoryType == "chore"
		featureFinished := tr.State == "finished" && story.StoryType != "chore"
		storyFinished := choreFinished || featureFinished

		if tr.State == "unscheduled" && started {
			started = false
		}

		if tr.State == "started" {
			startTime = transitionTime
			started = true
		} else if (storyFinished && !startTime.IsZero()) || (tr.State == "unstarted" && started) {
			transitionDuration, err := worktime.Duration(startTime, transitionTime)

			switch err.(type) {
			case worktime.StartAfterEnd:
				return 0, err
			case worktime.StartOnWeekend, worktime.EndOnWeekend:
			default:
				cycleTime += transitionDuration
			}
		}
	}
	if story.ID == 160005662 {
		fmt.Println(storyTransitions)
	}

	return cycleTime, nil
}