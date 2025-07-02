/*
Copyright 2024 Pisomind Inc.

Portions of this file are derived from devpod-provider-terraform:
https://github.com/loft-sh/devpod-provider-terraform/pkg/options/options.go
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

package options

import (
	"fmt"
	"os"
)

const (
	CLOUDINIT_SSH_KEY        = "CLOUDINIT_SSH_KEY"
	CLOUDINIT_USERNAME       = "CLOUDINIT_USERNAME"
	CLOUDINIT_PASSWORD       = "CLOUDINIT_PASSWORD"
	CLOUDINIT_IP             = "CLOUDINIT_IP"
	CLOUDINIT_GATEWAY        = "CLOUDINIT_GATEWAY"
	NODE_NAME                = "NODE_NAME"
	PROXMOX_API_URL          = "PROXMOX_API_URL"
	PROXMOX_API_TOKEN_ID     = "PROXMOX_API_TOKEN_ID"
	PROXMOX_API_TOKEN_SECRET = "PROXMOX_API_TOKEN_SECRET"
	PROXMOX_VM_ID            = "PROXMOX_VM_ID"
	TERRAFORM_PROJECT        = "TERRAFORM_PROJECT"
)

type Options struct {
	MachineID     string
	MachineFolder string

	// Proxmox
	NodeName              string
	ProxmoxApiUrl         string
	ProxmoxApiTokenId     string
	ProxmoxApiTokenSecret string
	ProxmoxVmId           string

	// Cloudinit
	CloudinitSshKey   string
	CloudinitUsername string
	CloudinitPassword string
	CloudinitIp       string
	CloudinitGateway  string
}

func ConfigFromEnv() (Options, error) {
	return Options{
		NodeName:              os.Getenv(NODE_NAME),
		ProxmoxApiUrl:         os.Getenv(PROXMOX_API_URL),
		ProxmoxApiTokenId:     os.Getenv(PROXMOX_API_TOKEN_ID),
		ProxmoxApiTokenSecret: os.Getenv(PROXMOX_API_TOKEN_SECRET),
		ProxmoxVmId:           os.Getenv(PROXMOX_VM_ID),
		CloudinitSshKey:       os.Getenv(CLOUDINIT_SSH_KEY),
		CloudinitUsername:     os.Getenv(CLOUDINIT_USERNAME),
		CloudinitPassword:     os.Getenv(CLOUDINIT_PASSWORD),
		CloudinitIp:           os.Getenv(CLOUDINIT_IP),
		CloudinitGateway:      os.Getenv(CLOUDINIT_GATEWAY),
	}, nil
}

func FromEnv() (*Options, error) {
	retOptions := &Options{}

	var err error

	retOptions.MachineFolder, err = FromEnvOrError("MACHINE_FOLDER")
	if err != nil {
		return nil, err
	}

	retOptions.MachineID, err = FromEnvOrError("MACHINE_ID")
	if err != nil {
		return nil, err
	}
	// prefix with devpod-
	retOptions.MachineID = "devpod-" + retOptions.MachineID

	retOptions.NodeName, err = FromEnvOrError(NODE_NAME)
	if err != nil {
		return nil, err
	}

	retOptions.ProxmoxApiUrl, err = FromEnvOrError(PROXMOX_API_URL)
	if err != nil {
		return nil, err
	}

	retOptions.ProxmoxApiTokenId, err = FromEnvOrError(PROXMOX_API_TOKEN_ID)
	if err != nil {
		return nil, err
	}

	retOptions.ProxmoxApiTokenSecret, err = FromEnvOrError(PROXMOX_API_TOKEN_SECRET)
	if err != nil {
		return nil, err
	}

	retOptions.ProxmoxVmId, err = FromEnvOrError(PROXMOX_VM_ID)
	if err != nil {
		return nil, err
	}

	retOptions.CloudinitSshKey, err = FromEnvOrError(CLOUDINIT_SSH_KEY)
	if err != nil {
		return nil, err
	}

	retOptions.CloudinitUsername, err = FromEnvOrError(CLOUDINIT_USERNAME)
	if err != nil {
		return nil, err
	}

	retOptions.CloudinitPassword, err = FromEnvOrError(CLOUDINIT_PASSWORD)
	if err != nil {
		return nil, err
	}

	retOptions.CloudinitIp, err = FromEnvOrError(CLOUDINIT_IP)
	if err != nil {
		return nil, err
	}

	retOptions.CloudinitGateway, err = FromEnvOrError(CLOUDINIT_GATEWAY)
	if err != nil {
		return nil, err
	}

	return retOptions, nil
}

func FromEnvOrError(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf(
			"couldn't find option %s in environment, please make sure %s is defined",
			name,
			name,
		)
	}

	return val, nil
}
