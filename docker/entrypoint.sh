#!/bin/sh
set -e

# Check the first argument passed to the script
if [ "$1" = "cli" ]; then
    # Shift to remove "cli" from arguments and pass the remaining to vertex
    shift
    exec /usr/local/bin/vertex "$@"
else
    # Run the vertex-server by default
    exec /usr/local/bin/vertex-server "$@"
fi
