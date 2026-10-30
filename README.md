# ClipDash

ClipDash is a small command-line note manager that works with your clipboard. Copy text, save it as a note with a name you choose, and retrieve it later without needing to remember where the note is stored.

## Requirements

- Go 1.27.3 or later to run from source or build the program.
- A working system clipboard for creating notes and copying notes back to it.

## Usage

Run commands from the project directory:

```powershell
go run . create meeting-notes
go run . list
go run . read meeting-notes
```

The commands are:

| Command | Description |
| --- | --- |
| `create <name>` | Saves the current text on your clipboard to a new `<name>.txt` note. |
| `list` | Prints the files in your notes folder. |
| `read <name>` | Copies the contents of `<name>.txt` to your clipboard. |

For example, copy some text to your clipboard and run `go run . create shopping-list`. Later, run `go run . read shopping-list` to put that note's contents back on your clipboard so you can paste them anywhere.

Note names are single file names, not paths. ClipDash adds the `.txt` extension automatically. It will not overwrite a note that already exists with the same name.

## Where notes are stored

ClipDash creates a `notes` folder inside your operating system's user configuration directory:

- **Windows:** `%AppData%\notes`
- **macOS:** `~/Library/Application Support/notes`
- **Linux:** `$XDG_CONFIG_HOME/notes`, or `~/.config/notes` when `XDG_CONFIG_HOME` is not set

You do not need to create this folder yourself; ClipDash creates it when saving the first note.

## Build

To compile the program for your current operating system:

```powershell
go build -o clipdash.exe .
```

Then run it on Windows like this:

```powershell
.\clipdash.exe create meeting-notes
.\clipdash.exe list
.\clipdash.exe read meeting-notes
```

On macOS or Linux, build without the `.exe` suffix:

```sh
go build -o clipdash .
```

Then run it with `./clipdash`, for example `./clipdash create meeting-notes`.
