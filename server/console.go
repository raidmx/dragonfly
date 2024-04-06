package server

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/STCraft/dragonfly/server/cmd"
	"github.com/STCraft/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
)

// console represents the Console Source that is responsible for execution of commands
var console = Console{log: logrus.New()}

// startConsole initialises the Console after which Commands can be sent from the console
func startConsole() {
	console.log.Formatter = &logrus.TextFormatter{ForceColors: true}
	console.log.Level = logrus.DebugLevel
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		if t := strings.TrimSpace(scanner.Text()); len(t) > 0 {
			name := strings.Split(t, " ")[0]
			if c, ok := cmd.ByAlias(name); ok {
				c.Execute(strings.TrimPrefix(strings.TrimPrefix(t, name), " "), console)
			} else {
				output := &cmd.Output{}
				output.Errorf("Could not find command '%s'", name)
				console.SendCommandOutput(output)
			}
		}
	}
}

// Console represents a Console Source that is used to send commands
type Console struct {
	log *logrus.Logger
}

// Name returns the name of console.
func (Console) Name() string { return "CONSOLE" }

// Position ...
func (Console) Position() mgl64.Vec3 { return mgl64.Vec3{} }

// SendCommandOutput prints out command outputs.
func (s Console) SendCommandOutput(o *cmd.Output) {
	for _, e := range o.Errors() {
		s.log.Error(text.ANSI(e))
	}
	for _, m := range o.Messages() {
		s.log.Info(text.ANSI(m))
	}
}

// SendMessage prints out message in console
func (s Console) SendMessage(message string) {
	message = format(fmt.Sprintln(message))
	s.log.Info(text.ANSI(message + "§r"))
}

// SendMessagef sends a formatted message using a specified format to console
func (s Console) SendMessagef(message string, args ...any) {
	s.SendMessage(fmt.Sprintf(message, args...))
}

// World ...
func (Console) World() *world.World { return nil }

// format is a utility function to format a list of values to have spaces between them, but no newline at the
// end, which is typically used for sending messages, popups and tips.
func format(text string) string {
	return strings.TrimSuffix(strings.TrimSuffix(text, "\n"), "\n")
}
