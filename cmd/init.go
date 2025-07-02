/*
Copyright 2024 Pisomind Inc.

Portions of this file are derived from devpod-provider-terraform:
https://github.com/loft-sh/devpod-provider-terraform/cmd/init.go
Copyright 2023 Loft Labs, Inc.
Licensed under the Apache License, Version 2.0

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"

	"github.com/pisomind/devpod-provider-proxmox/pkg/options"
	"github.com/pisomind/devpod-provider-proxmox/pkg/terraform"

	"github.com/loft-sh/devpod/pkg/config"
	"github.com/loft-sh/devpod/pkg/log"
	"github.com/loft-sh/devpod/pkg/provider"
	"github.com/spf13/cobra"
)

// InitCmd holds the cmd flags
type InitCmd struct{}

// NewInitCmd defines a init
func NewInitCmd() *cobra.Command {
	cmd := &InitCmd{}
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Init account",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.Run(
				context.Background(),
				provider.FromEnvironment(),
				log.Default,
			)
		},
	}

	return initCmd
}

// Run runs the init logic
func (cmd *InitCmd) Run(
	ctx context.Context,
	machine *provider.Machine,
	logs log.Logger,
) error {
	devpodPath, err := config.GetConfigDir()
	if err != nil {
		return err
	}

	terraformPath := devpodPath + "/bin/terraform"

	project, err := options.FromEnvOrError(options.TERRAFORM_PROJECT)
	if err != nil {
		return err
	}

	// create provider
	provider := &terraform.TerraformProvider{
		Log:     logs,
		Bin:     terraformPath,
		Project: project,
	}

	err = terraform.Install(provider)
	if err != nil {
		return err
	}

	return nil
}
