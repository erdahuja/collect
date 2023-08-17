db:
	go run app/scratch/db/main.go

sheetsconsumer:
	go run app/services/sheetsconsumer/main.go

api:
	go run app/services/collect-api/main.go

status-api:
	curl -il localhost:3000/v1/status

status-debug:
	curl -sS localhost:4000/debug/liveness

tidy:
	go mod tidy
	go mod vendor

# ==============================================================================
# Running tests within the local computer
# go install honnef.co/go/tools/cmd/staticcheck@latest
# go install golang.org/x/vuln/cmd/govulncheck@latest

test:
	CGO_ENABLED=0 go test -count=1 ./...
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...
	govulncheck ./...
