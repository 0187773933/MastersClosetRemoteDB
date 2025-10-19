#!/bin/bash

# ===== CONFIGURATION =====
GITHUB_USER="0187773933"
REPO_NAME="MastersClosetRemoteDB"
REPO_URL="https://github.com/${GITHUB_USER}/${REPO_NAME}.git"
REPO_PATH="/home/morphs/${REPO_NAME}"
HASH_FILE="/home/morphs/git.hash"
GO_BIN="/usr/local/go/bin/go"
# ==========================

sudo chown -R morphs:morphs /home/morphs/SAVE_FILES

REMOTE_HASH=$(${GO_BIN} run github.com/git-lfs/git-lfs@latest ls-remote ${REPO_URL} HEAD 2>/dev/null | awk '{print $1}')
# fallback if git-lfs not installed
if [ -z "$REMOTE_HASH" ]; then
    REMOTE_HASH=$(git ls-remote ${REPO_URL} HEAD | awk '{print $1}')
fi

if [ -f "$HASH_FILE" ]; then
    STORED_HASH=$(sudo cat "$HASH_FILE")
else
    STORED_HASH=""
fi

if [ "$REMOTE_HASH" == "$STORED_HASH" ]; then
    echo "No New Updates Available"
    cd "$REPO_PATH"
    LOG_LEVEL=debug exec "$REPO_PATH/server" "$@"
else
    echo "New updates available. Updating and rebuilding ${REPO_NAME}"
    echo "$REMOTE_HASH" | sudo tee "$HASH_FILE"

    cd /home/morphs
    sudo rm -rf "$REPO_PATH"
    git clone "$REPO_URL"
    sudo chown -R morphs:morphs "$REPO_PATH"
    cd "$REPO_PATH"

    $GO_BIN mod tidy

    GOOS=$($GO_BIN env GOOS)
    GOARCH=$($GO_BIN env GOARCH)
    echo "Building for ${GOOS}/${GOARCH}"

    $GO_BIN build -o "$REPO_PATH/server"

    LOG_LEVEL=debug exec "$REPO_PATH/server" "$@"
fi