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
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "envreplace",
	Short: "A tool that replaces environment variables in files",
	Long:  "A command-line tool that reads environment variables and replaces equally named variables in specified files",
	Run:   test,
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
	cobra.OnInitialize(initConfig)

	// Define flags
	rootCmd.Flags().StringSliceP("file", "f", []string{}, "file(s) to replace environment variables in")
	rootCmd.MarkFlagRequired("file")
	rootCmd.Flags().StringSliceP("prefix", "p", []string{}, "prefix(es) to filter environment variables by")

	// Bind flags to Viper
	viper.BindPFlag("file", rootCmd.Flags().Lookup("file"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}

func test(cmd *cobra.Command, args []string) {
	// Read environment variables
	envVars := os.Environ()
	// Iterate through environment variables and add them to the map
	for _, envVar := range envVars {
		pair := strings.Split(envVar, "=")
		log.Printf("%s=%s", pair[0], pair[1])
	}

	for _, filePath := range viper.GetStringSlice("file") {
		log.Printf("file: %s", filePath)
	}
}

func replace(cmd *cobra.Command, args []string) {
	// Read environment variables
	envVars := os.Environ()

	// Get file path from command-line flag or configuration file
	filePath := viper.GetString("file")
	if filePath == "" {
		log.Fatal("file path not specified")
	}

	// Open file
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Create a map to store the variable names and their corresponding values
	vars := make(map[string]string)

	// Iterate through environment variables and add them to the map
	for _, envVar := range envVars {
		pair := strings.Split(envVar, "=")
		vars[pair[0]] = pair[1]
	}

	// Iterate through each line of the file
	for scanner.Scan() {
		line := scanner.Text()

		// Iterate through variables in the map and check if they are present in the line
		for key, value := range vars {
			if strings.Contains(line, key) {
				// Replace the value of the variable in the line with the environment variable value
				newLine := strings.Replace(line, key+"=", key+"="+value, 1)
				// Move file pointer back to the beginning of the line
				_, err = file.Seek(-int64(len(line)), os.SEEK_CUR)
				if err != nil {
					log.Fatalf("failed to seek file: %v", err)
				}
				// Write the new line to the file
				_, err = fmt.Fprintln(file, newLine)
				if err != nil {
					log.Fatalf("failed to write to file: %v", err)
				}
			}
		}
	}
}
