#!/bin/bash

# Default value of n (number of last lines to display)
DEFAULT_N=100

# Check if the number of last lines is provided
if [[ $# -gt 0 ]]; then
  # Use the provided value of n
  n=$1
else
  # Use the default value of n
  n=$DEFAULT_N
fi

# Get the last n lines from the file within the Docker container
docker exec api_container tail -f -n "$n" ./archive_api.log