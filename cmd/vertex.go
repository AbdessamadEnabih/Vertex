// cmd/vertex.go
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

func Execute(globalState *state.State) {
    GlobalState = globalState
    rootCmd.AddCommand(setCmd, getCmd, deleteCmd, flushCmd, getAllCmd)

    // Start interactive command loop
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("vertex> ")

        input, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading input:", err)
            continue
        }
        
        input = strings.TrimSpace(input)

        if input == "exit" {
            fmt.Println("Exiting Vertex...")
            break
        }

        // Split the input into command and arguments
        args := strings.Split(input, " ")
        // Set the arguments for the root command
        rootCmd.SetArgs(args)

        // Execute the root command
        if err := rootCmd.Execute(); err != nil {
            fmt.Printf("Error executing command: %v\n", err)
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
		fmt.Printf("Error setting key: %v\n", err)
	}
}

func get(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		return
	}
	value, err := GlobalState.Get(args[0])
	if err != nil {
		fmt.Printf("Error getting key: %v\n", err)
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
		fmt.Printf("Error deleting key: %v\n", err)
	}
}

func flush(cmd *cobra.Command, args []string) {
	err := GlobalState.FlushAll()
	if err != nil {
		fmt.Printf("Error flushing state: %v\n", err)
	}
}

func getAll(cmd *cobra.Command, args []string) {
	values := GlobalState.GetAll()
	fmt.Println("All keys:", values)
}
