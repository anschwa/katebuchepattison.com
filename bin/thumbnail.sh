#!/bin/bash

# Create a thumbnail for every image in a directory

set -euo pipefail

mkdir -p "$1/thumbs"
mogrify -path "$1/thumbs" -thumbnail 400x -quality 50 "$1/*.jpg"
