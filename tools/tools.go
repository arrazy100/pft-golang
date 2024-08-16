package tools

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// ArgumentHandler defines the function signature for argument handler functions.
type ArgumentHandler func(args []string) error

// ArgumentRegistry holds the registered arguments and their handlers.
type ArgumentRegistry struct {
	handlers map[string]ArgumentHandler
}

// NewArgumentRegistry creates a new ArgumentRegistry.
func _newArgumentRegistry() *ArgumentRegistry {
	return &ArgumentRegistry{
		handlers: make(map[string]ArgumentHandler),
	}
}

// Register adds a new argument and its handler to the registry.
func (r *ArgumentRegistry) _register(arg string, handler ArgumentHandler) {
	r.handlers[arg] = handler
}

// Parse checks the command-line arguments and invokes the appropriate handler.
func (r *ArgumentRegistry) _parse() error {
	args := os.Args[1:] // Skip the program name

	for i := 0; i < len(args); i++ {
		if handler, exists := r.handlers[args[i]]; exists {
			// Pass the remaining arguments to the handler
			err := handler(args[i+1:])
			if err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func _compileProtoFiles(protoDir string) error {
	// Walk through the directory and subdirectories to find all .proto files
	err := filepath.Walk(protoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file has a .proto extension
		if filepath.Ext(path) == ".proto" {
			fmt.Printf("Compiling %s\n", path)

			// Compile the .proto file using protoc
			cmd := exec.Command("protoc", "--proto_path=.", "--go_out=.", "--go-grpc_out=.", path)
			output, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("failed to compile %s: %v\n%s", path, err, string(output))
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to compile proto files: %v", err)
		return err
	}

	return nil
}

func ParseCommand() bool {
	registry := _newArgumentRegistry()

	// Register the --compile-proto argument
	registry._register("--compile-proto", func(args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing directory argument for --compile-proto")
		}

		protoDir := args[0]
		err := _compileProtoFiles(protoDir)
		if err != nil {
			return fmt.Errorf("failed to compile proto files: %v", err)
		}

		fmt.Println("Proto files compiled successfully.")
		return nil
	})

	// Register more arguments and handlers as needed
	// registry.Register("--another-arg", anotherHandler)

	// Parse the arguments
	err := registry._parse()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	args := os.Args[1:]

	return len(args) != 0
}
