#!/bin/bash
set -euo pipefail

# Process images and create thumbnails
# WARNING: Before running this script you should strip EXIF data from
# the images and make sure each photo is oriented correctly.

usage="Usage: ./process-images.sh DIRECTORY"
if [ "$#" -ne 1 ]
then
    echo "${usage}"
    exit 1
fi

imgDir="${1}"
if ! [ -d "${imgDir}" ]
then
    echo "Error: ${imgDir} is not a directory"
    echo "${usage}"
    exit 1
fi

if ! [ -x "$(command -v exiftool)" ]
then
    echo "Error: exiftool is not available"
    exit 1
fi

if ! [ -x "$(command -v gm)" ]
then
    echo "Error: GraphicsMagick is not available"
    exit 1
fi
if ! [ -x "$(command -v gm)" ]
then
    echo "Error: GraphicsMagick is not available"
    exit 1
fi

# Reduce file size of all images
# WARNING this will overwrite the original files
gm mogrify -quality 50 "${imgDir}/*.jpg"

# Create thumbnails
mkdir -p "${imgDir}/thumbs"
find "${imgDir}" -type f -name '*.jpg' -exec cp '{}' "${imgDir}/thumbs" \;
gm mogrify -thumbnail 400x -quality 100 "${imgDir}/thumbs/*.jpg"
