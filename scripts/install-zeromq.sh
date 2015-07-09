#!/usr/bin/env bash
set -ex  # exit on error, show commands
wget https://github.com/zeromq/zeromq4-1/archive/v4.1.2.tar.gz
tar xzf v4.1.2.tar.gz
cd zeromq4-1-4.1.2
./autogen.sh
./configure --prefix=/usr
make
sudo make install

