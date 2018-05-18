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
	"context"
	"github.com/nshttpd/oci-tools/utils"
	"github.com/oracle/oci-go-sdk/core"
	"github.com/spf13/cobra"
	"os"
	"text/template"
)

const defaultListImageTmpl = "{{ .DisplayName }}   {{ .Id }}\n"

var operatingSystem string

// imagesCmd represents the images command
var imagesCmd = &cobra.Command{
	Use:     "images",
	Aliases: []string{"image"},
	Short:   "image api actions",
	Long:    `Run actions against the images core api`,
}

var listImagesCmd = &cobra.Command{
	Use:   "list",
	Short: "list images in a tenancies region",
	Long: `List images available in a tenancies region
additional filter can be applied if necessary, example:

oci-tool compute images list --operating-system "Oracle Linux"`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := core.NewComputeClientWithConfigurationProvider(config.Config())
		if err != nil {
			utils.ErrorMsg("error creating Compute Client", err)
			os.Exit(1)
		}
		c.SetRegion(cmd.Flag("region").Value.String())
		listImages(cmd, c)
	},
}

func init() {
	computeCmd.AddCommand(imagesCmd)
	imagesCmd.AddCommand(listImagesCmd)

	imagesCmd.PersistentFlags().StringVar(&operatingSystem, "operating-system", "", "limit images to operating system")
}

func listImages(cmd *cobra.Command, client core.ComputeClient) {
	cid := cmd.Flag("compartment-id").Value.String()
	if cid == "" {
		var err error
		cid, err = config.TenancyOCID()
		if err != nil {
			utils.ErrorMsg("error fetching tenancy ocid", err)
			os.Exit(1)
		}
	}

	req := core.ListImagesRequest{CompartmentId: &cid}
	if operatingSystem != "" {
		req.OperatingSystem = &operatingSystem
	}
	res, err := client.ListImages(context.Background(), req)
	if err != nil {
		utils.ErrorMsg("error fetching image list", err)
	} else {
		tmpl, err := template.New("imageList").Parse(defaultListImageTmpl)
		if err == nil {
			for _, i := range res.Items {
				err := tmpl.Execute(os.Stdout, i)
				if err != nil {
					utils.ErrorMsg("error processing item for template", err)
				}
			}
		} else {
			utils.ErrorMsg("error setting up output template", err)
		}
	}
}
