GUI_BINARY_NAME = mail-notifier-gui
DAEMON_BINARY_NAME = mail-notifier-daemon

.PHONY: build
build: build-daemon build-gui

.PHONY: daemon
build-daemon:
	go build -o $(DAEMON_BINARY_NAME) -v ./daemon

.PHONY: gui
build-gui:
	go build -o $(GUI_BINARY_NAME) -v ./gui

.PHONY: run-daemon
run-daemon: build-daemon
	./$(DAEMON_BINARY_NAME)

.PHONY: tidy
tidy:
	go mod tidy
