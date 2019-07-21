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
	"html/template"

	"os"

	"oci-tools/utils"

	"github.com/spf13/cobra"
)

const defaultListCompartmentTmpl = "{{ .Name }} - {{ .Id }}\n"

// compartmentsCmd represents the compartments command
var compartmentsCmd = &cobra.Command{
	Use:     "compartments",
	Aliases: []string{"compartment"},
	Short:   "access compartment information",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var listCompartmentsCmd = &cobra.Command{
	Use:   "list",
	Short: "list compartments",
	Long: `List compartments in the tenancy.

	oci-tool iam compartments list

`,
	Run: func(cmd *cobra.Command, args []string) {
		listCompartments(cmd)
	},
}

func init() {
	iamCmd.AddCommand(compartmentsCmd)
	compartmentsCmd.AddCommand(listCompartmentsCmd)
}

func listCompartments(cmd *cobra.Command) {
	comparts, err := config.GetCompartments()
	if err != nil {
		utils.ErrorMsg("error fetching compartments from API", err)
		return
	}

	tmpl, err := template.New("compartmentList").Parse(defaultListCompartmentTmpl)
	if err == nil {
		for _, c := range comparts.Compartments() {
			err := tmpl.Execute(os.Stdout, c)
			if err != nil {
				utils.ErrorMsg("error processing compartment item for template", err)
			}
		}
	} else {
		utils.ErrorMsg("error creating output template", err)
	}

}
