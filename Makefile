run:
	./pokedexcli

build:
	echo "Building..."
	go build

fullbuild: updateassets
	echo "Building..."
	go build

test:
	echo "Testing..."
	go test ./...

updateassets:
	echo "Updating assets..."
	chmod +x scripts/setup_assets.sh
	./scripts/setup_assets.sh

removelock:
	echo "Deleting .git/index.lock..."
	rm .git/index.lock