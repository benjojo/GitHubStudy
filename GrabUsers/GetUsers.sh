#!/bin/bash

curl http://data.githubarchive.org/2012-{01..12}-{01..31}-{0..23}.json.gz | gzip -d | php ParseJSON.php >> unsortedusers.txt
cat unsortedusers.txt | sort | uniq > usernames.txt
rm unsortedusers.txt