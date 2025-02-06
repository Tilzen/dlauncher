package launcher

import (
	"fmt"
	"os"

	"github.com/alecthomas/repr"
	"github.com/ncruces/zenity"
)

type CobraDisplay interface {
	Prompt(msg string) string
	Error(msg string)
	Info(msg string)
	Debug(params any)
}

type DisplayCLI struct {
	args []string
}
type DisplayGUI struct {
	args []string
}

func NewDisplay(useGUI bool, args []string) CobraDisplay {
	if useGUI {
		return DisplayGUI{args}
	}
	return DisplayCLI{args}
}

// DisplayCLI
func (cli DisplayCLI) Prompt(msg string) string {
	return cli.args[0]
}
func (cli DisplayCLI) Error(msg string) {
	cli.Info(msg)
}
func (cli DisplayCLI) Info(msg string) {
	fmt.Println(msg)
}
func (cli DisplayCLI) Debug(params any) {
	repr.Println(params)
}

// DisplayGUI
func (gui DisplayGUI) Prompt(msg string) string {
	resp, err := zenity.Entry(msg)
	if err != nil {
		zenity.Error(err.Error())
		os.Exit(1)
	}
	return resp
}
func (gui DisplayGUI) Error(msg string) {
	zenity.Error(msg)
}
func (gui DisplayGUI) Info(msg string) {
	zenity.Info(msg)
}
func (gui DisplayGUI) Debug(params any) {
	zenity.Info(repr.String(params))
}
