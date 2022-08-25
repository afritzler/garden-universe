// Copyright © 2018 Andreas Fritzler
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package renderer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/afritzler/garden-universe/pkg/gardener"
	"github.com/afritzler/garden-universe/pkg/stats"
	bg "github.com/afritzler/garden-universe/pkg/types"
	"github.com/gardener/gardener/pkg/apis/core"
)

func NewKubeRenderer(garden gardener.Gardener) Renderer {
	return &kuberenderer{garden: garden}
}

func (r *kuberenderer) GetGraph() ([]byte, error) {
	nodes := make(map[string]*bg.Node)
	links := make([]bg.Link, 0)
	shoots, err := r.garden.GetShoots()
	if err != nil {
		return nil, fmt.Errorf("failed to get shoots: %s", err)
	}
	seeds, err := r.garden.GetSeeds()
	if err != nil {
		return nil, fmt.Errorf("failed to get seeds: %s", err)
	}
	seednames := map[string]string{}
	// populate seeds
	for _, s := range *seeds {
		refs := s.GetOwnerReferences()
		name := "seed/" + s.GetObjectMeta().GetName()
		for _, r := range refs {
			if strings.Split(r.APIVersion, "/")[0] == core.GroupName && r.Kind == "Shoot" {
				name = "garden" + "/" + r.Name
			}
		}
		nodes[name] = &bg.Node{Id: name, Project: "", Name: name, Type: "seed"}
		seednames[s.GetObjectMeta().GetName()] = name
	}

	// populate nodes and links
	for _, s := range *shoots {
		namespace := s.GetObjectMeta().GetNamespace()
		shootname := fmt.Sprintf("%s/%s", namespace, s.GetObjectMeta().GetName())
		node, ok := nodes[shootname]
		status := ""
		size := stats.GetSizeOfShoot(s)
		if len(s.Status.LastErrors) > 0 {
			status = s.Status.LastErrors[0].Description
		}
		v := 0
		if ok {
			node.Project = namespace
			v = 1
		} else {
			nodes[shootname] = &bg.Node{Id: shootname, Project: namespace, Name: shootname, Type: "shoot", Status: status, Size: size}
		}
		seedName := "core"
		if s.Spec.SeedName != nil {
			seedName = *s.Spec.SeedName
		}
		projectMeta := seedName + namespace
		nodes[projectMeta] = &bg.Node{Id: projectMeta, Project: namespace, Name: namespace, Type: "project"}
		links = append(links, bg.Link{Source: shootname, Target: projectMeta, Value: v})
		links = append(links, bg.Link{Source: projectMeta, Target: seednames[seedName], Value: v})
	}
	data, err := json.MarshalIndent(bg.Graph{Nodes: r.values(nodes), Links: &links}, "", "	")
	if err != nil {
		return nil, fmt.Errorf("JSON marshaling failed: %s", err)
	}
	return data, nil
}

func (r *kuberenderer) values(nodes map[string]*bg.Node) *[]bg.Node {
	array := []bg.Node{}
	for _, n := range nodes {
		array = append(array, *n)
	}
	return &array
}
