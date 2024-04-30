package server

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/stcraft/dragonfly/server/cmd"
	"github.com/stcraft/dragonfly/server/world"
)

// Console represents the Console Source that is responsible for execution of commands
var Console ConsoleSource

// startConsole initialises the Console after which Commands can be sent from the console
func startConsole(logger Logger) {
	Console.logger = logger
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		if t := strings.TrimSpace(scanner.Text()); len(t) > 0 {
			name := strings.Split(t, " ")[0]
			if c, ok := cmd.ByAlias(name); ok {
				c.Execute(strings.TrimPrefix(strings.TrimPrefix(t, name), " "), Console)
			} else {
				output := &cmd.Output{}
				output.Errorf("Could not find command '%s'", name)
				Console.SendCommandOutput(output)
			}
		}
	}
}

// ConsoleSource represents a Console Source that is used to send commands
type ConsoleSource struct {
	logger Logger
}

// Name returns the name of console.
func (ConsoleSource) Name() string { return "CONSOLE" }

// Position ...
func (ConsoleSource) Position() mgl64.Vec3 { return mgl64.Vec3{} }

// SendCommandOutput prints out command outputs.
func (s ConsoleSource) SendCommandOutput(o *cmd.Output) {
	for _, e := range o.Errors() {
		s.logger.Errorf(text.ANSI(e))
	}
	for _, m := range o.Messages() {
		s.logger.Infof(text.ANSI(m))
	}
}

// SendMessage prints out message in console
func (s ConsoleSource) SendMessage(message string, args ...any) {
	message = format(fmt.Sprintf(message, args...))
	s.logger.Infof(text.ANSI(message + "§r"))
}

// World ...
func (ConsoleSource) World() *world.World { return nil }

// format is a utility function to format a list of values to have spaces between them, but no newline at the
// end, which is typically used for sending messages, popups and tips.
func format(text string) string {
	return strings.TrimSuffix(strings.TrimSuffix(text, "\n"), "\n")
}
