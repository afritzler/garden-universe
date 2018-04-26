// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

package pkg

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/afritzler/garden-universe/pkg/gardener"
	"github.com/afritzler/garden-universe/pkg/stats"
	bg "github.com/afritzler/garden-universe/pkg/types"
	"github.com/gardener/gardener/pkg/apis/garden"
)

type Renderer interface {
	GetGraph() ([]byte, error)
}

type renderer struct {
	garden gardener.Gardener
}

func NewRenderer(garden gardener.Gardener) Renderer {
	return &renderer{garden: garden}
}

func (r *renderer) GetGraph() ([]byte, error) {
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
			if strings.Split(r.APIVersion, "/")[0] == garden.GroupName && r.Kind == "Shoot" {
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
		if s.Status.LastError != nil {
			status = s.Status.LastError.Description
		}
		v := 0
		if ok {
			node.Project = namespace
			v = 1
		} else {
			nodes[shootname] = &bg.Node{Id: shootname, Project: namespace, Name: shootname, Type: "shoot", Status: status, Size: size}
		}
		projectMeta := *s.Spec.Cloud.Seed + namespace
		nodes[projectMeta] = &bg.Node{Id: projectMeta, Project: namespace, Name: namespace, Type: "project"}
		links = append(links, bg.Link{Source: shootname, Target: projectMeta, Value: v})
		links = append(links, bg.Link{Source: projectMeta, Target: seednames[*s.Spec.Cloud.Seed], Value: v})
	}
	data, err := json.MarshalIndent(bg.Graph{Nodes: r.values(nodes), Links: &links}, "", "	")
	if err != nil {
		return nil, fmt.Errorf("JSON marshaling failed: %s", err)
	}
	return data, nil
}

func (r *renderer) values(nodes map[string]*bg.Node) *[]bg.Node {
	array := []bg.Node{}
	for _, n := range nodes {
		array = append(array, *n)
	}
	return &array
}
