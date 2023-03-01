/*
Copyright © 2023 Finn Kristian Meland

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information for this application",
	Run: func(cmd *cobra.Command, args []string) {
		appName := viper.Get("appName")
		version := viper.Get("appVersion")
		branch := viper.Get("appBranch")
		commit := viper.Get("appCommit")
		appBuild := viper.Get("appBuild")
		goVersion := viper.Get("goVersion")

		log.Printf("%s\n\tversion: %s\n\tbranch: %s\n\tcommit: %s\n\tbuild: %s\n\tgo-version: %s", appName, version, branch, commit, appBuild, goVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
