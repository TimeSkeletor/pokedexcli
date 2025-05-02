run:
	./pokedexcli

build: setup
	echo "Building..."
	go build

test:
	echo "Testing..."
	go test ./...

setup:
	echo "Updating assets..."
	chmod +x scripts/setup_assets.sh
	./scripts/setup_assets.sh

removelock:
	echo "Deleting .git/index.lock..."
	rm .git/index.lock