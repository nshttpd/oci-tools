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

	"github.com/spf13/cobra"
	"github.com/oracle/oci-go-sdk/core"
	"os"
	"context"
)

// imagesCmd represents the images command
var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "get a list of images for a region",
	Long: `get a list of images for a region in a user friendly output
format.`,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imagesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listImages(client core.ComputeClient) {
	cid, err := config.TenancyOCID()
	if err != nil {
		fmt.Printf("error fetching tenancy ocid\n")
		fmt.Printf("error: %v\n", err)
	}
	req := core.ListImagesRequest{CompartmentId: &cid}
	res, err := client.ListImages(context.Background(), req)
	if err != nil {
		fmt.Printf("error fetching image list\n")
		fmt.Printf("error: %v\n", err)
	} else {
		for _, i := range res.Items {
			fmt.Printf("name : %s\n", *i.DisplayName)
		}
	}
}