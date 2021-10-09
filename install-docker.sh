#!/usr/bin/env bash

if [ -z $server_name ]; then
  read -p "please enter server_name(default:smzdm_reptile):" smzdm_reptile
fi
if [ -z $server_name ]; then
  server_name="smzdm_reptile"
fi

while :; do
  if [ ! -z $server_center_address ]; then
    break
  fi
  read -p "please enter server_center_address(required):" server_center_address
done

while :; do
  if [ ! -z $server_center_secret ]; then
    break
  fi
  read -p "please enter server_center_secret(required):" server_center_secret
done

echo
echo "server_name: $server_name"
echo "server_center_address: $server_center_address"
echo "server_center_secret: $server_center_secret"
echo 'input any key go on, or control+c over'
read

echo 'create volume'
docker volume create log
echo 'stop container'
docker stop $server_name
echo 'remove container'
docker rm $server_name
echo 'remove image'
docker rmi $server_name
echo 'docker build'
docker build -t $server_name .
echo 'docker run'
docker run -d \
  --restart=always \
  --name $server_name \
  -v log:/log \
  -e server_name=$server_name \
  -e server_center_address=$server_center_address \
  -e server_center_secret=$server_center_secret \
  $server_name

echo 'all finish'
