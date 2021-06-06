package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

const gitConf string = "gitconf.config"

func checkFileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		return false, nil
	}

	return err == nil, err
}

func atomicFileCopy(sourcePath, targetPath string) error {
	contents, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	tempPath := targetPath + ".tmp_gitconf"

	if err = os.WriteFile(tempPath, contents, 0666); err != nil {
		return err
	}

	return os.Rename(tempPath, targetPath)
}

func createConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := homeDir + "/.config/gitconf"

	err = os.MkdirAll(configPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return configPath, nil
}

func checkGitConfigFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	gitConfigPath := homeDir + "/.gitconfig"
	gitConfigExists, err := checkFileExists(gitConfigPath)

	if err != nil {
		return "", err
	}
	if !gitConfigExists {
		return "", errors.New(fmt.Sprintf(".gitconfig file does not exist at %s", gitConfigPath))
	}

	return gitConfigPath, nil
}

func currentProfile(configPath string) (string, error) {
	raw, err := os.ReadFile(configPath + "/" + gitConf)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(raw)), nil
}

func updateProfile(configPath, profile string) error {
	return os.WriteFile(configPath+"/"+gitConf, []byte(profile+"\n"), 0666)
}

func main() {
	showCmd := flag.NewFlagSet("show", flag.ExitOnError)
	setCmd := flag.NewFlagSet("set", flag.ExitOnError)

	missingSubMsg := "Expected 'show' or 'set' subcommands.\n"

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, missingSubMsg)
		os.Exit(1)
	}

	configPath, err := createConfigPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to locate config: %s.\n", err)
		os.Exit(1)
	}

	gitConfigPath, err := checkGitConfigFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to locate config: %s.\n", err)
		os.Exit(1)
	}

	profile, err := currentProfile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to extract current profile: %s.\n", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "show":
		showCmd.Parse(os.Args[2:])
		fmt.Printf("Current git config: %s\n", profile)
	case "set":
		setCmd.Parse(os.Args[2:])

		if nSetArgs := len(setCmd.Args()); nSetArgs != 1 {
			fmt.Fprintf(os.Stderr, "Expected exactly one profile set arg, got %d.\n", nSetArgs)
			os.Exit(1)
		}

		targetProfile := setCmd.Args()[0]
		if targetProfile == profile {
			fmt.Fprintf(os.Stdout, "Git config is already %s. No-op.\n", profile)
			os.Exit(0)
		}

		fmt.Println("set profile:", targetProfile)

		if err = atomicFileCopy(configPath+"/"+targetProfile+".gitconfig", gitConfigPath); err != nil { // TODO change target from configPath to $HOME
			fmt.Fprintf(os.Stderr, "Error setting new git config: %s\n", err)
			os.Exit(1)
		}

		profile = targetProfile

		if err = updateProfile(configPath, profile); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to update the current gitconf config file: %s\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, missingSubMsg)
	}

	os.Exit(0)
}
