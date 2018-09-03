// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// computeCmd represents the compute command
var computeCmd = &cobra.Command{
	Use:   "compute",
	Short: "Actions based off of the compute/core API",
	Long: `Commands that will get data from the compute/core API to display
back out said information. For example:

	oci-tool compute images list

Will print out a user friendly set of output for images.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// self chaining
		rootCmd.PersistentPreRun(cmd, args)
		if cmd.Flag("region").Value.String() == "" {
			fmt.Printf("region must be supplied as a parameter\n")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(computeCmd)
}
