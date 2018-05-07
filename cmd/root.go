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
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/nshttpd/oci-tool/oci"
	"github.com/oracle/oci-go-sdk/common"
)

const (
	defaultProfile = "DEFAULT"
	defaultConfig = "/.oci/config"
)

var (
	cfgFile string
	region string
	profile string
	compartment string
	config common.ConfigurationProvider
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "oci-tool",
	Short: "User Friendly Interactions with the Oracle OCI API",
	Long: `User friendly interactions with the Oracle OCI API to extract information 
that one may need or want on a daily basis. Unlike the OCI CLI application
this tool will bundle and concatenate information that is handy into one
execution for the user. For example: 

	oci-tool compute images

will gather all images across all compartments and produce output. The user can 
filter to only a single compartment if needed. It will also get the vnics and 
addresses for the hosts.
`,
	PersistentPreRun: func(cobra *cobra.Command, args []string) {
		if cobra.Flag("region").Value.String() == "" {
			fmt.Printf("region must be supplied as a parameter\n")
			os.Exit(1)
		}
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		v := cobra.Flag("profile").Value
		f := cobra.Flag("config").Value

		config = oci.CreateConfig(home + f.String(), v.String())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", defaultConfig, "config file (default is $HOME/.oci/config)")
	rootCmd.PersistentFlags().StringVar(&region, "region", "", "OCI region")
	rootCmd.PersistentFlags().StringVar(&profile, "profile", defaultProfile, "Config File Profile")
	rootCmd.PersistentFlags().StringVar(&compartment, "compartment", "", "Compartment Name")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

