package api

import "github.com/deis/controller-sdk-go/pkg/time"

// ProcessType represents the key/value mappings of a process type to a process inside
// a Heroku Procfile.
//
// See https://devcenter.heroku.com/articles/procfile
type ProcessType map[string]string

// Pods defines the structure of a process.
type Pods struct {
	Release string    `json:"release"`
	Type    string    `json:"type"`
	Name    string    `json:"name"`
	State   string    `json:"state"`
	Started time.Time `json:"started"`
}

// PodsList defines a collection of app pods.
type PodsList []Pods

func (p PodsList) Len() int           { return len(p) }
func (p PodsList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PodsList) Less(i, j int) bool { return p[i].Name < p[j].Name }

// PodType holds pods of the same type.
type PodType struct {
	Type     string
	PodsList PodsList
}

// PodTypes holds groups of pods organized by type.
type PodTypes []PodType

func (p PodTypes) Len() int           { return len(p) }
func (p PodTypes) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PodTypes) Less(i, j int) bool { return p[i].Type < p[j].Type }
