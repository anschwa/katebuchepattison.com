#!/bin/bash
set -euo pipefail

# Replace occurances of Unicode smartquotes with ASCII.
usage="Usage: smartquotes-to-ascii FILE"

if [ "$#" -ne 1 ]
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

# U+0022    QUOTATION MARK (")
# U+0027    APOSTROPHE (')
# U+2018    LEFT SINGLE QUOTATION MARK (‘)
# U+2019    RIGHT SINGLE QUOTATION MARK (’)
# U+201C    LEFT DOUBLE QUOTATION MARK (“)
# U+201D    RIGHT DOUBLE QUOTATION MARK (”)

# Replace smartquotes in-place, but create a backup of the original file
sed -i.bak "s/[\‘\’]/'/g; s/[\“\”]/\"/g" "${1}"
