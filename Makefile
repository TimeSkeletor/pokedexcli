run: setup
	echo "Running..."
	go build

	./pokedexcli

test:
	echo "Testing..."
	go test ./...

setup:
	chmod +x scripts/setup_assets.sh
	./scripts/setup_assets.sh
