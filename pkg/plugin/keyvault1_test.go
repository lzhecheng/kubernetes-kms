// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package plugin

import (
	"context"
	"fmt"
	"strings"
	"testing"

	kv "github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/kubernetes-kms/pkg/config"
)

func TestGetKey1(t *testing.T) {
	tests := []struct {
		desc                  string
		config                *config.AzureConfig
		vaultName             string
		keyName               string
		keyVersion            string
		keyVersionlessEnabled bool
		proxyMode             bool
		proxyAddress          string
		proxyPort             int
		managedHSM            bool
		expectedVaultURL      string
	}{
		{
			desc: "no error",
			config: &config.AzureConfig{
				TenantID:     "tenantid",
				ClientID:     "clientid",
				ClientSecret: "clientsecret"},
			vaultName:        "xxx",
			keyName:          "key1",
			keyVersion:       "keyversion",
			proxyMode:        false,
			expectedVaultURL: "https://xxxx.vault.azure.net/",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			kvClient, err := NewKeyVaultClient(test.config, test.vaultName, test.keyName, test.keyVersion, test.proxyMode, test.proxyAddress, test.proxyPort, test.managedHSM)
			if err != nil {
				t.Fatalf("newKeyVaultClient() failed with error: %v", err)
			}
			if kvClient == nil {
				t.Fatalf("newKeyVaultClient() expected kv client to not be nil")
			}
			if !strings.Contains(kvClient.GetUserAgent(), "k8s-kms-keyvault") {
				t.Fatalf("newKeyVaultClient() expected k8s-kms-keyvault user agent")
			}
			if kvClient.GetVaultURL() != test.expectedVaultURL {
				t.Fatalf("expected vault URL: %v, got vault URL: %v", test.expectedVaultURL, kvClient.GetVaultURL())
			}

			kvc := kvClient.(*KeyVaultClient)
			// version, err := kvc.GetLatestKeyVersion(context.TODO())

			// if err != nil {
			// 	t.Fatalf("GetLatestKeyVersion1() failed with error: %v", err)
			// }
			// fmt.Println("version is", version)

			somestr := "123"
			params := kv.KeyOperationsParameters{
				Algorithm: kv.RSAOAEP256,
				Value:     &somestr,
			}
			result, err := kvc.baseClient.Encrypt(context.TODO(), kvc.vaultURL, kvc.keyName, "", params)

			if err != nil {
				t.Fatalf("Encrypt() failed with error: %v", err)
			}

			if result.Kid != nil {
				fmt.Println("result is", *result.Kid)
				fmt.Println("result is", *result.Result)
			}
			fmt.Println("finish")
		})
	}
}
