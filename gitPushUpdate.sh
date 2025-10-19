#!/bin/bash
set -e

# Helper: test if string is integer
function is_int() { [[ "$1" =~ ^[0-9]+$ ]]; }

# Reset SSH identities
ssh-add -D >/dev/null 2>&1
ssh-add -k /Users/morpheous/.ssh/githubWinStitch >/dev/null 2>&1

# Ensure repo initialized
if [ ! -d .git ]; then
	git init
fi

# Git user config
git config user.name "0187773933"
git config user.email "collincerbus@student.olympic.edu"

# Ensure remote exists
if ! git remote | grep -q "^origin$"; then
	git remote add origin git@github.com:0187773933/MastersClosetRemoteDB.git
fi

# Check if working tree is clean
if git diff --quiet && git diff --cached --quiet; then
	echo "Nothing to commit — working tree clean."
	exit 0
fi

# Determine commit number
LastCommit=$(git log -1 --pretty="%B" 2>/dev/null | xargs || echo "0")
if is_int "$LastCommit"; then
	NextCommitNumber=$((LastCommit + 1))
else
	echo "Resetting commit number to 1"
	NextCommitNumber=1
fi

# Stage all changes
git add .

# Choose commit message and tag
if [ -n "$1" ]; then
	CommitMsg="$1"
	Tag="v1.0.$1"
else
	CommitMsg="$NextCommitNumber"
	Tag="v1.0.$NextCommitNumber"
fi

# Commit and tag
git commit -m "$CommitMsg"

# Delete local tag if it exists already (to avoid push rejection)
if git tag | grep -q "$Tag"; then
	git tag -d "$Tag" >/dev/null 2>&1
fi
git tag "$Tag"

# Push safely
git push origin master
git push --tags