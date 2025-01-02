package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/AbdessamadEnabih/Vertex/internal/cli/commands"
	"github.com/AbdessamadEnabih/Vertex/internal/persistance"
	vertex_log "github.com/AbdessamadEnabih/Vertex/pkg/log"
	"github.com/AbdessamadEnabih/Vertex/pkg/state"
	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

// Define the root command and subcommands
var rootCmd = &cobra.Command{
	Use:   "vertex",
	Short: "Run vertex commands",
}

// GlobalState is a pointer to the global state of the application. It is used to store
// key-value pairs in memory.
var GlobalState *state.State

// refreshInterval is the interval at which the global state is refreshed from persistence.
const refreshInterval = 60 * time.Second

// init initializes the CLI by loading the global state and adding the commands to the root command.
func init() {
	GlobalState, _ = persistance.Load()
	rootCmd.AddCommand(
		commands.NewGetCmd(GlobalState),
		commands.NewSetCmd(GlobalState),
		commands.NewDeleteCmd(GlobalState),
		commands.NewFlushCmd(GlobalState),
		commands.NewGetAllCmd(GlobalState),
		commands.NewUpdateCmd(GlobalState),
	)
}

// completer is a function that returns suggestions for the prompt based on the input text.
func completer(d prompt.Document) []prompt.Suggest {
	w := d.GetWordBeforeCursor()
	input := d.TextBeforeCursor()

	// Only show suggestions if no space is typed yet
	if !strings.Contains(input, " ") {
		commands := []prompt.Suggest{
			{Text: "set", Description: "Set a key-value pair"},
			{Text: "get", Description: "Get a value by key"},
			{Text: "delete", Description: "Delete a key"},
			{Text: "flush", Description: "Flush the entire state"},
			{Text: "all", Description: "Retrieve all keys and values"},
			{Text: "exit", Description: "Exit the program"},
		}
		return prompt.FilterHasPrefix(commands, w, true)
	}

	// Return empty suggestions after a command is typed
	return []prompt.Suggest{}
}

// executor is a function that executes the input command. It is called when the user presses
func executor(input string) {
	input = strings.TrimSpace(input)
	
	if input == "" {
		return
	}

	if input == "exit" {
		fmt.Println("\nExiting Vertex...")
		persistance.Save(GlobalState)
		os.Exit(0)
	}

	// Split the input into command and arguments
	args := strings.Split(input, " ")
	// Set the arguments for the root command
	rootCmd.SetArgs(args)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

// Execute is the main function that runs the CLI application. It reads input from the user
func Execute(GlobalState *state.State) {
	// Start a goroutine to periodically refresh the global state.
	go refreshState()

	// Start the prompt
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("vertex > "),
		prompt.OptionTitle("Vertex CLI"),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionDescriptionBGColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSelectedDescriptionBGColor(prompt.LightGray),
	)

	p.Run()
}

// The `refreshState` function periodically loads and updates the global state from persistence.
func refreshState() {
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for range ticker.C {
		state, err := persistance.Load()
		if err != nil {
			vertex_log.Log("Error while loading state: "+err.Error(), "ERROR")
		} else {
			GlobalState = state
		}
	}
}

func main() {
	// Handle SIGINT and SIGTERM signals to save the state before exiting
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	
	go func() {
		<-c
		persistance.Save(GlobalState)
		os.Exit(1)
	}()

	Execute(GlobalState)
}
