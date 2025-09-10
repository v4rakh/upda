#!/usr/bin/env bash

REPOS=(
    "promhippie/hcloud_exporter"
    "crazy-max/diun"
)

for REPO in "${REPOS[@]}"; do
    SAFE_REPO_NAME=$(echo "$REPO" | tr '/' '_')
    LAST_RELEASE_FILE="last_release_${SAFE_REPO_NAME}.txt"

    LATEST_TAG=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$LATEST_TAG" ]; then
        echo "[$REPO] Failed to fetch latest release tag."
        continue
    fi

    if [ -f "$LAST_RELEASE_FILE" ]; then
        LAST_TAG=$(cat "$LAST_RELEASE_FILE")
    else
        LAST_TAG=""
    fi

    if [ "$LATEST_TAG" != "$LAST_TAG" ]; then
        echo "[$REPO] New release detected: $LATEST_TAG"
        echo "$LATEST_TAG" >"$LAST_RELEASE_FILE"
        HOSTNAME=$(hostname)
        JSON_PAYLOAD=$(
            cat <<EOF
        {
        "application": "$REPO",
        "host": "$HOSTNAME",
        "version": "$LATEST_TAG"
        }
EOF
        )
        upda webhook send "$JSON_PAYLOAD"
    else
        echo "[$REPO] No new release. Latest is $LATEST_TAG"
    fi
done
