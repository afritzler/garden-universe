// Copyright © 2018 Andreas Fritzler <andreas.fritzler@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package stats

import (
	"encoding/json"
	"fmt"

	"github.com/afritzler/garden-universe/pkg/gardener"
	bg "github.com/afritzler/garden-universe/pkg/types"
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
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
	var noOfNodes int32
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

func GetSizeOfShoot(shoot v1beta1.Shoot) int32 {
	var size int32
	for _, w := range shoot.Spec.Provider.Workers {
		size += w.Maximum
	}
	return size
}
