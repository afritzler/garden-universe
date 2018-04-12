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
		//		fmt.Printf("%s\n", s.GetObjectMeta().GetName())
		nodes[s.GetObjectMeta().GetName()] = &bg.Node{Id: s.GetObjectMeta().GetName(), Project: "", Seed: true}
	}

	// populate nodes
	//var shootnodes = make([]bg.Node, len(shoots.Items))
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
