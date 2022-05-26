#!/bin/bash 

NOCOLOR='\033[0m'
GREEN='\033[0;32m'
RED='\033[0;31m'

if [[ $EUID -ne 0 ]]; then 
  echo -e "${RED}[ERROR]${NOCOLOR} Please run this script in root"
else 
  mkdir /opt/wgManagerAPI 

  echo -e "${GREEN}[INFO: CREATE]${NOCOLOR} root configuration folder"

  cp src/config/template.json /opt/wgManagerAPI/config.json

  echo -e "${GREEN}[INFO: COPY]${NOCOLOR} template configuration file"

  mkdir /opt/wgManagerAPI/wg

  echo -e "${GREEN}[INFO: CREATE]${NOCOLOR} database folder"

  mkdir /opt/wgManagerAPI/logs 

  echo -e "${GREEN}[INFO: CREATE]${NOCOLOR} logs folder"
fi