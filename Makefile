compile: webhook

webhook:
	GOOS=linux go build -o bin/webhook functions/webhook/*.go

