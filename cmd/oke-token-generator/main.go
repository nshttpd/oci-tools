package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"oci-tools/oci"
	"oci-tools/utils"
	"os"
	"time"

	"github.com/mitchellh/go-homedir"

	"github.com/oracle/oci-go-sdk/common"
	flag "github.com/spf13/pflag"
)

type ExecStatus struct {
	Token               string `json:"token"`
	ExpirationTimestamp string `json:"expirationTimestamp"`
}

type ExecCredential struct {
	ApiVersion string     `json:"apiVersion"`
	Kind       string     `json:"kind"`
	Status     ExecStatus `json:"status"`
}

const (
	TimeFormat     = "Mon, 02 Jan 2006 15:04:05 GMT"
	defaultProfile = "DEFAULT"
	defaultConfig  = "/.oci/config"
)

func createSignature(config oci.ClientConfig, clusterId string) {
	client, err := common.NewClientWithConfig(config.Config())
	if err != nil {
		panic(err)
	}

	signer := client.Signer

	uri := fmt.Sprintf("https://containerengine.%s.oraclecloud.com/cluster_request/%s", config.Region, clusterId)
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		panic(err)
	}

	req.Header.Set("date", time.Now().UTC().Format(TimeFormat))

	err = signer.Sign(req)

	if err != nil {
		panic(err)
	}

	tokenRequest, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		panic(err)
	}

	query := tokenRequest.URL.Query()
	query.Set("authorization", req.Header.Get("authorization"))
	query.Set("date", req.Header.Get("Date"))

	// TODO check if IP and add OBO

	tokenRequest.URL.RawQuery = query.Encode()
	b := tokenRequest.URL.String()

	b64 := base64.StdEncoding.EncodeToString([]byte(b))

	token := ExecCredential{
		ApiVersion: "client.authentication.k8s.io/v1beta1",
		Kind:       "ExecCredential",
		Status: ExecStatus{
			Token:               b64,
			ExpirationTimestamp: time.Now().Add(time.Duration(4) * time.Minute).UTC().Format(time.RFC3339Nano),
		},
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(token)

	if err != nil {
		panic(err)
	}
}

func main() {

	var ociConfig oci.ClientConfig

	region := flag.String("region", "", "OCI region the cluster lives in")
	clusterId := flag.String("cluster-id", "", "OKE Cluster OCID")
	profile := flag.String("profile", defaultProfile, "OCI Config Profile")
	config := flag.String("config", defaultConfig, "config file (default is $HOME/.oci/config)")
	flag.Parse()

	home, err := homedir.Dir()
	if err != nil {
		fmt.Printf("error finding users home directory : %s\n", err)
		os.Exit(1)
	}

	ociConfig, err = oci.CreateConfig(fmt.Sprintf("%s%s", home, *config), *profile, *region)

	if err != nil {
		utils.ErrorMsg("error getting OCI config", err)
		os.Exit(1)
	}

	createSignature(ociConfig, *clusterId)

	os.Exit(0)
}
