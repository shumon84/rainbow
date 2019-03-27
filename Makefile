.PHONY: all
all:
	go build
	time ./rainbow-table 10000 3000

.PHONY: deploy
deploy:
	GOOS=linux GOARCH=amd64 go build
	scp rainbow-table home:~/.local/bin
	ssh home "time ~/.local/bin/rainbow-table 10000 3000"
