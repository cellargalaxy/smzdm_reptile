#!/usr/bin/env bash

if [ -z $retry ];then
    read -p "please enter retry(default:2):" retry
fi
if [ -z $retry ];then
    retry="2"
fi

if [ -z $maxPage ];then
    read -p "please enter max page(default:10):" maxPage
fi
if [ -z $maxPage ];then
    maxPage="10"
fi

if [ -z $timeout ];then
    read -p "please enter timeout(default:5[s]):" timeout
fi
if [ -z $timeout ];then
    timeout="5"
fi

if [ -z $sleep ];then
    read -p "please enter sleep(default:2[s]):" sleep
fi
if [ -z $sleep ];then
    sleep="2"
fi

if [ -z $listenPort ];then
    read -p "please enter listen port(default:8088):" listenPort
fi
if [ -z $listenPort ];then
    listenPort="8088"
fi

while :
do
    if [ ! -z $wxPushAddress ];then
        break
    fi
    read -p "please enter wx push address(required):" wxPushAddress
done

while :
do
    if [ ! -z $wxToken ];then
        break
    fi
    read -p "please enter wx token(required):" wxToken
done

echo 'retry:'$retry
echo 'maxPage:'$maxPage
echo 'timeout:'$timeout
echo 'sleep:'$sleep
echo 'listenPort:'$listenPort
echo 'wxPushAddress:'$wxPushAddress
echo 'wxToken:'$wxToken
echo 'input any key go on,or control+c over'
read

echo 'stop container'
docker stop smzdm_reptile
echo 'remove container'
docker rm smzdm_reptile
echo 'remove image'
docker rmi smzdm_reptile
echo 'docker build'
docker build -t smzdm_reptile .
echo 'docker run'
docker run -d \
--restart=always \
--name smzdm_reptile \
-p $listenPort:8088 \
-e RETRY=$retry \
-e MAX_PAGE=$maxPage \
-e TIMEOUT=$timeout \
-e SLEEP=$sleep \
-e WX_PUSH_ADDRESS=$wxPushAddress \
-e WX_TOKEN=$wxToken \
-v smzdm_reptile:/resources \
smzdm_reptile

echo 'all finish'