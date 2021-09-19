build:
		go build -v

run: build
	./gmail-notifier start
