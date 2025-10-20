# i3-scripts-go

A comprehensive collection of Go-based scripts for enhancing the i3 window manager experience. This is a work-in-progress golang version of [i3-scripts](https://github.com/modernpacifist/i3-scripts).

## Overview

This project provides a suite of utilities to extend i3 window manager functionality with optimized Go implementations. Each script is designed to handle specific window management tasks, from basic operations like killing containers to advanced features like floating window management and workspace manipulation.

## Dependencies

**Debian 12:**
- VolumeControl: `pulseaudio` (install: `apt install pulseaudio`)
- VolumeControl: `pactl` (install: `apt install pulseaudio-utils`)
- General: `i3-wm` and `i3-input` for i3 window manager integration
- Optional: `upx-ucl` for binary compression (install: `apt install upx-ucl`)

## Build Instructions

```bash
# Build all scripts with optimizations
make build

# Build with maximum optimizations
make build-optimized

# Build with compression (requires UPX)
make build-compressed

# Install to system (requires root)
sudo make install

# Clean build artifacts
make clean
```

## Scripts Documentation

### Window Management Scripts

#### 1. **BackAndForthContainers** (`cmd/back_and_forth_containers/`)
**Purpose**: Implements Alt-Tab like functionality for i3 containers, allowing quick switching between the current and previously focused containers.

**Features**:
- Maintains a history of the last 3 focused containers
- Can run as a daemon to continuously track focus changes
- Stores container history in `~/.BackAndForthContainer.json`
- Uses TCP port 63333 for daemon communication

**Usage**:
```bash
# Switch to previously focused container
./bin/BackAndForthContainers

# Run as daemon to track focus changes
./bin/BackAndForthContainers -daemon
```

#### 2. **KillContainer** (`cmd/kill_container/`)
**Purpose**: Terminates the currently focused container/window.

**Usage**:
```bash
./bin/KillContainer
```

#### 3. **LockContainer** (`cmd/lock_container/`)
**Purpose**: Locks or prevents interaction with the currently focused container.

**Usage**:
```bash
./bin/LockContainer
```

#### 4. **MarkContainer** (`cmd/mark_container/`)
**Purpose**: Adds marks to containers for easy identification and navigation.

**Usage**:
```bash
./bin/MarkContainer
```

### Floating Window Management Scripts

#### 5. **ResizeFloatContainer** (`cmd/resize_float_container/`)
**Purpose**: Resizes floating windows using vim-style directional keys.

**Features**:
- Supports vim-style directions: h (left), j (down), k (up), l (right), w (widen)
- Accepts positive or negative values for expansion/contraction
- Input validation for direction and value format

**Usage**:
```bash
# Resize floating window left by 50 pixels
./bin/ResizeFloatContainer resize h -50

# Resize floating window down by 30 pixels
./bin/ResizeFloatContainer resize j +30
```

#### 6. **MoveFloatContainer** (`cmd/move_float_container/`)
**Purpose**: Moves floating containers to predefined positions.

**Usage**:
```bash
# Move to position (numeric argument determines position)
./bin/MoveFloatContainer position <position_number>
```

#### 7. **ManageFloatContainer** (`cmd/manage_float_container/`)
**Purpose**: Advanced floating container management with save/restore functionality.

**Features**:
- Save current container state with marks
- Restore containers by mark
- Show/hide containers by mark
- Update container states

**Usage**:
```bash
# Save current container state
./bin/ManageFloatContainer --save

# Restore container by mark
./bin/ManageFloatContainer --restore <mark>

# Show container by mark
./bin/ManageFloatContainer --show <mark>

# Update container by mark
./bin/ManageFloatContainer --update <mark>
```

#### 8. **DiagonalResize** (`cmd/diagonal_resize/`)
**Purpose**: Resizes floating windows diagonally (both width and height simultaneously).

**Usage**:
```bash
# Increase size diagonally by 20 pixels
./bin/DiagonalResize size +20

# Decrease size diagonally by 15 pixels
./bin/DiagonalResize size -15
```

#### 9. **MarginResize** (`cmd/margin_resize/`)
**Purpose**: Resizes windows by adjusting their margins/borders in specific directions.

**Usage**:
```bash
# Resize border in specified direction
./bin/MarginResize border <direction>
```

### Workspace Management Scripts

#### 10. **SwapWorkspaces** (`cmd/swap_workspaces/`)
**Purpose**: Swaps the contents of two workspaces.

**Usage**:
```bash
./bin/SwapWorkspaces
```

#### 11. **RenameWorkspace** (`cmd/rename_workspace/`)
**Purpose**: Provides an interactive interface to rename the current workspace.

**Usage**:
```bash
./bin/RenameWorkspace
```

#### 12. **ChangeWorkspaceIndex** (`cmd/change_workspace_index/`)
**Purpose**: Changes the current workspace to a specified index number.

**Usage**:
```bash
# Switch to workspace 5
./bin/ChangeWorkspaceIndex index 5
```

### System Control Scripts

#### 13. **VolumeControl** (`cmd/volume_control/`)
**Purpose**: Comprehensive audio volume management using PulseAudio.

**Features**:
- Toggle mute/unmute
- Adjust volume with +/- values
- Round volume to nearest 5
- Maximum volume limit (100 by default)

**Usage**:
```bash
# Toggle mute/unmute
./bin/VolumeControl toggle

# Increase volume by 10
./bin/VolumeControl adjust +10

# Decrease volume by 5
./bin/VolumeControl adjust -5

# Round volume to nearest 5
./bin/VolumeControl round
```

#### 14. **ChangeMonitorBrightness** (`cmd/change_monitor_brightness/`)
**Purpose**: Adjusts monitor brightness levels.

**Usage**:
```bash
# Increase brightness by 0.1
./bin/ChangeMonitorBrightness --change 0.1

# Decrease brightness by 0.2
./bin/ChangeMonitorBrightness --change -0.2
```

#### 15. **KeyboardLayout** (`cmd/keyboard_layout/`)
**Purpose**: Cycles through multiple keyboard layouts.

**Features**:
- Supports multiple layouts separated by forward slashes
- Cycles through layouts in sequence

**Usage**:
```bash
# Cycle through US and Russian layouts
./bin/KeyboardLayout cycle us/ru

# Cycle through multiple layouts
./bin/KeyboardLayout cycle us/ru/de/fr
```

### Utility Scripts

#### 16. **StickyToggle** (`cmd/sticky_toggle/`)
**Purpose**: Toggles the sticky state of the focused container (makes it appear on all workspaces).

**Usage**:
```bash
./bin/StickyToggle
```

#### 17. **FocusContainer** (`cmd/focus_container/`)
**Purpose**: Placeholder for future focus management functionality.

**Status**: Currently contains TODO comments for planned functionality to focus/hide floating containers by mark.

## Library Components

### Internal Libraries (`internal/`)

#### `internal/i3scripts/common.go`
Common utilities for i3 operations including:
- Tree and workspace retrieval functions
- Node and output management
- Notification system integration
- i3-input dialog utilities
- Command execution helpers

#### `internal/config/common.go`
Configuration management utilities for various scripts.

### Public Libraries (`pkg/`)

#### `pkg/i3scripts/i3scripts.go`
Public API for i3 operations with similar functionality to internal version but with different error handling approaches.

#### `pkg/i3scripts/optimized.go`
Performance-optimized versions of common i3 operations.

## Performance Analysis

The project includes comprehensive performance analysis in `PERFORMANCE_ANALYSIS.md` covering:
- Build time optimizations
- Binary size reductions
- Runtime performance improvements
- Memory usage optimization

## Development

### Code Structure
- `cmd/`: Main application entry points
- `internal/`: Internal packages and utilities
- `pkg/`: Public packages for external use
- `Makefile`: Build automation and optimization

### Build Optimizations
The project uses advanced Go build optimizations:
- Link-time optimizations (`-ldflags="-s -w"`)
- Trimmed paths (`-trimpath`)
- Parallel builds
- Optional UPX compression

### Testing
```bash
# Test all binaries
make test

# Run build benchmarks
make benchmark

# Check binary sizes
make sizes
```

## Contributing

This is an active project with ongoing development. Many scripts contain TODO comments indicating planned improvements and features.

## License

This project follows the same license as the original [i3-scripts](https://github.com/modernpacifist/i3-scripts) project.
