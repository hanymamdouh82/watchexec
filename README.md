# WatchExec

WatchExec is a lightweight and configurable file watcher that monitors specified directories and executes predefined commands when changes are detected. It supports directory-specific intervals and allows excluding certain files and folders from monitoring.

## Features

- **Configurable directory watching** – Define directories to monitor in a config file.
- **Automated command execution** – Run specified commands when changes occur.
- **Custom watch intervals** – Each directory can have its own monitoring interval.
- **Exclusion support** – Ignore specified files or directories from triggering actions.
- **Lightweight and efficient** – Built with Go for minimal resource usage.

## Installation

```sh
# Clone the repository
git clone https://github.com/hanymamdouh82/watchexec.git
cd watchexec

# Build the binary
go build -o watchexec
```

## Usage

```sh
./watchexec -c config.yml
```

### Configuration File

WatchExec uses a `yml` configuration file to define watched directories, commands, and exclusions.

#### Example `config.json`

```yml
dirs:
  /home/user/docs:
    bin: ls # command / binary to excute once changes detected on the directory
    args: # list of arguments
      - -a
      - -l
    stdout: true # pipe result to stdout
    delay: 3 # delay between excution and watch rounds
    exclude: # exclusion list of dirs and files
      - .git # directory (no need to include /)
      - readme.md # file

  /home/user/Downloads:
    bin: rsync
    args:
      - myfile.txt
      - /home/user/backup
    delay: 10
    stdout: false
    exclude:
      - partial.download
```

- `dirs` – List of directories to monitor.
- `bin` – Command / binary to execute when a change is detected.
- `stdout` – pipe result to stdou
- `args` – list of command arguments
- `delay` – Time (in seconds) between checks.
- `exclude` – List of files or subdirectories to ignore.

## Command-Line Options

| Flag | Description                    | Default      |
| ---- | ------------------------------ | ------------ |
| `-c` | Path to the configuration file | `conf.yml` |
| `-v` | Verbose; logs all operations   |              |

## Contributing

Pull requests are welcome! If you have feature suggestions or find a bug, please open an issue.

## License

This project is licensed under the MIT License.

## Author

[Hany Mamdouh](https://github.com/hanymamdouh82)
