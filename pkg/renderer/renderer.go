// Copyright © 2018 Andreas Fritzler
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

package renderer

import (
	"github.com/afritzler/garden-universe/pkg/gardener"
	promclient "github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type Renderer interface {
	GetGraph() ([]byte, error)
}

type promrenderer struct {
	client promclient.Client
	api    v1.API
}

type kuberenderer struct {
	garden gardener.Gardener
}
