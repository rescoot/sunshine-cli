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
sunshine config set default_scooter 3    # Set default scooter
sunshine config path                     # Print config file location
```

Config file: `~/.config/sunshine/config.yaml`. Override per-command with `--server` and `--scooter` flags.

## Scooter selection

The scooter is resolved in order:

1. `--scooter` flag
2. `default_scooter` from config
3. Auto-detect (if you have exactly one scooter)

## Usage

```bash
# Scooter status
sunshine status                          # Full scooter details

# List scooters
sunshine scooters list [--limit 20] [--offset 0]
sunshine scooters show

# Control commands
sunshine lock
sunshine unlock
sunshine honk
sunshine blinkers <left|right|both|off>
sunshine seatbox
sunshine ping
sunshine state                           # Request telemetry refresh
sunshine locate
sunshine hibernate

# Alarm system
sunshine alarm                           # Show alarm state
sunshine alarm arm                       # Arm alarm
sunshine alarm disarm                    # Disarm alarm
sunshine alarm trigger [--duration 5s]   # Sound alarm
sunshine alarm stop                      # Silence active alarm

# Navigation
sunshine navigate <lat> <lng> [title]    # Set destination
sunshine navigate show                   # Show current destination
sunshine navigate clear                  # Clear destination

# Trips
sunshine trips list [--limit 20] [--offset 0]
sunshine trips show <trip-id>
```

All commands accept `--json` for machine-readable output.

## Raw API access

```bash
sunshine api /scooters                              # GET
sunshine api /scooters/3                            # GET with path
sunshine api /scooters/3/trips?limit=5              # GET with query params
sunshine api /scooters/3/lock -X POST               # POST
sunshine api /scooters/3/blinkers -d '{"state":"left"}'  # POST with JSON body
sunshine api /scooters/3/alarm -X POST -d '{"duration":"10s"}'
```

Paths are relative to `/api/v1/`. Method defaults to GET, or POST when `-d` is given.

## Man pages

```bash
make man              # Generate man pages to man/
make install-man      # Install to system man path
man sunshine          # View main man page
```

## License

[AGPL-3.0](LICENSE)
