lint:
		golangci-lint run --config=.golangci.yml ./...

test:	lint
		go test -tags integration -cover -failfast ./...

binary:
		 go build -a -o deliveroo-cron

