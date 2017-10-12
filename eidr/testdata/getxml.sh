#!/bin/bash

DOIS=`cat dois.txt`

for c in $DOIS; do
	filename=`echo -n $c | sed -e s@^[0-9\.]*/@@`
	curl -o $filename -L https://doi.org/$c?locatt=type:Full 
done
