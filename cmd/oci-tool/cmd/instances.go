// Copyright © 2018 Steve Brunton <sbrunton@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package cmd

import (
	"fmt"

	"github.com/nshttpd/oci-tools/oci"
	"github.com/nshttpd/oci-tools/utils"
	"github.com/spf13/cobra"
	"os"
	"text/template"
)

const defaultListInstanceTmpl = `{{ .Instance.DisplayName }} {{ .Instance.LifecycleState }}
{{ range $key, $value := .Vnics }}- {{ .PublicIp }}	- {{ .PrivateIp }}{{ end }}
	{{ .Instance.Id }}
`

// instancesCmd represents the instances command
var instancesCmd = &cobra.Command{
	Use:     "instances",
	Aliases: []string{"instance"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var listInstancesCmd = &cobra.Command{
	Use:   "list",
	Short: "list instances in tenancy",
	Long: `List instances in the tenancy. This will list
all of them unless a compartment-id is specified. 

	oci-tool compute instances list
`,
	Run: func(cmd *cobra.Command, args []string) {
		listInstances(cmd)
	},
}

func init() {
	computeCmd.AddCommand(instancesCmd)
	instancesCmd.AddCommand(listInstancesCmd)
}

func listInstances(cmd *cobra.Command) {
	cid := cmd.Flag("compartment-id").Value.String()

	comparts, err := config.GetCompartments()
	if err != nil {
		utils.ErrorMsg("error fetching compartments from API", err)
		return
	}

	// get the instances in each Compartment
	var computes []*oci.Compute

	// can maybe goroutines this to make it faster?
	for _, c := range comparts.Compartments() {
		if cid != "" && *c.Id != cid {
			continue
		}
		cs, err := config.GetComputeInstances(c)
		if err != nil {
			utils.ErrorMsg(fmt.Sprintf("error fetching Computes for cid : %s", *c.Id), err)
		} else {
			computes = append(computes, cs...)
		}
	}

	tmpl, err := template.New("imageList").Parse(defaultListInstanceTmpl)
	if err == nil {
		for _, i := range computes {
			err := tmpl.Execute(os.Stdout, i)
			if err != nil {
				utils.ErrorMsg("error processing item for template", err)
			}
		}
	} else {
		utils.ErrorMsg("error setting up output template", err)
	}

}
