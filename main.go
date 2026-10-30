package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/atotto/clipboard"
)

var validArgs = []string{"create", "read", "list"} //the list of args that the program allows you to use

func ValidateArg(arg string) bool {
	return slices.Contains(validArgs, strings.ToLower(arg))
}

// The function that does the work for you based on the arg you pass in.
func DoWorkForArg(arg string) {
	lowerArg := strings.ToLower(arg)
	switch lowerArg {
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a name for the note")
			return
		}
		if err := createNote(os.Args[2]); err != nil {
			fmt.Println("Could not create note:", err)
			return
		}
		fmt.Println("Note created")
	case "read":
		if len(os.Args) < 3 {
			fmt.Println("Please provide the name of the note to read")
			return
		}
		if err := readNote(os.Args[2]); err != nil {
			fmt.Println("Could not read note:", err)
			return
		}
		fmt.Println("Note copied to clipboard")
	case "list":
		if err := listNotes(); err != nil {
			fmt.Println("Could not list notes:", err)
		}
	default:
		fmt.Println("Unknown argument")
	}
}

func notesDirectory() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user configuration directory: %w", err)
	}
	return filepath.Join(configDir, "notes"), nil
}

func notePath(name string) (string, error) {
	if name == "" || name == "." || name == ".." || filepath.Base(name) != name {
		return "", fmt.Errorf("invalid note name %q", name)
	}
	dir, err := notesDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name+".txt"), nil
}

func createNote(name string) error {
	content, err := clipboard.ReadAll()
	if err != nil {
		return fmt.Errorf("read text from clipboard: %w", err)
	}
	path, err := notePath(name)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create notes directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("note %q already exists", name)
		}
		return fmt.Errorf("create note file: %w", err)
	}
	if _, err := file.WriteString(content); err != nil {
		closeErr := file.Close()
		removeErr := os.Remove(path)
		if closeErr != nil || removeErr != nil {
			return fmt.Errorf("write note file: %w (cleanup close error: %v; remove error: %v)", err, closeErr, removeErr)
		}
		return fmt.Errorf("write note file: %w", err)
	}
	if err := file.Close(); err != nil {
		if removeErr := os.Remove(path); removeErr != nil {
			return fmt.Errorf("close note file: %w (remove incomplete file: %v)", err, removeErr)
		}
		return fmt.Errorf("close note file: %w", err)
	}
	return nil
}

func listNotes() error {
	dir, err := notesDirectory()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No notes found")
			return nil
		}
		return fmt.Errorf("read notes directory: %w", err)
	}

	found := false
	for _, entry := range entries {
		if !entry.IsDir() {
			fmt.Println(entry.Name())
			found = true
		}
	}
	if !found {
		fmt.Println("No notes found")
	}
	return nil
}

func readNote(name string) error {
	path, err := notePath(name)
	if err != nil {
		return err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read note file: %w", err)
	}
	if err := clipboard.WriteAll(string(content)); err != nil {
		return fmt.Errorf("copy note to clipboard: %w", err)
	}
	return nil
}

func main() {
	if len(os.Args) == 1 {
		panic("Not enough args")
	}
	res := ValidateArg(os.Args[1])
	if res {
		DoWorkForArg(os.Args[1])
	} else {
		fmt.Println("arg ", os.Args[1], "is not a valid argument")
	}
}
