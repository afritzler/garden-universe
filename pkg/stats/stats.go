package stats

import (
	"encoding/json"
	"fmt"

	"github.com/afritzler/garden-universe/pkg/gardener"
	bg "github.com/afritzler/garden-universe/pkg/types"
	"github.com/gardener/gardener/pkg/apis/garden/v1beta1"
)

type Stats interface {
	GetStats() (*bg.Stats, error)
	GetStatsJSON() ([]byte, error)
}

type stats struct {
	garden gardener.Gardener
}

func NewStats(garden gardener.Gardener) Stats {
	return &stats{garden: garden}
}

func (s *stats) GetStats() (*bg.Stats, error) {
	shoots, err := s.garden.GetShootList()
	if err != nil {
		return nil, err
	}
	noOfNodes := 0
	for _, shoot := range shoots.Items {
		noOfNodes += GetSizeOfShoot(shoot)
	}
	return &bg.Stats{NoOfShoots: len(shoots.Items), NoOfNodes: noOfNodes}, nil
}

func (s *stats) GetStatsJSON() ([]byte, error) {
	data, err := s.GetStats()
	if err != nil {
		return nil, err
	}
	json, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return nil, fmt.Errorf("JSON marshaling of stats object failed: %s", err)
	}
	return json, nil
}

func GetSizeOfShoot(shoot v1beta1.Shoot) int {
	size := 0
	if shoot.Spec.Cloud.OpenStack != nil {
		for _, w := range shoot.Spec.Cloud.OpenStack.Workers {
			size += w.AutoScalerMin
		}
	}
	if shoot.Spec.Cloud.AWS != nil {
		for _, w := range shoot.Spec.Cloud.AWS.Workers {
			size += w.AutoScalerMin
		}
	}
	if shoot.Spec.Cloud.Azure != nil {
		for _, w := range shoot.Spec.Cloud.Azure.Workers {
			size += w.AutoScalerMin
		}
	}
	if shoot.Spec.Cloud.GCP != nil {
		for _, w := range shoot.Spec.Cloud.GCP.Workers {
			size += w.AutoScalerMin
		}
	}
	return size
}
