package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"gopkg.in/ini.v1"
	"log"
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
	for i := len(profiles) - 1; i >= 0; i-- {
		if profiles[i] == "old_modified_default" {
			profiles = append(profiles[:i], profiles[i+1:]...)
		}

	}
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

func setDefaultProfile(defaultProfile string, awsConfigPath string, unset bool) {
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

	config.SaveTo(awsConfigPath)
}

func unsetDefault(awsConfigPath string) {
	config, err := ini.Load(awsConfigPath)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error loading cinfig file: %s", awsConfigPath))
		os.Exit(-1)
	}
	if config.HasSection("default") {
		if config.Section("default").HasKey("awspDefault") {
			config.DeleteSection("default")
			if config.HasSection("old_modified_default") {
				config.NewSection("default")
				for _, element := range config.Section("old_modified_default").KeyStrings() {
					config.Section("default").Key(element).SetValue(config.Section("old_modified_default").Key(element).String())
				}
				config.DeleteSection("old_modified_default")
			}
		} else if config.HasSection("old_modified_default") {
			prompt := promptui.Select{
				Label: "The current default profile was not created using awsp. Would you like to delete the current default profile and revert to a previous default?",
				Items: []string{"Yes", "No"},
			}
			_, result, err := prompt.Run()
			if err != nil {
				log.Fatalf("Prompt failed %v\n", err)
			}
			if result == "Yes" {
				config.DeleteSection("default")
				config.NewSection("default")
				for _, element := range config.Section("old_modified_default").KeyStrings() {
					config.Section("default").Key(element).SetValue(config.Section("old_modified_default").Key(element).String())
				}
				config.DeleteSection("old_modified_default")
			}
		}
	} else if config.HasSection("old_modified_default") {
		config.NewSection("default")
		for _, element := range config.Section("old_modified_default").KeyStrings() {
			config.Section("default").Key(element).SetValue(config.Section("old_modified_default").Key(element).String())
		}
	}
	config.SaveTo(awsConfigPath)
}
func main() {
	homeDir, _ := os.UserHomeDir()
	awsConfigPath := fmt.Sprintf("%s/.aws/config", homeDir)
	awsCredsPath := fmt.Sprintf("%s/.aws/credentials", homeDir)

	//unsetCmd := flag.NewFlagSet("unset", flag.ExitOnError)
	//unsetEnable := unsetCmd.Bool("unset", false, "unset")
	if len(os.Args) > 1 {
		if os.Args[1] == "unset" {
			fmt.Println(fmt.Sprintf("Unsetting defaults in [ %s ]", awsConfigPath))
			unsetDefault(awsConfigPath)
			fmt.Println(fmt.Sprintf("Unsetting defaults in [ %s ]", awsCredsPath))
			unsetDefault(awsCredsPath)
		} else {
			fmt.Println("Expected unset subcommand.")
			os.Exit(1)
		}
	} else {
		precheck(awsConfigPath)
		profiles := getAwsConfig(awsConfigPath)
		selectedProfile := selectProfile(profiles)
		setDefaultProfile(selectedProfile, awsConfigPath, false)
		setDefaultProfile(selectedProfile, awsCredsPath, false)
	}
}
