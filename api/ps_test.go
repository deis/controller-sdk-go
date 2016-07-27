package api

import (
	"sort"
	"testing"

	"github.com/deis/controller-sdk-go/pkg/time"
)

func TestPodsListSorted(t *testing.T) {
	pods := PodsList{
		{"", "web", "web.fsdfgh4", "up", time.Time{}},
		{"", "web", "web.asdfgh1", "up", time.Time{}},
		{"", "web", "web.csdfgh3", "up", time.Time{}},
		{"", "web", "web.bsdfgh2", "up", time.Time{}},
	}

	sort.Sort(pods)

	expectedPodNames := []string{"web.asdfgh1", "web.bsdfgh2", "web.csdfgh3", "web.fsdfgh4"}

	for i, pod := range pods {
		if expectedPodNames[i] != pod.Name {
			t.Errorf("Expected pods to be sorted %v, Got %v", expectedPodNames[i], pod.Name)
		}
	}
}

func TestPodTypesSorted(t *testing.T) {
	podTypes := PodTypes{
		{"worker", PodsList{}},
		{"web", PodsList{}},
		{"clock", PodsList{}},
	}

	sort.Sort(podTypes)
	expectedPodTypes := []string{"clock", "web", "worker"}

	for i, podType := range podTypes {
		if expectedPodTypes[i] != podType.Type {
			t.Errorf("Expected pod types to be sorted %v, Got %v at index %v", expectedPodTypes[i], podType.Type, i)
		}
	}
}
