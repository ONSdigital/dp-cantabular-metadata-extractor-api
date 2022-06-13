#!/usr/bin/env bash 
# STM forked this from an existing copy

#set -e

##################### VARIABLES ##########################

# prompt colours
GREEN="\e[32m"
RESET="\e[0m"

# services
SERVICES="
The-Train,9147c1fd7c3158059ec8e5ef7354c8aaaaf582a4
babbage,dc7fe09edc188d1589360ad3cf74e8ecbef5c069
dp-api-router,79c062a93180763f2036d086c7558a8b3e8a5183
dp-cantabular-api-ext,a564bfa6ecab9acd9d1b2dc7c0014e74b5ce5774
dp-cantabular-csv-exporter,57a14a73076d72ad26b06ea78680acdd2eec2c95
dp-cantabular-dimension-api,1432fc302a41089de5050f531dbaf2cf1228050e
dp-cantabular-filter-flex-api,f86d1eedd65f08ba8d89688b8579a03e2f69152c
dp-cantabular-metadata-exporter,036f8a9e261c329ebd1e66124e8e568bcca60b5a
dp-cantabular-metadata-service,da353edc16c6240c9d859064486724b01a1ce14c
dp-cantabular-server,abc98c5004bfb44b5499fd073a009dcd594af82f
dp-cantabular-xlsx-exporter,4b30806bb3062fa2367f538c966bb69a597c9491
dp-compose,58902f633997b7703e898a74a0ba72e59006d4b9
dp-dataset-api,91ce52a4b9fd72f813bfb03546b6cdbe3a710abd
dp-download-service,286dfefa44ae48d584fe83555674dec0408b571e
dp-filter-api,77b815b2433bd6578809be2bc9463989deb1a4fd
dp-frontend-dataset-controller,89f49220392c6d2a10aabf4de403509e2045832e
dp-frontend-filter-flex-dataset,d699162bd75b836c677349b9a7eafbde6b3cd6fd
dp-frontend-router,0386e60c36248c8dd034947326b202d2d84774f2
dp-import-api,13e04509e7da30fe989a17493fd488b06f8f361b
dp-import-cantabular-dataset,0cc80b45f88e13bcf9e5606375fa9f7c95cd58a3
dp-import-cantabular-dimension-options,0f87c20d6f13748f84e2779a701f8468ca9be23d
dp-publishing-dataset-controller,027b4a070e538c5e73b39e1d858328a7d135828e
dp-recipe-api,ec16d92b27c2c6d6b31cf8a8797e52f24a86e330
florence,f2ae973fc5d5f2e8c54f33320176ffe7a3d2c04f
zebedee,f4c7da4cb0c0abebb7a926942cbb308fe2142d3c
"

# current directory
DIR="$PWD"

# directories
DP_BABBAGE_DIR="$DIR/babbage"
DP_CANTABULAR_API_EXT_DIR="$DIR/dp-cantabular-api-ext"
DP_COMPOSE_DIR="$DIR/dp-compose"
DP_CANTABULAR_IMPORT_DIR="$DP_COMPOSE_DIR/cantabular-import"
DP_CANTABULAR_SERVER_DIR="$DIR/dp-cantabular-server"
DP_CANTABULAR_METADATA_SERVER_DIR="$DIR/dp-cantabular-metadata-service"
DP_FLORENCE_DIR="$DIR/florence"
DP_FRONTEND_DATASET_CONTROLLER_DIR="$DIR/dp-frontend-dataset-controller"
DP_FRONTEND_ROUTER_DIR="$DIR/dp-frontend-router"
DP_THE_TRAIN_DIR="$DIR/The-Train"
DP_ZEBEDEE_DIR="$DIR/zebedee"
#DP_RECIPE_API_IMPORT_RECIPES_DIR="$DIR/dp-recipe-api/import-recipes" # STM
#DP_DATASET_API_IMPORT_SCRIPT_DIR="$DIR/dp-dataset-api/import-script" # STM

ACTION=$1

##################### FUNCTIONS ##########################
logSuccess() {
    echo -e "$GREEN ${1} $RESET"
}

splash() {
    echo "Start Cantabular Services (SCS)"
    echo ""
    echo "usual workflow 'clone','setup'"
    echo ""
    echo "Simple script to run cantabular import service locally and all the dependencies"
    echo ""
    echo "This script should be executed from your 'ons' roo/workspace folder"
    echo "Additionally source this as an alias on your '.bashrc' or '.zshrc', to use it from anywhere on your environment"
    echo "  e.g., alias scs='$HOME/ons/scs.sh'"
    echo ""
    echo "Partial List of commands (see source for full): "
    echo "   chown         - fix ownership from root for 'rm' etc."
    echo "   clone         - git clone all the required GitHub repos"
    echo "   down          - stop running the containers via docker-compose"
    echo "   help          - splash screen with all these options"
    echo "   init-db       - preparing db services. Run this once"
    echo "   goodclone     - clone at good working version"
    echo "   pull          - git pull the latest from your remote repos"
    echo "   quickclone    - faster, shallow git clone"
    echo "   rm            - remove repos"
    echo "   rmdocker      - agressively remove docker instances etc."
    echo "   setup         - preparing services. Run this once, before 'up'"
    echo "   up            - run the containers via docker-compose"
}

