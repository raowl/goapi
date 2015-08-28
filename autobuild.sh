#!/bin/bash

cur=`pwd`
./goapi &
inotifywait -mqr --timefmt '%d/%m/%y %H:%M' --format '%T %w %f' \
   -e modify ./ | while read date time dir file; do
    ext="${file##*.}"
    if [[ "$ext" = "go" ]]; then
    	killall goapi
        echo "$file changed @ $time $date, rebuilding..."
        go build
        ./goapi &
    fi
done
