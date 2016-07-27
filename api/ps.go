package api

import "github.com/deis/controller-sdk-go/pkg/time"

// Pods defines the structure of a process.
type Pods struct {
	Release string    `json:"release"`
	Type    string    `json:"type"`
	Name    string    `json:"name"`
	State   string    `json:"state"`
	Started time.Time `json:"started"`
}

type PodsList []Pods

func (p PodsList) Len() int           { return len(p) }
func (p PodsList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PodsList) Less(i, j int) bool { return p[i].Name < p[j].Name }

type PodType struct {
	Type     string
	PodsList PodsList
}

type PodTypes []PodType

func (p PodTypes) Len() int           { return len(p) }
func (p PodTypes) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PodTypes) Less(i, j int) bool { return p[i].Type < p[j].Type }
