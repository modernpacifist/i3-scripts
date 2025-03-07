fmt:
	find . -type f -name "*.go" | xargs -n 1 go fmt

back_and_forth_containers:
	go build -o ./bin/BackAndForthContainers ./cmd/back_and_forth_containers/main.go

keyboard_layout:
	go build -o ./bin/KeyboardLayout ./cmd/keyboard_layout/main.go

change_monitor_brightness:
	go build -o ./bin/ChangeMonitorBrightness ./cmd/change_monitor_brightness/main.go

change_workspace_index:
	go build -o ./bin/ChangeWorkspaceIndex ./cmd/change_workspace_index/main.go

diagonal_resize:
	go build -o ./bin/DiagonalResize ./cmd/diagonal_resize/main.go

kill_container:
	go build -o ./bin/KillContainer ./cmd/kill_container/main.go

lock_container:
	go build -o ./bin/LockContainer ./cmd/lock_container/main.go

mark_container:
	go build -o ./bin/MarkContainer ./cmd/mark_container/main.go

move_float_container:
	go build -o ./bin/MoveFloatContainer ./cmd/move_float_container/main.go

manage_float_container:
	go build -o ./bin/ManageFloatContainer ./cmd/manage_float_container/main.go

rename_workspace:
	go build -o ./bin/RenameWorkspace ./cmd/rename_workspace/main.go

scale_float_window:
	go build -o ./bin/ScaleFloatWindow ./cmd/scale_float_window/main.go

sticky_toggle:
	go build -o ./bin/StickyToggle ./cmd/sticky_toggle/main.go

swap_workspaces:
	go build -o ./bin/SwapWorkspaces ./cmd/swap_workspaces/main.go

volume_control:
	go build -o ./bin/VolumeControl ./cmd/volume_control/main.go

build:
	make back_and_forth_containers
	make keyboard_layout
	make change_monitor_brightness
	make change_workspace_index
	make diagonal_resize
	make kill_container
	make lock_container
	make mark_container
	make move_float_container
	make manage_float_container
	make rename_workspace
	make scale_float_window
	make sticky_toggle
	make swap_workspaces
	make volume_control

install:
	@if [ $$(id -u) != 0 ]; then echo "You must run install with root privileges"; exit 1; fi

	make build

	find ./bin/ -type f -executable -exec mv {} /bin/ \;
