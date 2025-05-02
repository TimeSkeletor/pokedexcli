#!/bin/bash

ASSETS_DIR="assets/pokemon-sprites"

if [ -d "$ASSETS_DIR" ]; then
    echo "Updating sprite repo..."
    git -C "$ASSETS_DIR" pull
else
    echo "Cloning sprite repo..."
    git clone https://github.com/PokeAPI/sprites.git "$ASSETS_DIR"
fi
