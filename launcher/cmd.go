package launcher

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var CONFIG Config

func runCmdGetFlagValues(cmd *cobra.Command, display CobraDisplay) (Shortcut, Executable, []string, error) {
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

func addCmdGetFlagValues(cmd *cobra.Command, display CobraDisplay) (string, Shortcut, error) {
	s := Shortcut{}
	name, err := cmd.Flags().GetString("shortcut-name")
	if err != nil {
		return name, s, err
	}
	if name == "" {
		name = display.Prompt("Shortcut name")
	}
	s.Template, err = cmd.Flags().GetString("shortcut-name")
	if err != nil {
		return name, s, err
	}
	if s.Template == "" {
		s.Template = display.Prompt("Shortcut template")
	}

	return name, s, nil
}

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Go and launcht it!",
	Args:  cobra.NoArgs,
}

var runCmd = &cobra.Command{
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
		shortcut, executable, params, err := runCmdGetFlagValues(cmd, display)
		if err != nil {
			display.Error(err.Error())
			return
		}
		err = RunCommand(shortcut, executable, params)
		if err != nil {
			display.Error(err.Error())
			return
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

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a command to your config file",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		useGUI, _ := cmd.Flags().GetBool("use-gui")
		display := NewDisplay(useGUI, args)
		name, shortcut, err := addCmdGetFlagValues(cmd, display)
		if err != nil {
			display.Error(err.Error())
			return
		}
		err = CONFIG.AddShortcut(name, shortcut)
		if err != nil {
			display.Error(err.Error())
			return
		}
		fmt.Printf("Done! Added shortcut '%s' -> '%s'", name, shortcut.Template)
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
	runCmd.Flags().BoolP("use-gui", "g", false, "Uses GUI instead of CLI")
	runCmd.Flags().StringP("executable-name", "e", "", "The program that should execute your command template.")
	runCmd.Flags().StringP("shortcut-name", "s", "", "The name of the shortcut.")
	runCmd.Flags().StringArrayP("params", "p", []string{}, "(optional) The params for the command.")

	addCmd.Flags().BoolP("use-gui", "g", false, "Uses GUI instead of CLI")
	addCmd.Flags().StringP("shortcut-name", "s", "", "The name of the shortcut.")
	addCmd.Flags().StringP("shortcut-template", "t", "", "The template for the shortcut.")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(initCmd)

	var err error
	CONFIG, err = ParseConfig()
	if err != nil {
		panic(err)
	}
}

func Execute() {
	if err := runCmd.Execute(); err != nil {
		panic(err)
	}
}
