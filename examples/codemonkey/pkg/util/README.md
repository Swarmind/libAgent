## Package: `main`

This package defines a command-line application for processing files based on user input and environment variables. It reads file paths from the command line, checks if they exist, and performs operations like printing their contents or calculating checksums. The program exits with an error message if any required arguments are missing or invalid.

**Imports:**

*   `fmt`: For formatted output to the console.
*   `os`: For interacting with the operating system (command-line arguments, file existence checks).
*   `crypto/md5`: For calculating MD5 checksums of files.
*   `io`: For reading from files.
*   `log`: For logging errors and exiting the program.

**Environment Variables:**

The application uses `DEBUG` environment variable to control verbose output. If set to "true", it prints additional debugging information during file processing.

**Command-Line Arguments:**

The program expects one or more file paths as command-line arguments. It checks if each provided path exists before attempting to process the file. If a file does not exist, an error message is printed, and the program exits with a non-zero exit code.

**Function Summary:**

### `main()`

This function parses command-line arguments, iterates through them, and performs actions based on whether each file exists:
1.  If DEBUG environment variable is set to "true", it prints verbose output about processing each file.
2.  It checks if the provided file path exists using `os.Stat()`. If not, it logs an error message with `log.Fatalf()` and exits.
3.  For existing files, it reads their contents into a byte slice using `io.ReadFile()`.
4.  It calculates the MD5 checksum of the file's content using `crypto/md5` package.
5.  Finally, it prints the file path along with its calculated MD5 hash to standard output.

**Edge Cases:**

*   If no command-line arguments are provided, the program exits with an error message.
*   The application does not handle symbolic links or other special file types beyond basic existence checks.
*   Large files may consume significant memory when read into a byte slice using `io.ReadFile()`. Consider streaming for very large files to avoid memory issues.

**Project Package Structure:**

-   `main.go`

**Relations Between Code Entities:**

The main function orchestrates the entire process, relying on standard library packages like `os`, `crypto/md5`, and `io` to perform file system operations and checksum calculations. The DEBUG environment variable controls verbosity without altering core functionality. Error handling is done via `log.Fatalf()`, which terminates execution upon failure.