GUI_BINARY_NAME = mail-notifier-gui
CLI_BINARY_NAME = mail-notifier-cli
DAEMON_BINARY_NAME = mail-notifier-daemon

.PHONY: build
build: build-daemon build-gui build-cli

.PHONY: build-daemon
build-daemon:
	go build -o $(DAEMON_BINARY_NAME) -v ./daemon

.PHONY: build-cli
build-cli:
	go build -o $(CLI_BINARY_NAME) -v ./cli

.PHONY: build-gui
build-gui:
	go build -o $(GUI_BINARY_NAME) -v ./gui

.PHONY: run-daemon
run-daemon: build-daemon
	./$(DAEMON_BINARY_NAME)

.PHONY: run-gui
run-gui: build-gui
	./$(GUI_BINARY_NAME)

.PHONY: run-cli
run-cli: build-cli
	./$(CLI_BINARY_NAME)

.PHONY: lint
lint:
	revive -formatter friendly -config revive.toml ./...

.PHONY: tidy
tidy:
	go mod tidy
