#!/bin/bash

HASH_FILE="/home/morphs/git.hash"
REMOTE_HASH=$(git ls-remote https://github.com/0187773933/MastersClosetRemoteDB.git HEAD | awk '{print $1}')

if [ -f "$HASH_FILE" ]; then
	STORED_HASH=$(sudo cat "$HASH_FILE")
else
	STORED_HASH=""
fi

if [ "$REMOTE_HASH" == "$STORED_HASH" ]; then
	echo "No New Updates Available"
	cd /home/morphs/MastersClosetRemoteDB
	LOG_LEVEL=debug exec /home/morphs/MastersClosetRemoteDB/server "$@"
else
	echo "New updates available. Updating and Rebuilding Go Module"
	echo "$REMOTE_HASH" | sudo tee "$HASH_FILE"
	cd /home/morphs
	sudo rm -rf /home/morphs/MastersClosetRemoteDB
	git clone "https://github.com/0187773933/MastersClosetRemoteDB.git"
	sudo chown -R morphs:morphs /home/morphs/MastersClosetRemoteDB
	cd /home/morphs/MastersClosetRemoteDB
	/usr/local/go/bin/go mod tidy
	GOOS=linux GOARCH=amd64 /usr/local/go/bin/go build -o /home/morphs/MastersClosetRemoteDB/server
	LOG_LEVEL=debug exec /home/morphs/MastersClosetRemoteDB/server "$@"
fi