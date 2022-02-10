#!/bin/bash
set -euo pipefail

# Normalize an HTML document by removing or replacing unwanted characters and entities.
usage="Usage: normalize-html FILE"

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

# Code points:
# U+0022    QUOTATION MARK (")
# U+0027    APOSTROPHE (')
# U+2018    LEFT SINGLE QUOTATION MARK (‘)
# U+2019    RIGHT SINGLE QUOTATION MARK (’)
# U+201C    LEFT DOUBLE QUOTATION MARK (“)
# U+201D    RIGHT DOUBLE QUOTATION MARK (”)
# U+2026    HORIZONTAL ELLIPSIS (…)
# U+2014    EM DASH (—)

# Entities:
# &#8216;   LEFT SINGLE QUOTE
# &#8217;   RIGHT SINGLE QUOTE
# &#8220;   LEFT DOUBLE QUOTE
# &#8221;   RIGHT DOUBLE QUOTE
# &#8230;   HORIZONTAL ELLIPSIS
# &#8211;   EN DASH
# &quot;    DOUBLE QUOTE (")

sed \
    --regexp-extended \
    --in-place \
    "s/[\‘\’]/'/g; s/[\“\”]/\"/g; s/&#8216;|&#8217;/'/g; s/&#8220;|&#8221;|&quot;/\"/g; s/&#8230;/…/g; s/&#8211;/--/g" \
    "${1}"
