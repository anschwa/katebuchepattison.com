#!/bin/bash

set -euo pipefail

# Process images for a blog post
# (1) Strip EXIF data from all images
# (2) Reduce image size
# (3) Generate thumbnails

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

if ! [ -x "$(command -v mogrify)" ]
then
    echo "Error: mogrify is not available"
    exit 1
fi

# Strip all EXIF data from images in a directory
# WARNING this will overwrite the original files
exiftool -overwrite_original -all:all= -r "${imgDir}"

# Reduce file size of all images
# WARNING this will overwrite the original files
mogrify -quality 75 "${imgDir}/*.jpg"

# Crete thumbnails
mkdir -p "${imgDir}/thumbs"
mogrify -path "${imgDir}/thumbs" -thumbnail 400x -quality 50 "${imgDir}/*.jpg"
