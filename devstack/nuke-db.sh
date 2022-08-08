#!/bin/ksh
#this has to run from the root of the ONS journey directory

# stop & rm
docker rm -f cantabular-metadata-pub_mongodb_1
# rebuild
./scs-md.sh up
sleep 10
# fresh mongo db
./scs-md.sh init-db
# delete collections
sudo rm -rf $zebedee_root/zebedee/collections/*
