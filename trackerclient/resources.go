package trackerclient

import "time"

type Story struct {
	ID           int           `json:"id"`
	StoryType    string        `json:"story_type"`
	Name         string        `json:"name"`
	CurrentState string        `json:"current_state"`
	Labels       []Label       `json:"labels"`
	CycleTime    time.Duration `json:"-"`
}

type Label struct {
	Name string `json:"name"`
}

type Stories []Story

type StoryTransition struct {
	State      string `json:"state"`
	OccurredAt string `json:"occurred_at"`
}

type StoryTransitions []StoryTransition

func (s Story) HasLabel(label string) bool {
	for _, l := range s.Labels {
		if l.Name == label {
			return true
		}
	}

	return false
}

func (s Story) HasALabelFrom(labels []string) bool {
	for _, l := range labels {
		if s.HasLabel(l) {
			return true
		}
	}

	return false
}

