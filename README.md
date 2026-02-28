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
# Set your Sunshine server (default: https://sunshine.rescoot.org)
sunshine config set server https://sunshine.example.com

# Set a default scooter so you don't need to pass an ID every time
sunshine config set default_scooter 1

sunshine config show
```

Config is stored in `~/.config/sunshine/config.yaml`.

## Usage

```bash
# List scooters
sunshine scooters list
sunshine scooters show <id>

# Control commands
sunshine lock [id]
sunshine unlock [id]
sunshine honk [id]
sunshine blinkers [id] <left|right|both|warning|off>
sunshine seatbox [id]
sunshine ping [id]
sunshine state [id]
sunshine locate [id]
sunshine alarm [id] arm|disarm|trigger
sunshine hibernate [id]

# Trips
sunshine trips [id] [--limit N]

# Navigation destination
sunshine destination show [id]
sunshine destination set [id] --lat LAT --lng LNG [--address ADDR]
sunshine destination clear [id]
```

All commands accept `--json` for machine-readable output.

## License

[AGPL-3.0](LICENSE)
