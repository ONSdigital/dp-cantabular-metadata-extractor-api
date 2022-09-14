#!/bin/bash
#
# Temporary Fix (TM) to patch the cantabular server to use old synth data
# XXX This should go away once we get the synth data with UR which matches metadata
# This is run in the dp-cantabular-server directory

# extract from encrypted tarball
make setup

echo $PWD

cp ../dp_synth_config_1.dat cantabular/data/input/

if ! [[ -f "cantabular/data/input/dp_synth_config_1.dat" ]]; then
    echo "need to copy dp_synth_config_1.dat under $PWD/cantabular/data/input"
    exit 1
fi

patch -p1 < ../dp-cantabular-metadata-extractor-api/devstack/synth.patch 
