#!/usr/bin/env bash
set -ex  # exit on error, show commands
wget https://github.com/jedisct1/libsodium/releases/download/1.0.3/libsodium-1.0.3.tar.gz
tar xzf libsodium-1.0.3.tar.gz
cd libsodium-1.0.3
./autogen.sh
./configure --prefix=/usr
make
sudo make install

