package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AbdessamadEnabih/Vertex/internal/persistance"
	vertex_log "github.com/AbdessamadEnabih/Vertex/pkg/log"
	"github.com/AbdessamadEnabih/Vertex/pkg/state"
	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

// Define the root command and subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "vertex",
		Short: "Run vertex commands",
	}

	setCmd = &cobra.Command{
		Use:       "set",
		Short:     "Set a key-value pair",
		Example:   `set key value`,
		ValidArgs: []string{"key", "value"},
		Args:      cobra.ExactArgs(2),
		Run:       set,
	}

	getCmd = &cobra.Command{
		Use:       "get",
		Short:     "Get a value by key",
		Example:   `get key`,
		ValidArgs: []string{"key"},
		Args:      cobra.ExactArgs(1),
		Run:       get,
	}

	deleteCmd = &cobra.Command{
		Use:       "delete",
		Short:     "Delete a key",
		Example:   `delete key`,
		ValidArgs: []string{"key"},
		Args:      cobra.ExactArgs(1),
		Run:       delete,
	}

	flushCmd = &cobra.Command{
		Use:     "flush",
		Short:   "Flush the entire state",
		Example: `flush`,
		Run:     flush,
	}

	getAllCmd = &cobra.Command{
		Use:     "all",
		Short:   "Retrieve all keys and values",
		Example: `all`,
		Run:     getAll,
	}
)

// GlobalState is a pointer to the global state of the application. It is used to store
// key-value pairs in memory.
var GlobalState *state.State

// refreshInterval is the interval at which the global state is refreshed from persistence.
const refreshInterval = 60 * time.Second

func init() {
	// Load the global state from persistence. The function returns two values - the loaded
	// state and an error. In this case, the underscore `_` is used to ignore the error value,
	// assuming that the operation will succeed without any errors.
	GlobalState, _ = persistance.Load()

	// Add the subcommands to the root command
	rootCmd.AddCommand(setCmd, getCmd, deleteCmd, flushCmd, getAllCmd)
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
		fmt.Printf("Unknown Command: %v\n", err)
		fmt.Printf("Run 'vertex --help' for usage.\n")
	}
}

// Execute is the main function that runs the CLI application. It reads input from the user
func Execute(GlobalState *state.State) {
	// Start a goroutine to periodically refresh the global state.
	go refreshState()

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

func validateArgs(args []string, expectedArgs int) {
	if len(args) != expectedArgs {
		fmt.Printf("Expected %d arguments, got %d\n", expectedArgs, len(args))
		return
	}
}

// The following functions are the handlers for the CLI commands.
func set(cmd *cobra.Command, args []string) {
	validateArgs(args, 2)

	err := GlobalState.Set(args[0], args[1])
	if err != nil {
		fmt.Printf("Unable to set the key %v: %v\n", args[0], err)
	}
}

func get(cmd *cobra.Command, args []string) {
	validateArgs(args, 1)

	value, err := GlobalState.Get(args[0])
	if err != nil {
		fmt.Printf("Failed to retrieve the key %s :  %v\n", args[0], err)
	} else {
		fmt.Println("Value:", value)
	}
}

func delete(cmd *cobra.Command, args []string) {
	validateArgs(args, 1)

	err := GlobalState.Delete(args[0])
	if err != nil {
		fmt.Printf("Failed to delete key %s : %v\n", args[0], err)
	}
}

func flush(cmd *cobra.Command, args []string) {
	err := GlobalState.FlushAll()
	if err != nil {
		fmt.Printf("Failed to flush data: %v\n", err)
	}
}

func getAll(cmd *cobra.Command, args []string) {
	values := GlobalState.GetAll()
	fmt.Println("All keys:", values)
}

func main() {
	// Execute the CLI commands with the loaded global state.
	Execute(GlobalState)
}
