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
	"log"

	bg "github.com/afritzler/garden-universe/pkg/types"
	gardenclientset "github.com/gardener/gardener/pkg/client/garden/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func GetGraph(kubeconfig string) ([]byte, error) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	gardenset, err := gardenclientset.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	nodes := make(map[string]*bg.Node)
	links := make([]bg.Link, 0)

	shoots, err := gardenset.GardenV1beta1().Shoots("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	seeds, err := gardenset.GardenV1beta1().Seeds().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	// populate seeds
	for _, s := range seeds.Items {
		nodes[s.GetObjectMeta().GetName()] = &bg.Node{Id: s.GetObjectMeta().GetName(), Project: "", Seed: true}
	}

	// populate nodes
	for _, s := range shoots.Items {
		shootname := s.GetObjectMeta().GetName()
		node, ok := nodes[shootname]
		v := 0
		if ok {
			node.Project = s.GetObjectMeta().GetNamespace()
			v = 1
		} else {
			nodes[shootname] = &bg.Node{Id: shootname, Project: s.GetObjectMeta().GetNamespace(), Seed: false}
		}
		links = append(links, bg.Link{Source: shootname, Target: *s.Spec.Cloud.Seed, Value: v})
	}
	data, err := json.MarshalIndent(bg.Graph{Nodes: values(nodes), Links: &links}, "", "	")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
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
