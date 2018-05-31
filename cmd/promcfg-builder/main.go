package main

import (
	"flag"
	"os"
	"github.com/mitchellh/go-homedir"
	"github.com/nshttpd/oci-tools/utils"
	"github.com/nshttpd/oci-tools/oci"
	"sync"
	"fmt"
	"text/template"
	"github.com/oracle/oci-go-sdk/identity"
	"strings"
)

const (
	defaultProfile = "DEFAULT"
	defaultConfig  = "/.oci/config"
)

func getRegionComputes(config oci.ClientConfig) []*oci.Compute {
	comparts, err := config.GetCompartments()
	if err != nil {
		utils.ErrorMsg("error fetching compartments from API", err)
	}

	computes := make([]*oci.Compute, 0)
	mux := &sync.Mutex{}

	var wg sync.WaitGroup
	for _, cpts := range comparts.Compartments() {
		wg.Add(1)
		go func(compart identity.Compartment) {
			defer wg.Done()
			fmt.Printf("processing compartment : %s\n", *compart.Name)
			cs, err := config.GetComputeInstances(compart)
			if err != nil {
				utils.ErrorMsg(fmt.Sprintf("error fetching Computes for cid : %s", *compart.Name), err)
			} else {
				mux.Lock()
				fmt.Printf("adding compartment %s to slice with %d values\n", *compart.Name, len(cs))
				computes = append(computes, cs...)
				fmt.Printf("lengh of computes slice is now %d\n", len(computes))
				mux.Unlock()
			}
			fmt.Printf("done processing compartment : %s\n", *compart.Name)
		}(cpts)
	}
	wg.Wait()

	return computes

}

func writeConfig(computes []*oci.Compute, tmpl *template.Template, promCfg *string) {

	f, err := os.Create(*promCfg)
	if err != nil {
		utils.ErrorMsg("error opening output file for prometheus config", err)
		return
	}
	defer f.Close()

	err = tmpl.Execute(f, &computes)
	if err != nil {
		utils.ErrorMsg("error processing template", err)
	}

	return
}

func main() {
	region := flag.String("region", "", "region to query")
	profile := flag.String("profile", defaultProfile, "profile to use")
	promTmpl := flag.String("template", "", "template file")
	promCfg := flag.String("output", "prometheus.yml", "output file name")
	flag.Parse()

	if *region == "" {
		utils.ErrorMsg("missing region value. region must be specified", nil)
		os.Exit(1)
	}

	if *promTmpl == "" {
		utils.ErrorMsg("missing template value. template must be specified", nil)
		os.Exit(1)
	}

	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	tmpl := template.New(*promTmpl)
	tmpl = tmpl.Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles(*promTmpl)

	if err != nil {
		utils.ErrorMsg("error reading and parsing promCfg template", err)
		os.Exit(1)
	}

	home, err := homedir.Dir()
	if err != nil {
		utils.ErrorMsg("error finding users home directory", err)
		os.Exit(1)
	}

	config, err := oci.CreateConfig(home+defaultConfig, *profile, *region)
	if err != nil {
		utils.ErrorMsg("error getting OCI config", err)
		os.Exit(1)
	}
	config.Region = *region
	computes := getRegionComputes(config)

	writeConfig(computes, tmpl, promCfg)


	return;
}