#!/bin/bash

DB=$1
BLIDS=$2

cmdArray=(
  "DELETE FROM blog_authors WHERE blog_id in ($BLIDS);"
  "DELETE FROM blog_media WHERE blog_id in ($BLIDS);"
  "DELETE FROM blog_tags WHERE blog_id in ($BLIDS);"
  "DELETE FROM blogs WHERE id in ($BLIDS);" 
)


for ((i = 0; i < ${#cmdArray[@]}; )); do
    cmd=""
    while [[ ${cmdArray[i]} != *";" ]]; do
        cmd+="${cmdArray[i]} "
        ((i++))
    done
    cmd+="${cmdArray[i]}"
    ((i++))
    echo "$cmd" | sqlite3 $DB
done

