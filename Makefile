run:
	./pokedexcli

build: setup
	echo "Building..."
	go build

test:
	echo "Testing..."
	go test ./...

setup:
	chmod +x scripts/setup_assets.sh
	./scripts/setup_assets.sh
