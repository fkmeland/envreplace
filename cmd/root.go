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
	"io"
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
	Run:   process,
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
	rootCmd.Flags().BoolP("verbose", "v", false, "verbose output")

	// Bind flags to Viper
	viper.BindPFlag("file", rootCmd.Flags().Lookup("file"))
	viper.BindPFlag("prefix", rootCmd.Flags().Lookup("prefix"))
	viper.BindPFlag("verbose", rootCmd.Flags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Set some default values
	viper.SetDefault("appName", "envreplace")
	viper.SetDefault("appVersion", "v0.0.0")
	viper.SetDefault("appBranch", "n/a")
	viper.SetDefault("appCommit", "n/a")
	viper.SetDefault("appBuild", "n/a")
	viper.SetDefault("goVersion", "n/a")

	// Read in environment variables that match for the application
	viper.AutomaticEnv()
}

// verbose prints a message if the verbose flag is set
func verbose(msg string) {
	if viper.GetBool("verbose") {
		log.Println(msg)
	}
}

func process(cmd *cobra.Command, args []string) {
	// Read environment variables
	envVars := os.Environ()

	// Print verbose message if verbose flag is set
	verbose("processing environment variables")

	// If prefixes are provided, filter environment variables by prefix
	prefixes := viper.GetStringSlice("prefix")
	if len(prefixes) > 0 {
		// Print verbose message if verbose flag is set
		verbose(fmt.Sprintf("filtering environment variables by %d prefix(es): %s", len(prefixes), prefixes))
		var filteredEnvVars []string
		for _, envVar := range envVars {
			pair := strings.Split(envVar, "=")
			for _, prefix := range prefixes {
				if strings.HasPrefix(pair[0], prefix) {
					// Print verbose message if verbose flag is set
					verbose(fmt.Sprintf("added prefixed environment variable: %s", strings.Split(envVar, "=")[0]))
					filteredEnvVars = append(filteredEnvVars, envVar)
				}
			}
		}
		envVars = filteredEnvVars
	}

	// Print verbose message if verbose flag is set
	verbose(fmt.Sprintf("read total of %d environment variables", len(envVars)))

	// Get file path from command-line flag or configuration file
	files := viper.GetStringSlice("file")
	if len(files) == 0 {
		log.Fatal("files not provided")
	}

	// Print verbose message if verbose flag is set
	verbose(fmt.Sprintf("processing %d file(s): %s", len(files), files))

	// Iterate through each file path
	for _, filePath := range files {
		// Print verbose message if verbose flag is set
		verbose(fmt.Sprintf("processing file: %s", filePath))

		// Open file
		file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
		if err != nil {
			log.Fatalf("failed to open file: %v", err)
		}
		defer file.Close()

		// Create a temporary file to hold the modified contents
		tempFile, err := os.CreateTemp("", "tempfile")
		if err != nil {
			log.Fatalf("failed to create temporary file: %v", err)
		}
		defer tempFile.Close()

		// Create a scanner to read the file line by line
		scanner := bufio.NewScanner(file)

		// Create a map to store the variable names and their corresponding values
		vars := make(map[string]string)

		// Iterate through environment variables and add them to the map
		for _, envVar := range envVars {
			pair := strings.Split(envVar, "=")
			vars[pair[0]] = fmt.Sprintf("\"%s\"", strings.Join(pair[1:], "="))
		}

		// Create a set to keep track of the keys that have been written to the temporary file
		writtenLines := make(map[string]bool)

		// Iterate through each line of the file
		for scanner.Scan() {
			line := scanner.Text()

			// Iterate through variables in the map and check if they are present in the line
			for key, value := range vars {
				if strings.HasPrefix(line, key) {
					// Print verbose message if verbose flag is set
					verbose(fmt.Sprintf("updating %s to value %s in file %s", key, value, filePath))

					// Replace the value of the variable in the line with the environment variable value
					newLine := strings.Replace(line, line, key+"="+value, 1)

					// Print verbose message if verbose flag is set
					verbose(fmt.Sprintf("writing new line to file: %s", newLine))

					// Write the new line to the file
					_, err := tempFile.WriteString(newLine + "\n")
					if err != nil {
						log.Fatalf("failed to write to file: %v", err)
					}

					// Mark the line as written
					writtenLines[line] = true
				}
			}

			if !writtenLines[line] {
				// Write the original line to the temporary file
				_, err := tempFile.WriteString(line + "\n")
				if err != nil {
					log.Fatalf("failed to write to file: %v", err)
				}

				// Mark the line as written
				writtenLines[line] = true
			}
		}

		// Check for any errors that occurred while scanning the file
		if err := scanner.Err(); err != nil {
			log.Fatalf("failed to read file: %v", err)
		}

		// Truncate the original file
		err = file.Truncate(0)
		if err != nil {
			panic(err)
		}

		// Seek to the beginning of the file
		_, err = file.Seek(0, 0)
		if err != nil {
			panic(err)
		}

		// Copy the modified contents from the temporary file to the original file
		_, err = tempFile.Seek(0, 0)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(file, tempFile)
		if err != nil {
			panic(err)
		}

		// Remove the temporary file
		err = os.Remove(tempFile.Name())
		if err != nil {
			panic(err)
		}
	}
}
