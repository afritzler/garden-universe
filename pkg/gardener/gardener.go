package gardener

import (
	"fmt"

	"github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	garden "github.com/gardener/gardener/pkg/client/garden/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Gardener interface {
	GetClientset() *garden.Clientset
	GetVersion() string
	GetShootList() (*v1beta1.ShootList, error)
	GetShoots() (*[]v1beta1.Shoot, error)
	GetSeedList() (*v1beta1.SeedList, error)
	GetSeeds() (*[]v1beta1.Seed, error)
}

type gardener struct {
	config    *restclient.Config
	clientset *garden.Clientset
	version   *version.Info
}

func NewGardener(kubeconfig string) (Gardener, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %s", err)
	}
	gardenClientset, err := garden.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create garden clientset: %s", err)
	}
	version, err := gardenClientset.ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get server version: %s", err)
	}
	return &gardener{config: config, clientset: gardenClientset, version: version}, nil
}

func (g *gardener) GetClientset() *garden.Clientset {
	return g.clientset
}

func (g *gardener) GetVersion() string {
	return g.version.String()
}

func (g *gardener) GetShootList() (*v1beta1.ShootList, error) {
	shootlist, err := g.clientset.GardenV1beta1().Shoots("").List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get shootlist: %s", err)
	}
	return shootlist, nil
}

func (g *gardener) GetShoots() (*[]v1beta1.Shoot, error) {
	shootlist, err := g.GetShootList()
	if err != nil {
		return nil, fmt.Errorf("failed to get shoots: %s", err)
	}
	return &shootlist.Items, nil
}

func (g *gardener) GetSeedList() (*v1beta1.SeedList, error) {
	seedlist, err := g.clientset.GardenV1beta1().Seeds().List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get seedlist: %s", err)
	}
	return seedlist, nil
}

func (g *gardener) GetSeeds() (*[]v1beta1.Seed, error) {
	seedlist, err := g.GetSeedList()
	if err != nil {
		return nil, fmt.Errorf("failed to get seeds: %s", err)
	}
	return &seedlist.Items, nil
}
