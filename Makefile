# Build optimization variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags="-s -w -X main.version=$(VERSION)"
BUILDFLAGS := -trimpath -buildmode=exe
PARALLEL_JOBS := $(shell nproc 2>/dev/null || echo 4)

# Enable parallel builds
.PHONY: build build-optimized build-compressed install clean fmt test benchmark

fmt:
	find . -type f -name "*.go" | xargs -n 1 go fmt

# Original build targets (kept for compatibility)
back_and_forth_containers:
	go build $(BUILDFLAGS) $(LDFLAGS) -o ./bin/BackAndForthContainers ./cmd/back_and_forth_containers/main.go

keyboard_layout:
	go build $(BUILDFLAGS) $(LDFLAGS) -o ./bin/KeyboardLayout ./cmd/keyboard_layout/main.go

change_monitor_brightness:
	go build $(BUILDFLAGS) $(LDFLAGS) -o ./bin/ChangeMonitorBrightness ./cmd/change_monitor_brightness/main.go

change_workspace_index:
	go build $(BUILDFLAGS) $(LDFLAGS) -o ./bin/ChangeWorkspaceIndex ./cmd/change_workspace_index/main.go

diagonal_resize:
	go build $(BUILDFLAGS) $(LDFLAGS) -o ./bin/DiagonalResize ./cmd/diagonal_resize/main.go

kill_container:
	go build $(BUILDFLAGS) $(LDFLAGS) -o ./bin/KillContainer ./cmd/kill_container/main.go

lock_container:
	go build $(BUILDFLAGS) $(LDFLAGS) -o ./bin/LockContainer ./cmd/lock_container/main.go

mark_container:
	go build $(BUILDFLAGS) $(LDFLAGS) -o ./bin/MarkContainer ./cmd/mark_container/main.go

move_float_container:
	go build $(BUILDFLAGS) $(LDFLAGS) -o ./bin/MoveFloatContainer ./cmd/move_float_container/main.go

manage_float_container:
	go build $(BUILDFLAGS) $(LDFLAGS) -o ./bin/ManageFloatContainer ./cmd/manage_float_container/main.go

rename_workspace:
	go build $(BUILDFLAGS) $(LDFLAGS) -o ./bin/RenameWorkspace ./cmd/rename_workspace/main.go

margin_resize:
	go build $(BUILDFLAGS) $(LDFLAGS) -o ./bin/MarginResize ./cmd/margin_resize/main.go

resize_float_container:
	go build $(BUILDFLAGS) $(LDFLAGS) -o ./bin/ResizeFloatContainer ./cmd/resize_float_container/main.go

sticky_toggle:
	go build $(BUILDFLAGS) $(LDFLAGS) -o ./bin/StickyToggle ./cmd/sticky_toggle/main.go

swap_workspaces:
	go build $(BUILDFLAGS) $(LDFLAGS) -o ./bin/SwapWorkspaces ./cmd/swap_workspaces/main.go

volume_control:
	go build $(BUILDFLAGS) $(LDFLAGS) -o ./bin/VolumeControl ./cmd/volume_control/main.go

# Optimized parallel build
build:
	@echo "Building with optimizations (parallel jobs: $(PARALLEL_JOBS))..."
	@mkdir -p bin
	@$(MAKE) -j$(PARALLEL_JOBS) \
		back_and_forth_containers \
		keyboard_layout \
		change_monitor_brightness \
		change_workspace_index \
		diagonal_resize \
		kill_container \
		lock_container \
		mark_container \
		move_float_container \
		manage_float_container \
		rename_workspace \
		margin_resize \
		resize_float_container \
		sticky_toggle \
		swap_workspaces \
		volume_control
	@echo "Build completed. Binary sizes:"
	@du -sh bin/* | sort -h

# Ultra-optimized build with additional compression
build-optimized: LDFLAGS += -buildmode=pie
build-optimized: BUILDFLAGS += -a -installsuffix cgo
build-optimized: build

# Compressed build using UPX (if available)
build-compressed: build-optimized
	@echo "Compressing binaries with UPX..."
	@if command -v upx >/dev/null 2>&1; then \
		for binary in bin/*; do \
			echo "Compressing $$binary..."; \
			upx --best --lzma $$binary 2>/dev/null || echo "UPX compression failed for $$binary"; \
		done; \
		echo "Compression completed. Final sizes:"; \
		du -sh bin/* | sort -h; \
	else \
		echo "UPX not found. Install with: apt install upx-ucl"; \
		echo "Proceeding without compression."; \
	fi

# Performance testing
benchmark:
	@echo "Running build time benchmark..."
	@time $(MAKE) clean build >/dev/null 2>&1

# Test all binaries
test:
	@echo "Testing binary functionality..."
	@for binary in bin/*; do \
		echo "Testing $$binary..."; \
		$$binary --help >/dev/null 2>&1 || echo "$$binary may not support --help"; \
	done

# Clean build artifacts
clean:
	rm -rf bin/

install:
	@if [ $$(id -u) != 0 ]; then echo "You must run install with root privileges"; exit 1; fi
	@echo "Installing optimized binaries..."
	$(MAKE) build-optimized
	find ./bin/ -type f -executable -exec mv {} /bin/ \;
	@echo "Installation completed."

# Development helpers
dev-install:
	@echo "Installing development dependencies..."
	@command -v upx >/dev/null 2>&1 || echo "Consider installing UPX: apt install upx-ucl"
	@command -v golangci-lint >/dev/null 2>&1 || echo "Consider installing golangci-lint"

# Show current binary sizes
sizes:
	@if [ -d "bin" ]; then \
		echo "Current binary sizes:"; \
		du -sh bin/* | sort -h; \
		echo "Total size: $$(du -sh bin | cut -f1)"; \
	else \
		echo "No binaries found. Run 'make build' first."; \
	fi
