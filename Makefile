
PHONY release debug

release:
	go build -ldflags="-H windowsgui"
debug:
	go build