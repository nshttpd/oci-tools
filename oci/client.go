package oci

import (
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
)

func CreateConfig(file string, profile string) common.ConfigurationProvider {
	config, err := common.ConfigurationProviderFromFileWithProfile(file, profile, "")
	if err != nil {
		fmt.Printf("error creating config for profile %s from file %s\n", profile, file)
		fmt.Printf("error : %v\n", err)
		config = nil
	}
	return config
}