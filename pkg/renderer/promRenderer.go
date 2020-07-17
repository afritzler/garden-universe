package renderer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	bg "github.com/afritzler/garden-universe/pkg/types"
	promclient "github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

const (
	shootQuery        = "garden_shoot_info"
	nodeQuery         = "sum(shoot:kube_node_info:count) by (name, project)"
	availabilityQuery = "garden_shoot_condition{condition='APIServerAvailable'}"
)

func NewPromRenderer(client promclient.Client) Renderer {
	api := v1.NewAPI(client)
	return &promrenderer{client: client, api: api}
}

func (r *promrenderer) GetGraph() ([]byte, error) {
	var (
		shootVector        model.Vector
		nodeVector         model.Vector
		availabilityVector model.Vector

		prometheusQueryErrors []error
		errLock               sync.RWMutex
		wg                    sync.WaitGroup

		// Define function to exec a Prometheus query and to store the result in the passed Vector.
		execPrometheusQuery = func(query string, vector *model.Vector) {
			defer wg.Done()
			result, _, err := r.api.Query(context.Background(), query, time.Now())
			switch result.Type() {
			case model.ValVector:
				*vector = result.(model.Vector)
			default:
				fmt.Println("Could handle prometheus query")
			}
			if err != nil {
				errLock.Lock()
				defer errLock.Unlock()
				prometheusQueryErrors = append(prometheusQueryErrors, err)
			}
		}
	)

	wg.Add(3)
	go execPrometheusQuery(shootQuery, &shootVector)
	go execPrometheusQuery(nodeQuery, &nodeVector)
	go execPrometheusQuery(availabilityQuery, &availabilityVector)
	wg.Wait()

	// Check if any prometheus queries failed
	for _, e := range prometheusQueryErrors {
		if e != nil {
			return nil, e
		}
	}

	nodes := make(map[string]*bg.Node)
	links := make([]bg.Link, 0)

	// populate seeds and shoots
	for _, shoot := range shootVector {
		seedname := "garden/" + string(shoot.Metric["seed"])
		shootname := fmt.Sprintf("%s/%s", string(shoot.Metric["project"]), string(shoot.Metric["name"]))
		project := string(shoot.Metric["project"])
		region := string(shoot.Metric["region"])

		// If this seed hasn't been found yet add it to existing
		if _, exists := nodes[seedname]; !exists {
			nodes[seedname] = &bg.Node{Id: seedname, Project: project, Name: seedname, Type: "seed"}
		}

		// TODO: improve soil cluster discovery
		if strings.Contains(seedname, "soil") {
			nodes[shootname] = &bg.Node{Id: shootname, Project: project, Name: shootname, Type: "seed", Region: region}
		} else {
			nodes[shootname] = &bg.Node{Id: shootname, Project: project, Name: shootname, Type: "shoot", Region: region}
		}
		links = append(links, bg.Link{Source: shootname, Target: nodes[seedname].Id})
	}

	// Populate shoots with node info
	for _, shoot := range nodeVector {
		shootname := fmt.Sprintf("%s/%s", string(shoot.Metric["project"]), string(shoot.Metric["name"]))
		if node, exists := nodes[shootname]; exists {
			node.Size = int32(shoot.Value)
		} else {
			fmt.Printf("found node info for non existent shoot %s\n", shootname)
		}
	}

	// Populate availability
	for _, shoot := range availabilityVector {
		shootname := fmt.Sprintf("%s/%s", string(shoot.Metric["project"]), string(shoot.Metric["name"]))
		if node, exists := nodes[shootname]; exists {
			node.Availability = shoot.Value != 0.0
		} else {
			fmt.Printf("found availability info for non existent shoot %s\n", shootname)
		}
	}

	data, err := json.MarshalIndent(bg.Graph{Nodes: r.values(nodes), Links: &links}, "", "	")
	if err != nil {
		return nil, fmt.Errorf("JSON marshaling failed: %s", err)
	}
	return data, nil
}

func (r *promrenderer) values(nodes map[string]*bg.Node) *[]bg.Node {
	array := []bg.Node{}
	for _, n := range nodes {
		array = append(array, *n)
	}
	return &array
}
