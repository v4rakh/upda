#!/usr/bin/env bash

set -e

# source and target file definition, absolute paths
ENV_FILE="/usr/share/nginx/app/.env"
RUNTIME_FILE="/usr/share/nginx/html/conf/runtime-config.js"

echo "[INFO] Using env file $ENV_FILE"
echo "[INFO] Using runtime file $RUNTIME_FILE"

# Recreate config file
rm -rf $RUNTIME_FILE
touch $RUNTIME_FILE

# Add assignment
echo "const runtime_config = Object.freeze({" >> $RUNTIME_FILE

# Read each line in env file
# Each line represents key=value pairs

while read -r line || [[ -n "$line" ]];
do
  # Split env variables by character `=`
  if printf '%s\n' "$line" | grep -q -e '='; then
    varName=$(printf '%s\n' "$line" | sed -e 's/=.*//')
    varValue=$(printf '%s\n' "$line" | sed -e 's/^[^=]*=//')
  fi

  # Read value of current variable if exists as Environment variable
  value=$(printf '%s\n' "${!varName}")
  # Otherwise use value from env file
  [[ -z $value ]] && value=${varValue}

  # Append configuration property to JS file
  echo "  $varName: \"$value\"," >> $RUNTIME_FILE
done < $ENV_FILE

echo "});" >> $RUNTIME_FILE

echo "Object.defineProperty(window, 'runtime_config', {" >> $RUNTIME_FILE
echo "value: runtime_config," >> $RUNTIME_FILE
echo "writable: false" >> $RUNTIME_FILE
echo "});" >> $RUNTIME_FILE

echo "[INFO] $RUNTIME_FILE is ready"
echo "[INFO] Starting web server..."
echo ""
exec "$@"
