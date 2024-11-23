build:
	go build -o kubectl-save

install: build cleanup
	sudo cp kubectl-save /usr/local/bin