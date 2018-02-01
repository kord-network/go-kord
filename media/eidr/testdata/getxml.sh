#!/bin/bash

set -e

while read id; do
  echo "$(date): getting ${id}..."
  curl -fsSLo "$(basename "${id}")" "https://doi.org/${id}?locatt=type:Full"
done < dois.txt
