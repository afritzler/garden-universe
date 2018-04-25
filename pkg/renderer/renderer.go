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

	bg "github.com/afritzler/garden-universe/pkg/types"
	"github.com/gardener/gardener/pkg/apis/garden"
	"github.com/gardener/gardener/pkg/apis/garden/v1beta1"
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

	seednames := map[string]string{}

	// populate seeds
	for _, s := range seeds.Items {
		//nodeName := fmt.Sprintf("%s/%s", "seed", s.GetObjectMeta().GetName())
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
	for _, s := range shoots.Items {
		namespace := s.GetObjectMeta().GetNamespace()
		shootname := fmt.Sprintf("%s/%s", namespace, s.GetObjectMeta().GetName())
		node, ok := nodes[shootname]
		status := ""
		size := getSizeOfShoot(s)
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

func getSizeOfShoot(s v1beta1.Shoot) int {
	size := 0
	if s.Spec.Cloud.OpenStack != nil {
		for _, w := range s.Spec.Cloud.OpenStack.Workers {
			size += w.AutoScalerMin
		}
	}
	if s.Spec.Cloud.AWS != nil {
		for _, w := range s.Spec.Cloud.AWS.Workers {
			size += w.AutoScalerMin
		}
	}
	if s.Spec.Cloud.Azure != nil {
		for _, w := range s.Spec.Cloud.Azure.Workers {
			size += w.AutoScalerMin
		}
	}
	if s.Spec.Cloud.GCP != nil {
		for _, w := range s.Spec.Cloud.GCP.Workers {
			size += w.AutoScalerMin
		}
	}
	return size
}
