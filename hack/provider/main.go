/*
Copyright 2024 Pisomind Inc.

Portions of this file are derived from devpod-provider-terraform:
https://github.com/loft-sh/devpod-provider-terraform/hack/provider/main.go
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

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

var checksumMap = map[string]string{
	"./release/devpod-provider-proxmox-linux-amd64":       "##CHECKSUM_LINUX_AMD64##",
	"./release/devpod-provider-proxmox-linux-arm64":       "##CHECKSUM_LINUX_ARM64##",
	"./release/devpod-provider-proxmox-darwin-amd64":      "##CHECKSUM_DARWIN_AMD64##",
	"./release/devpod-provider-proxmox-darwin-arm64":      "##CHECKSUM_DARWIN_ARM64##",
	"./release/devpod-provider-proxmox-windows-amd64.exe": "##CHECKSUM_WINDOWS_AMD64##",
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Expected version as argument")
		os.Exit(1)
		return
	}

	content, err := os.ReadFile("./hack/provider/provider.yaml")
	if err != nil {
		panic(err)
	}

	replaced := strings.Replace(string(content), "##VERSION##", os.Args[1], -1)
	for k, v := range checksumMap {
		checksum, err := File(k)
		if err != nil {
			panic(fmt.Errorf("generate checksum for %s: %v", k, err))
		}

		replaced = strings.Replace(replaced, v, checksum, -1)
	}

	fmt.Print(replaced)
}

// File hashes a given file to a sha256 string
func File(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return strings.ToLower(hex.EncodeToString(hash.Sum(nil))), nil
}
