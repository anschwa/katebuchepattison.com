#!/bin/bash

set -euo pipefail

# Strip all EXIF data from images in a directory
# WARNING this will overwrite the original files

exiftool -overwrite_original -all:all= -r "$1"
