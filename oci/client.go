package oci

import (
	"github.com/oracle/oci-go-sdk/common"
)

type ClientConfig struct {
	config common.ConfigurationProvider
}

func CreateConfig(file string, profile string) (ClientConfig, error) {
	c, err := common.ConfigurationProviderFromFileWithProfile(file, profile, "")
	if err != nil {
		return ClientConfig{}, err
	}
	return ClientConfig{config: c}, nil
}

func (c *ClientConfig) Config() common.ConfigurationProvider {
	return c.config
}

func (c *ClientConfig) TenancyOCID() (string, error) {
	return c.config.TenancyOCID()
}
