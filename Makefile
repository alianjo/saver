build:
	go build -o kubectl-save

install: build
	sudo cp kubectl-save /usr/local/bin