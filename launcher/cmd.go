package launcher

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var CONFIG Config

func getFlagValues(cmd *cobra.Command, display CobraDisplay) (Shortcut, Executable, []string, error) {
	s := Shortcut{}
	e := Executable{}
	params := []string{}
	executableName, err := cmd.Flags().GetString("executable-name")
	if err != nil {
		return s, e, params, err
	}
	if executableName == "" {
		return s, e, params, fmt.Errorf("the executable-name must be provided")
	}
	shortcutName, err := cmd.Flags().GetString("shortcut-name")
	if err != nil {
		return s, e, params, err
	}
	if shortcutName == "" {
		shortcutName = display.Prompt(fmt.Sprintf("[%s] Shortcut name", executableName))
	}
	params, err = cmd.Flags().GetStringArray("params")
	if err != nil {
		return s, e, params, err
	}
	s, err = CONFIG.GetShortcut(shortcutName)
	if err != nil {
		return s, e, params, err
	}
	e, err = CONFIG.GetExecutable(executableName)
	if err != nil {
		return s, e, params, err
	}
	if s.HasParams() && len(params) == 0 {
		promptResponse := display.Prompt("Params for template, comma separated")
		params = strings.Split(promptResponse, ",")
	}
	return s, e, params, nil
}

var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a shortcut",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		CONFIG, err = ParseConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		useGUI, _ := cmd.Flags().GetBool("use-gui")
		display := NewDisplay(useGUI, args)
		shortcut, executable, params, err := getFlagValues(cmd, display)
		if err != nil {
			display.Error(err.Error())
			panic(err)
		}
		err = RunCommand(shortcut, executable, params)
		if err != nil {
			display.Error(err.Error())
			panic(err)
		}
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a default configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		err := CreateDefaultConfig()
		if err != nil {
			fmt.Printf("Error creating config file: %v\n", err)
			os.Exit(1)
		}
		configPath, err := getConfigFilePath()
		if err != nil {
			fmt.Println("Default configuration file created successfully.")
		} else {
			fmt.Printf("Default configuration file created successfully at %s\n", configPath)
		}

	},
}

func init() {
	rootCmd.Flags().BoolP("use-gui", "g", false, "Uses GUI instead of CLI")
	rootCmd.Flags().StringP("executable-name", "e", "", "The program that should execute your command template.")
	rootCmd.Flags().StringP("shortcut-name", "s", "", "The name of the shortcut.")
	rootCmd.Flags().StringArrayP("params", "p", []string{}, "(optional) The params for the command.")
	rootCmd.AddCommand(initCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
