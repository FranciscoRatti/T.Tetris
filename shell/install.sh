#!/bin/bash
DIR=$(pwd)/TTetris

# CONFIG
if [ ! -d ~/.config/ttetris ]; then
  mkdir ~/.config/ttetris
fi

cp $DIR/resources/config.json ~/.config/ttetris/

# ESTATICOS
if [ ! -d /usr/share/ttetris ]; then
  sudo mkdir /usr/share/ttetris
fi

sudo cp -r $DIR/resources/audio /usr/share/ttetris/
sudo cp $DIR/resources/icon.png /usr/share/ttetris/

sudo chmod -R 775 /usr/share/ttetris/audio

# BINARIO
sudo cp $DIR/bin/linux /usr/bin/tetris

# ENTRADA
sudo cp $DIR/resources/ttetris.desktop /usr/share/applications/
sudo update-desktop-database