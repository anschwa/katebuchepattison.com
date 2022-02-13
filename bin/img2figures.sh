#!/bin/bash
set -euo pipefail

# Create HTML figures for each image found in IMG_DIR
usage="Usage: img2figures POST IMG_DIR"

if [ "$#" -ne 2 ]
then
    echo "${usage}"
    exit 1
fi

if ! [ -f "${1}" ]
then
    echo "Error: File does not exist: ${1}"
    echo "${usage}"
    exit 1
fi


if ! [ -d "${2}" ]
then
    echo "Error: Directory does not exist: ${2}"
    echo "${usage}"
    exit 1
fi

echo '<h2>Photos</h2>' >> "${1}"

for x in $(/bin/ls "${2}")
do
    f=$(basename "${x}")

    cat << EOF >> "${1}"
  <figure>
    <a href="img/${f}">
      <img
        src="img/thumbs/${f}"
        alt=""
      />
    </a>

    <figcaption></figcaption>
  </figure>

EOF
done
