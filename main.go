package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

func precheck(awsConfigPath string) {

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

func getAwsConfig(awsConfigPath string) []string {
	config, err := ini.Load(awsConfigPath)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error loading cinfig file: %s", awsConfigPath))
		os.Exit(-1)
	}
	sections := config.SectionStrings()
	return sections
}

func selectProfile(profiles []string) string {
	prompt := promptui.Select{
		Label: "Select Profile",
		Items: profiles,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(-1)
	}

	fmt.Printf("You choose %q\n", result)
	return result
}

func setDefaultProfile(defaultProfile string, awsConfigPath string) {
	config, err := ini.Load(awsConfigPath)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error loading cinfig file: %s", awsConfigPath))
		os.Exit(-1)
	}
	if config.HasSection("default") {
		if !config.Section("default").HasKey("awspDefault") {
			for _, element := range config.Section("default").KeyStrings() {
				config.Section("old_modified_default").Key(element).SetValue(config.Section("default").Key(element).String())
			}
		}
	}
	config.DeleteSection("default")
	config.NewSection("default")
	for _, element := range config.Section(defaultProfile).KeyStrings() {
		config.Section("default").Key(element).SetValue(config.Section(defaultProfile).Key(element).String())
	}
	config.Section("default").Key("awspDefault").SetValue("true")
	//config.Section("default").SetBody(config.Section(defaultProfile))
	config.SaveTo(awsConfigPath)
	//	todo - add awspDefault to an element in the default section as well as in the section selected
}
func main() {
	homeDir, _ := os.UserHomeDir()
	awsConfigPath := fmt.Sprintf("%s/.aws/config", homeDir)
	awsCredsPath := fmt.Sprintf("%s/.aws/credentials", homeDir)

	precheck(awsConfigPath)
	profiles := getAwsConfig(awsConfigPath)
	selectedProfile := selectProfile(profiles)
	setDefaultProfile(selectedProfile, awsConfigPath)
	setDefaultProfile(selectedProfile, awsCredsPath)
}
