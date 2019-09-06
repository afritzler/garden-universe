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

package utils

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// GetKubeConfigFlagOrEnv returns either the kubeconfig from a flag or uses the KUBECONFIG env var.
func GetKubeConfigFlagOrEnv(cmd *cobra.Command) string {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = cmd.Flag("kubeconfig").Value.String()
		fmt.Printf("using kubeconfig: %s\n", kubeconfig)
	} else {
		fmt.Printf("using KUBECONFIG env: %s\n", kubeconfig)
	}
	return kubeconfig
}
