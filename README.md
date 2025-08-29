# tmux-githop

Fast tmux session hopping between git repos

## Features

- Interactive repository picker with fuzzy search
- Quick tmux session switching
- Automatic session creation and management
- Works with GHQ-managed repositories

## Installation

### Package managers

```bash
# Debian/Ubuntu
sudo apt install ./tmux-githop_*_amd64.deb

# Red Hat/Fedora  
sudo dnf install ./tmux-githop_*_x86_64.rpm

# Alpine
sudo apk add --allow-untrusted ./tmux-githop_*_x86_64.apk
```

#### Nix/NixOS

##### Run directly

`nix run github:callumio/tmux-githop`

##### Flake input

```nix
{
  inputs = {
    tmux-githop.url = "github:callumio/tmux-githop";
  };

  outputs = { self, tmux-githop, ... }: {
    # Use tmux-githop.packages.${system}.default in your configurations
  };
}
```

### Pre-built binaries

```bash
# Download latest release from https://github.com/callumio/tmux-githop/releases
# Choose the appropriate archive for your system (linux/darwin, amd64/arm64)
tar xzf tmux-githop_*_*.tar.gz
sudo mv tmux-githop /usr/local/bin/
```

### Build from source

```bash
git clone https://github.com/callumio/tmux-githop.git
cd tmux-githop && go build ./cmd/tmux-githop
```

## Development

```bash
git clone https://github.com/callumio/tmux-githop.git
cd tmux-githop

# Use Nix for full dev environment (recommended)
nix develop

# Or manually
go mod download
go build ./cmd/tmux-githop
go test ./...
```

### Contributing

1. Fork and clone
2. Create feature branch
3. Make changes and test
4. Submit pull request

### Releases

Automated with GoReleaser. Push a version tag to create a release:

```bash
git tag v1.0.0 && git push origin v1.0.0
```

## License

See [LICENSE](LICENSE) file for details.
