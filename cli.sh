#!/bin/bash -e

sudo rm -rf /usr/share/nginx/html/*.html

tar -cvf  htmls/scrm.tar htmls/*.html
sudo cp -r htmls/* /usr/share/nginx/html
