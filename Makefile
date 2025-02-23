test:
	@go test -v

bench:
	@go test -bench=. -v -benchmem

cover:
	@mkdir -p _dist
	@go test -coverprofile=_dist/coverage.out -v
	@go tool cover -html=_dist/coverage.out -o _dist/coverage.html
