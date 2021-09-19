build:
		go build -v

run: build
	./gmail-notifier start

# (build but with a smaller binary)
dist:
	go build -ldflags="-w -s" -gcflags=all=-l -v

# (even smaller binary)
pack: dist
	upx ./gmail-notifier

