#!/bin/bash

URL="https://www.mobafire.com/images/champion/square"

xargs -a ./champs.txt -d '\n' -I {} \
	wget $URL/{}.png