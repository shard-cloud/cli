.PHONY: completions
completions:
	mkdir -p completions
	@go run ./cmd/shardcloud/main.go completion zsh >"completions/completions-zsh.zsh"
	@go run ./cmd/shardcloud/main.go completion fish >"completions/completions-fish.fish"
	@go run ./cmd/shardcloud/main.go completion bash >"completions/completions-bash.bash"