cloneServices() {
    cd "$DIR" || exit
    for service in $SERVICES; do
        repo="${service%,*}"
        git clone git@github.com:ONSdigital/"${repo}".git 2> /dev/null
        logSuccess "Cloned $repo"
    done
}

goodCloneServices() {
    cd "$DIR" || exit
    for service in $SERVICES; do
        repo="${service%,*}"
        sha="${service#*,}"
        git clone git@github.com:ONSdigital/"${repo}".git 2> /dev/null
        cd "$repo" || exit
        git reset --hard "$sha"
        cd ..
        logSuccess "Cloned $repo"
    done
}

rmServices() {
    cd "$DIR" || exit
    doChown
    chmod +w -R .
    for service in $SERVICES; do
        repo="${service%,*}"
        rm -rf "$repo"
    done
}


quickCloneServices() {
    cd "$DIR" || exit
    for service in $SERVICES; do
        repo="${service%,*}"
        git clone --depth 1 git@github.com:ONSdigital/"${repo}".git 2> /dev/null
        logSuccess "Cloned $repo"
    done
}


pull() {
    for repo in */ ; do
        cd "$DIR/$repo" || exit
        git pull 2> /dev/null
        logSuccess "'$repo' updated"
    done
}


initDB() {
    echo "Importing Recipes & Dataset documents..."
    cd "$DP_CANTABULAR_IMPORT_DIR" || exit
    make init-db
    logSuccess "Importing Recipes & Dataset documents... Done."
}

setupServices () {

    echo "Clean..."
    cd "$DP_CANTABULAR_IMPORT_DIR" || exit
    make full-clean
    logSuccess "Clean... Done."

    echo "Make Assets for dp-frontend-router..."
    cd "$DP_FRONTEND_ROUTER_DIR" || exit
    make assets
    logSuccess "Make Assets for dp-frontend-router... Done."

    echo "Generate prod for dp-frontend-dataset-controller..."
    cd "$DP_FRONTEND_DATASET_CONTROLLER_DIR" || exit
    make generate-prod
    logSuccess "Generate prod for dp-frontend-dataset-controller... Done."

    echo "Build florence..."
    cd "$DP_FLORENCE_DIR" || exit
    make build
    logSuccess "Build florence...  Done."

    echo "Build zebedee..."
    cd "$DP_ZEBEDEE_DIR" || exit
    make build
    logSuccess "Build zebedee...  Done."

    echo "Build babbage..."
    cd "$DP_BABBAGE_DIR" || exit
    make build
    logSuccess "Build babbage...  Done."

    echo "Build the-train..."
    cd "$DP_THE_TRAIN_DIR" || exit
    make build
    logSuccess "Build the-train... Done."

    echo "Preparing dp-cantabular-server..."
    cd "$DP_CANTABULAR_SERVER_DIR" || exit
    make setup
    logSuccess "Preparing dp-cantabular-server... Done."

    echo "Preparing dp-cantabular-metadata-service..."
    cd "$DP_CANTABULAR_METADATA_SERVER_DIR" || exit
    make setup
    logSuccess "Preparing dp-cantabular-metadata-service... Done."

    echo "Preparing dp-cantabular-api-ext..."
    cd "$DP_CANTABULAR_API_EXT_DIR" || exit
    make setup
    logSuccess "Preparing dp-cantabular-api-ext... Done."
    
    upServices

    initDB
}

upServices () {
    echo "Starting dp cantabular import..."
    cd "$DP_CANTABULAR_IMPORT_DIR" || exit
    make start-detached
    echo "Starting dp cantabular import... Done."
    pollFlorence
    logSuccess "Florence is available at http://localhost:8081/florence"
    logSuccess "         if 1st time accessing it the credentials are: florence@magicroundabout.ons.gov.uk / Doug4l"
    logSuccess "You may need to clear cookies"
}


downServices () {
    echo "Stopping base services..."
    cd "$DP_COMPOSE_DIR" || exit
    docker-compose down
    logSuccess "Stopping base services... Done."

    echo "Stopping dp cantabular import..."
    cd "$DP_CANTABULAR_IMPORT_DIR" || exit
    make stop
    logSuccess "Stopping dp cantabular import... Done."
}

doChown () {
    sudo chown -R "$USER": "$DIR"
}

rmDocker() {
    # Stop all containers
    docker stop $(docker ps -a -q)
    # Delete all containers
    docker rm $(docker ps -a -q)
    # Delete all images
    docker rmi $(docker images -q)
    # clean up
    docker volume prune
    docker volume rm $(docker volume ls -qf dangling=true)
    yes | docker image prune -a
}

pollFlorence() {
    while true
    do
        if nc -z 127.0.0.1 8081
        then
            echo "florence nearly ready!"
            break
        fi
        echo "waiting for florence..."
        sleep 10
    done
}

#####################    MAIN    #########################

case $ACTION in 
"chown") doChown;;
"clone") cloneServices;;
"down") downServices;;
"help") splash;;
"init-db") initDB;;
"pull") pull;;
"quickclone") quickCloneServices;;
"goodclone") goodCloneServices;;
"rm") rmServices ;;
"rmdocker") rmDocker ;;
"setup") setupServices;;
"up") upServices;;
*) echo "$ACTION - invalid action"; splash;;
esac
