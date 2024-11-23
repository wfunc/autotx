#!/bin/bash


is_screen_exists() {
    local nodename="$1"
    local screen_output=$(screen -ls | grep "$nodename")
    if [ -n "$screen_output" ]; then
        return 0
    else
        return 1
    fi
}

srvname=m

if is_screen_exists $srvname;then
	read -p "A screen session with name '$srvname' exists. Do you want to attach to it? (y/n): " choice
        if [ "$choice" == "y" ]; then
            screen -r "$srvname"
	else
            echo "Bye"
            exit 0
        fi
else
screen -dmSL $srvname -s /srv/start.sh
screen -r $srvname
fi