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
	log "github.com/AbdessamadEnabih/Vertex/pkg/log"
	"github.com/AbdessamadEnabih/Vertex/pkg/datastore"
	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

// Define the root command and subcommands
var rootCmd = &cobra.Command{
	Use:   "vertex",
	Short: "Run vertex commands",
}

// GlobalDataStore is a pointer to the global datastore of the application. It is used to store
// key-value pairs in memory.
var GlobalDataStore *datastore.DataStore

// refreshInterval is the interval at which the global datastore is refreshed from persistence.
const refreshInterval = 60 * time.Second

// init initializes the CLI by loading the global datastore and adding the commands to the root command.
func init() {
	GlobalDataStore, _ = persistance.Load()
	rootCmd.AddCommand(
		commands.NewGetAllCmd(GlobalDataStore),
		commands.NewGetCmd(GlobalDataStore),
		commands.NewSetCmd(GlobalDataStore),
		commands.NewUpdateCmd(GlobalDataStore),
		commands.NewDeleteCmd(GlobalDataStore),
		commands.NewFlushCmd(GlobalDataStore),
	)
}

// completer is a function that returns suggestions for the prompt based on the input text.
func completer(d prompt.Document) []prompt.Suggest {
	// Get the word before the cursor and the text before the cursor
	w := d.GetWordBeforeCursor()
	input := d.TextBeforeCursor()

	// Only show suggestions if no space is typed yet
	if !strings.Contains(input, " ") {
		commands := []prompt.Suggest{}

		// Add the commands to the suggestions
		for _, cmd := range rootCmd.Commands() {
			commands = append(commands, prompt.Suggest{Text: cmd.Use, Description: cmd.Short})
		}

		commands = append(commands, prompt.Suggest{Text: "exit", Description: "Exit the program"})

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
		persistance.Save(GlobalDataStore)
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
func Execute(GlobalDataStore *datastore.DataStore) {
	// Start a goroutine to periodically refresh the global datastore.
	go refreshDataStore()

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

// The `refreshDataStore` function periodically loads and updates the global datastore from persistence.
func refreshDataStore() {
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for range ticker.C {
		datastore, err := persistance.Load()
		if err != nil {
			log.Log("Error while loading datastore: "+err.Error(), "ERROR")
		} else {
			GlobalDataStore = datastore
		}
	}
}

func main() {
	// Handle SIGINT and SIGTERM signals to save the datastore before exiting
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		<-c
		persistance.Save(GlobalDataStore)
		os.Exit(1)
	}()

	Execute(GlobalDataStore)
}
