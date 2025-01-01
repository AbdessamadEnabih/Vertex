package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AbdessamadEnabih/Vertex/internal/persistance"
	vertex_log "github.com/AbdessamadEnabih/Vertex/pkg/log"
	"github.com/AbdessamadEnabih/Vertex/pkg/state"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vertex",
	Short: "Run vertex commands",
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a key-value pair",
	Run:   set,
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a value by key",
	Run:   get,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a key",
	Run:   delete,
}

var flushCmd = &cobra.Command{
	Use:   "flush",
	Short: "Flush the entire state",
	Run:   flush,
}

var getAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Retrieve all keys and values",
	Run:   getAll,
}

var GlobalState *state.State

const refreshInterval = 60

func main() {
	// Load the global state from persistence. The function returns two values - the loaded
	// state and an error. In this case, the underscore `_` is used to ignore the error value,
	// assuming that the operation will succeed without any errors.
	GlobalState, _ = persistance.Load()

	// Execute the CLI commands with the loaded global state.
	Execute(GlobalState)
}

func Execute(GlobalState *state.State) {
	// Add subcommands to the root command.
	rootCmd.AddCommand(setCmd, getCmd, deleteCmd, flushCmd, getAllCmd)

	// Start a goroutine to periodically refresh the global state.
	go refreshState()
    
	// Create a new reader to read input from the standard input (stdin).
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("vertex > ")

		// Read a line of input from the user.
		input, err := reader.ReadString('\n')
		if err != nil {
			// If the error is due to EOF (end of file), exit the program.
			if err.Error() == "EOF" {
				fmt.Println("\nExiting Vertex...")
				break
			}

			fmt.Println("Error reading input:", err)
			continue
		}
	    
		// Trim any leading or trailing whitespace from the input.
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		if input == "exit" {
			fmt.Println("\nExiting Vertex...")
			break
		}

		// Split the input into command and arguments
		args := strings.Split(input, " ")
		// Set the arguments for the root command
		rootCmd.SetArgs(args)

		// Execute the root command
		if err := rootCmd.Execute(); err != nil {
			fmt.Printf("Unkown Command: %v\n", err)
			fmt.Printf("Run 'vertex --help' for usage.\n")
		}
	}

	defer persistance.Save(GlobalState)
}

// The `refreshState` function periodically loads and updates the global state from persistence.
func refreshState() {
	ticker := time.NewTicker(refreshInterval * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		state, err := persistance.Load()
		if err != nil {
			vertex_log.Log("Error while loading state: " + err.Error(), "ERROR")
		} else {
			GlobalState = state
		}
	}
}

func set(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		cmd.Usage()
		return
	}
	err := GlobalState.Set(args[0], args[1])
	if err != nil {
		fmt.Printf("Unable to set the key %v: %v\n", args[0], err)
	}
}

func get(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		return
	}
	value, err := GlobalState.Get(args[0])
	if err != nil {
		fmt.Printf("Failed to retrieve the key %s :  %v\n", args[0], err)
	} else {
		fmt.Println("Value:", value)
	}
}

func delete(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		return
	}
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
