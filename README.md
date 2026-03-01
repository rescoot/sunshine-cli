# sunshine-cli

Command-line interface for [Sunshine](https://github.com/rescoot/sunshine), the rescoot scooter management platform.

## Install

```bash
# From source
go install github.com/rescoot/sunshine-cli@latest

# Or build locally
git clone git@github.com:rescoot/sunshine-cli.git
cd sunshine-cli
make build
```

Pre-built binaries for Linux, macOS, and Windows are available on the [releases page](https://github.com/rescoot/sunshine-cli/releases).

## Authentication

The CLI uses OAuth2 with PKCE to authenticate against your Sunshine instance:

```bash
sunshine auth login              # Opens browser for OAuth flow
sunshine auth status             # Show current auth state
sunshine auth logout             # Clear stored credentials
```

Tokens are stored in `~/.config/sunshine/tokens.json`.

## Configuration

```bash
sunshine config                          # Show current config
sunshine config set server https://...   # Set server URL
sunshine config set default_scooter 3    # Skip the [id] argument
sunshine config path                     # Print config file location
```

Config file: `~/.config/sunshine/config.yaml`. Override per-command with `--server` and `--scooter` flags.

## Scooter ID resolution

The `[id]` argument is optional. The CLI resolves the scooter ID in order:

1. Positional argument
2. `--scooter` flag
3. `default_scooter` from config
4. Auto-detect (if you have exactly one scooter)

## Usage

```bash
# List scooters
sunshine scooters list [--limit 20] [--offset 0]
sunshine scooters show [id]

# Control commands
sunshine lock [id]
sunshine unlock [id]
sunshine honk [id]
sunshine blinkers [id] <left|right|both|off>
sunshine seatbox [id]
sunshine ping [id]
sunshine state [id]
sunshine locate [id]
sunshine hibernate [id]

# Alarm system
sunshine alarm [id]                      # Show alarm state
sunshine alarm arm [id]                  # Arm alarm
sunshine alarm disarm [id]              # Disarm alarm
sunshine alarm trigger [id] [--duration 5s]  # Sound alarm
sunshine alarm stop [id]                 # Silence active alarm

# Navigation
sunshine navigate [id] <lat> <lng> [title]   # Set destination
sunshine navigate show [id]                   # Show current destination
sunshine navigate clear [id]                  # Clear destination

# Trips
sunshine trips [id] [--limit 20] [--offset 0]
```

All commands accept `--json` for machine-readable output.

## License

[AGPL-3.0](LICENSE)
