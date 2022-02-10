#!/bin/bash
set -euo pipefail

# Convert WordPress HTML into HTML template

if ! [ -f "${1}" ]
then
    echo "Error: file does not exist: ${1}"
    exit 1
fi

# Cleanup weird WP characters issues
sed --regexp-extended --in-place 's:…\.:…:g; s:\.{2}:\.:g' "${1}"

# Add newlines around <p>...</p>
sed --regexp-extended --in-place 's:<p>:<p>\n:g; s:</p>:\n</p>\n:g' "$1"

# Append boilerplate
cat << EOF >> "${1}"
{{define "content"}}
<article class="blog-post">
  <h1>Example</h1>
  <time datetime="2015-08-24">August 24, 2015</time>

</article>
{{end}}

<figure>
  <a href="">
    <img
      src=""
      alt=""
    />
  </a>

   <figcaption></figcaption>
 </figure>
EOF
