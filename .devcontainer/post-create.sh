#!/bin/sh

# Running inside docker
if ! [ -f "/.dockerenv" ]
then
    echo "Error: This script must be run inside a docker container"
    exit 1
fi

source .env

# Git Configuration
git config --global --add safe.directory ${WORKSPACE_FOLDER}

git config --global user.name "${GITHUB_USER}"
git config --global user.email "${GITHUB_EMAIL}"

. ${NVM_DIR}/nvm.sh
nvm install --lts
npm install -g @commitlint/cli @commitlint/config-conventional
pre-commit install
