package oci

import (
	"bufio"
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/oracle/oci-go-sdk/common"
)

type ClientConfig struct {
	config common.ConfigurationProvider
	Region string
}

func CreateConfig(file string, profile string, region string) (ClientConfig, error) {
	c, err := common.ConfigurationProviderFromFileWithProfile(file, profile, "")
	if err != nil {
		return ClientConfig{}, err
	}

	t, _ := c.TenancyOCID()
	u, _ := c.UserOCID()
	f, _ := c.KeyFingerprint()
	k, _ := c.PrivateRSAKey()
	pk := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(k),
	}

	var b bytes.Buffer
	o := bufio.NewWriter(&b)
	err = pem.Encode(o, pk)
	if err != nil {
		return ClientConfig{}, fmt.Errorf("error encoding private key back to string : %v", err)
	}
	o.Flush()

	if region == "" {
		region, _ = c.Region()
	}

	r := common.NewRawConfigurationProvider(t, u, region, f, string(b.Bytes()[:len(b.Bytes())]), nil)

	return ClientConfig{config: r, Region: region}, nil
}

func (c *ClientConfig) Config() common.ConfigurationProvider {
	return c.config
}

func (c *ClientConfig) TenancyOCID() (string, error) {
	return c.config.TenancyOCID()
}
