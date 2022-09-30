package main

import (
	"fmt"
	"os"
	"strings"
)

func precheck() {
	homeDir, _ := os.UserHomeDir()
	awsConfigPath := fmt.Sprintf("%s/.aws/config", homeDir)
	_, err := os.Stat(awsConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("'" + awsConfigPath + "' file does not exist.")
			// exit the application
			os.Exit(-1)
		}
	}
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		if variable[0] == "DEFAULT_AWS_PROFILE" || variable[0] == "AWS_PROFILE" {
			print(fmt.Sprintf("The [ %s ] variable is set. Please unset the variable and rerun.\n", variable[0]))
			os.Exit(-1)
		}
	}
}

func getAwsConfig() {

}
func main() {
	precheck()

}
