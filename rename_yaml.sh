#!/bin/bash
# Rename yaml files: replace all Chinese punctuation with English equivalents, space -> _

DIR="/c/Users/29326/Desktop/prismxscan/core/plugins/exploits"
count=0

for file in "$DIR"/*.yaml; do
    base=$(basename "$file")
    newname="$base"

    # space -> _
    newname="${newname// /_}"
    # Chinese brackets
    newname="${newname//（/_(}"
    newname="${newname//）/)}"
    # Chinese colon
    newname="${newname//：/:}"
    # Chinese comma -> English comma
    newname="${newname//，/,}"
    # Chinese period -> English period
    newname="${newname//。/.}"
    # Chinese enumeration comma -> underscore
    newname="${newname//、/_}"
    # Chinese exclamation
    newname="${newname//！/!}"
    # Chinese question mark
    newname="${newname//？/?}"
    # Chinese semicolon
    newname="${newname//；/;}"
    # Chinese double quotes
    newname="${newname//\"/_}"
    newname="${newname//\"/_}"
    # Chinese single quotes
    newname="${newname//\'/_}"
    newname="${newname//\'/_}"
    # Chinese square brackets
    newname="${newname//【/_[}"
    newname="${newname//】/]}"
    # Chinese angle brackets
    newname="${newname//《/_<}"
    newname="${newname//》/>_}"
    # Chinese em dash
    newname="${newname//—/-}"
    # Chinese ellipsis
    newname="${newname//…/...}"
    # Chinese middle dot
    newname="${newname//·/.}"
    # Chinese angle quotes 〈〉
    newname="${newname//〈/<}"
    newname="${newname//〉/>}"
    # Chinese brackets 〔〕
    newname="${newname//〔/_[}"
    newname="${newname//〕/]}"
    # Chinese brackets 〖〗
    newname="${newname//〖/_[}"
    newname="${newname//〗/]}"
    # Remove consecutive double underscores
    while [[ "$newname" == *__* ]]; do
        newname="${newname//__/_}"
    done
    # Remove leading/trailing underscores before extension
    newname="${newname%.yaml}"
    newname="${newname##_}"
    newname="${newname%_}"
    newname="${newname}.yaml"

    if [ "$newname" != "$base" ]; then
        mv "$file" "$DIR/$newname"
        echo "Renamed: $base -> $newname"
        ((count++))
    fi
done

echo "Done. $count files renamed."
