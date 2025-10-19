#!/bin/bash
set -euo pipefail

# Helper
is_int() { [[ "$1" =~ ^[0-9]+$ ]]; }

# SSH identities
ssh-add -D >/dev/null 2>&1
ssh-add -k /Users/morpheous/.ssh/githubWinStitch >/dev/null 2>&1

# Initialize if needed
[ -d .git ] || git init

git config user.name  "0187773933"
git config user.email "collincerbus@student.olympic.edu"

# Ensure remote
if ! git remote | grep -qx "origin"; then
	git remote add origin git@github.com:0187773933/MastersClosetRemoteDB.git
fi

# Skip if nothing changed
if git diff --quiet && git diff --cached --quiet; then
	echo "Nothing to commit — working tree clean."
	exit 0
fi

# Get numeric commit number
LastCommit=$(git log -1 --pretty="%B" 2>/dev/null | xargs || echo "0")
if is_int "$LastCommit"; then
	NextCommitNumber=$((LastCommit + 1))
else
	echo "Resetting commit number to 1"
	NextCommitNumber=1
fi

# Stage and commit
git add .
if [ -n "${1:-}" ]; then
	CommitMsg="$1"
	Tag="v1.0.$1"
else
	CommitMsg="$NextCommitNumber"
	Tag="v1.0.$NextCommitNumber"
fi
git commit -m "$CommitMsg"

# Remove tag locally and remotely if exists
if git tag | grep -qx "$Tag"; then
	git tag -d "$Tag" >/dev/null 2>&1
fi
if git ls-remote --tags origin | grep -q "refs/tags/$Tag$"; then
	git push --delete origin "$Tag" >/dev/null 2>&1 || true
fi

git tag "$Tag"

# Push safely
git push origin master
git push origin --tags
