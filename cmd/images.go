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

	"context"
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
		c, err := core.NewComputeClientWithConfigurationProvider(config)
		if err != nil {
			fmt.Printf("error creating Compute Client\n")
			fmt.Printf("error : %v\n", err)
			os.Exit(1)
		}
		c.SetRegion(cmd.Flag("region").Value.String())
		listImages(c)
	},
}

func init() {
	computeCmd.AddCommand(imagesCmd)
	imagesCmd.AddCommand(listImagesCmd)

	imagesCmd.PersistentFlags().StringVar(&operatingSystem, "operating-system", "", "limit images to operating system")
}

func listImages(client core.ComputeClient) {
	cid, err := config.TenancyOCID()
	if err != nil {
		fmt.Printf("error fetching tenancy ocid\n")
		fmt.Printf("error: %v\n", err)
	} else {
		req := core.ListImagesRequest{CompartmentId: &cid}
		if operatingSystem != "" {
			req.OperatingSystem = &operatingSystem
		}
		res, err := client.ListImages(context.Background(), req)
		if err != nil {
			fmt.Printf("error fetching image list\n")
			fmt.Printf("error: %v\n", err)
		} else {
			tmpl, err := template.New("imageList").Parse(defaultListImageTmpl)
			if err == nil {
				for _, i := range res.Items {
					err := tmpl.Execute(os.Stdout, i)
					if err != nil {
						fmt.Println("error processing item for template")
						fmt.Printf("error : %v\n", err)
					}
				}
			} else {
				fmt.Println("error setting up output template")
				fmt.Printf("error : %v\n", err)
			}
		}
	}
}
