#!/bin/bash

# This bootstrap script sets up the container
#   to be a usable development environment.

# Add nodesource repo
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo bash -

# Install system dependencies
sudo apt-get install -y --no-install-recommends \
        curl \
        zsh \
        git \
        neovim \
        nano \
        golang-go \
        nodejs

# Set up nvim plugin manager
curl -fLo "${XDG_DATA_HOME:-$HOME/.local/share}"/nvim/site/autoload/plug.vim --create-dirs \
       https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
