# sort-imports

The `sort-imports` tool automatically sorts Go import statements into three groups: standard library imports, 
external imports, and internal imports, based on configurable rules. This tool aims to improve code readability and 
maintain consistency across Go projects.

This README provides instructions on integrating `sort-imports` with the GoLand IDE, enabling automatic sorting 
of imports every time a Go file is saved.

## Installation

### Using `go install` (Recommended)

You can directly install the `sort-imports` tool using the `go install` command:

```sh
go install github.com/wundersasha/sort-imports@latest
```
This command will download and install the sort-imports executable to your `$GOPATH/bin` directory.

### Building and Installing Locally
If you prefer to build and install sort-imports from the source, follow these steps:

Clone the repository:

```shell
git clone https://github.com/wundersasha/sort-imports.git
cd sort-imports
```
Build the tool:

```shell
go build -o sort-imports .
```
This command compiles the program and generates an executable named sort-imports in the current directory.

#### (Optional) Move the executable to a directory in your PATH to run it from anywhere:

```shell
mv sort-imports /usr/local/bin/
```
**Note**: The destination directory (/usr/local/bin/) might differ based on your operating system and environment 
setup. Choose a directory that is in your PATH.

## Usage
Before using sort-imports, set the `GO_INTERNAL_PATH` environment variable to specify the import path of your 
internal packages. This helps sort-imports distinguish between external and internal imports.

```shell
export GO_INTERNAL_PATH="your/internal/package/path"`
```

### Goland IDE integration

Integrate sort-imports with GoLand by setting up a File Watcher that automatically runs the tool on Go file saves.

#### Setting Up a File Watcher
1. Open GoLand and navigate to **File | Settings** for Windows and Linux, or **GoLand | Settings** for macOS.
2. In the **Settings/Preferences** dialog, select **Tools | File Watchers**.
3. Click the + button to add a new watcher and select **custom** template from the list of templates.
4. Configure the new **File Watcher** with the following settings:
   * **Name**: `Sort Imports`
   * **File type**: `Go`
   * **Scope**: Choose **Current File** to apply the watcher to the file you're currently editing.
   * **Program**: Specify the path to the `sort-imports` executable. If you've installed it via `go install`, 
   it should be in your `$GOPATH/bin` directory. You can find the exact path by running `which sort-imports` (**Unix**) 
   or `where sort-imports` (**Windows**) in your terminal.
   * **Arguments**: Enter `$FilePath$` to pass the current file path to the sort-imports tool.
   * **Output paths**: Leave this field empty, as `sort-imports` modifies the file in place.
   * **Working directory**: Use `$ProjectFileDir$` to ensure the tool runs in the context of your project directory.
   * **Environment variables**: Set `GO_INTERNAL_PATH` environment variable here to ensure the tool would work correctly.
   * In the **Advanced Options** section, adjust the settings according to your preferences. It's recommended 
   to check **Auto-save edited files to trigger the watcher** for seamless integration.
5. Click **OK** or **Apply** to save the new **File Watcher**.
6. It is recommended to use **Project** scope for this **File Watcher** to ensure the `GO_INTERNAL_PATH` environment
variable is specific for every single project.

With the **File Watcher** set up, `sort-imports` will automatically run every time you save a Go file in GoLand, 
sorting the import statements according to the predefined rules.

### Local usage

To sort the imports in a local Go file, simply run:

```shell
sort-imports path/to/yourfile.go
```

## Contributing
Contributions to `sort-imports` are welcome! Whether it's bug reports, feature requests, or code contributions, 
please feel free to open an issue or submit a pull request on GitHub.

## Fork the repository.
* Create a new branch for your feature or fix.
* Commit your changes.
* Push to the branch.
* Submit a pull request.

## License
`sort-imports` is open-sourced software licensed under the **MIT license**.
