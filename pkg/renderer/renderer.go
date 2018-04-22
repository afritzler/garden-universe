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

	bg "github.com/afritzler/garden-universe/pkg/types"
	gardenclientset "github.com/gardener/gardener/pkg/client/garden/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func GetGraph(kubeconfig string) ([]byte, error) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %s", err)
	}

	gardenset, err := gardenclientset.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create garden clientset: %s", err)
	}

	nodes := make(map[string]*bg.Node)
	links := make([]bg.Link, 0)

	shoots, err := gardenset.GardenV1beta1().Shoots("").List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get shoots: %s", err)
	}
	seeds, err := gardenset.GardenV1beta1().Seeds().List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get seeds: %s", err)
	}

	// populate seeds
	for _, s := range seeds.Items {
		nodeName := fmt.Sprintf("%s/%s", "seed", s.GetObjectMeta().GetName())
		nodes[s.GetObjectMeta().GetName()] = &bg.Node{Id: s.GetObjectMeta().GetName(), Project: "", Name: nodeName, Seed: true}
	}

	// populate nodes
	for _, s := range shoots.Items {
		shootname := s.GetObjectMeta().GetName()
		nodeName := fmt.Sprintf("%s/%s", s.GetObjectMeta().GetNamespace(), s.GetObjectMeta().GetName())
		node, ok := nodes[shootname]
		status := ""
		if s.Status.LastError != nil {
			status = s.Status.LastError.Description
		}
		v := 0
		if ok {
			node.Project = s.GetObjectMeta().GetNamespace()
			v = 1
		} else {
			nodes[shootname] = &bg.Node{Id: shootname, Project: s.GetObjectMeta().GetNamespace(), Name: nodeName, Seed: false, Status: status}
		}
		links = append(links, bg.Link{Source: shootname, Target: *s.Spec.Cloud.Seed, Value: v})
	}
	data, err := json.MarshalIndent(bg.Graph{Nodes: values(nodes), Links: &links}, "", "	")
	if err != nil {
		return nil, fmt.Errorf("JSON marshaling failed: %s", err)
	}
	return data, nil
}

func values(nodes map[string]*bg.Node) *[]bg.Node {
	array := []bg.Node{}
	for _, n := range nodes {
		array = append(array, *n)
	}
	return &array
}
