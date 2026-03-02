#!/usr/bin/env bash
# Checks public releases from GitHub
# Requires: curl, jq

set -e;
REPOS=(
    "promhippie/hcloud_exporter"
    "crazy-max/diun"
)

for REPO in "${REPOS[@]}"; do
    SAFE_REPO_NAME=$(echo "$REPO" | tr '/' '_')
    LAST_RELEASE_FILE="last_release_${SAFE_REPO_NAME}.txt"
    JSON=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest")

    if echo "$JSON" | jq -e '.message? // empty' >/dev/null; then
        echo "[$REPO] Failed to fetch latest release: $(echo "$JSON" | jq -r '.message')"
        continue
    fi

    LATEST_TAG=$(echo "$JSON" | jq -r '.tag_name // empty')

    if [ -z "$LATEST_TAG" ] || [ "$LATEST_TAG" = "null" ]; then
        echo "[$REPO] No tag_name in latest release."
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
        # use environment for URL, webhook ID and webhook's token
        upda webhook send --application "$REPO" --application-version "$LATEST_TAG" --host "$HOSTNAME" --metadata "hub_link=https://github.com/$REPO/releases"
    else
        echo "[$REPO] No new release. Latest is $LATEST_TAG"
    fi
done
