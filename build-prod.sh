#!/bin/bash

# Check if the argument is provided
if [[ $# -eq 0 ]]; then
  # No argument provided, checkout to the develop branch
  git checkout main
else
  # Parse the argument
  case "$1" in
    -v|--version)
      # Check if the branch name is provided
      if [[ -z $2 ]]; then
        echo "Branch name not provided."
        exit 1
      fi
      # Checkout to the specified branch
      git checkout "release-$2"
      ;;
    *)
      echo "Invalid option: $1"
      exit 1
      ;;
  esac
fi
git pull
docker compose  --file docker-compose.prod.yml up -d --build