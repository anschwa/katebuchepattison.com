#!/bin/bash
set -euo pipefail

# Download images from HTML templates made from WordPress pages
usage="Usage: wget-img TMPL OUTDIR"

if [ "$#" -ne 2 ]
then
    echo "${usage}"
    exit 1
fi

if ! [ -f "${1}" ]
then
    echo "Error: file does not exist: ${1}"
    echo "${usage}"
    exit 1
fi

if ! [ -d "${2}" ]
then
    echo "Error: Output directory does not exist: ${2}"
    echo "${usage}"
    exit 1
fi

wget --force-html --input-file="${1}" --directory-prefix="${2}"
