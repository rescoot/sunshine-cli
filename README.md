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

Edit `~/.config/sunshine/config.yaml` to set defaults:

```yaml
server: https://sunshine.rescoot.org    # default
default_scooter: 3                      # skip the [id] argument
```

You can also override per-command with `--server` and `--scooter` flags.

## Usage

```bash
# List scooters
sunshine scooters list
sunshine scooters show <id>

# Control commands (id is optional if default_scooter is set)
sunshine lock [id]
sunshine unlock [id]
sunshine honk [id]
sunshine blinkers [id] <left|right|both|off>
sunshine seatbox [id]
sunshine ping [id]
sunshine state [id]
sunshine locate [id]
sunshine alarm [id] [--duration 5s]
sunshine hibernate [id]

# Trips
sunshine trips [id]

# Navigation destination
sunshine destination show [id]
sunshine destination set [id] <lat> <lng> [--address ADDR]
sunshine destination clear [id]
```

All commands accept `--json` for machine-readable output.

## License

[AGPL-3.0](LICENSE)
