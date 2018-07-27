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
	"os"

	"context"

	"fmt"

	"github.com/nshttpd/oci-tools/utils"
	"github.com/oracle/oci-go-sdk/identity"
	"github.com/spf13/cobra"
)

// idpCmd represents the idp command
var idpCmd = &cobra.Command{
	Use:   "idp",
	Short: "access idp information",
	Long: `Deal with IDP level services for listing and deleting
them if something goes wrong through the UI while trying to set up
for Federation:

oci-tool iam idp list

oci-tool iam idp delete

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var idpListCmd = &cobra.Command{
	Use:   "list",
	Short: "list idp configs",
	Long: `List idp configs in the tenancy.

	oci-tool iam idp list

`,
	Run: func(cmd *cobra.Command, args []string) {
		listIdentityProviders(cmd)
	},
}

func init() {
	iamCmd.AddCommand(idpCmd)
	idpCmd.AddCommand(idpListCmd)
}

func listIdentityProviders(cmd *cobra.Command) {
	iam, err := identity.NewIdentityClientWithConfigurationProvider(config.Config())
	if err != nil {
		utils.ErrorMsg("error creating Identity Client", err)
		os.Exit(1)
	}
	tid, err := config.TenancyOCID()
	if err != nil {
		utils.ErrorMsg("error grabbing tenancy ocid from config", err)
		os.Exit(1)
	}
	req := identity.ListIdentityProvidersRequest{
		Protocol:      identity.ListIdentityProvidersProtocolSaml2,
		CompartmentId: &tid,
	}
	res, err := iam.ListIdentityProviders(context.Background(), req)
	fmt.Printf("items length: %d\n", len(res.Items))
	if err != nil {
		utils.ErrorMsg("error fetching list of identity providers from API", err)
		os.Exit(1)
	}
	for _, i := range res.Items {
		fmt.Println(i)
	}
}
