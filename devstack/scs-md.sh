#!/usr/bin/env bash
# STM forked this from an existing version by Paulo Monteiro

#set -x
#set -e

##################### VARIABLES ##########################

# prompt colours
GREEN="\e[32m"
RESET="\e[0m"

# services
# special dp-compose & dp-cantabular-server
SERVICES=(
    "The-Train,9147c1fd7c3158059ec8e5ef7354c8aaaaf582a4|make build"
    "babbage,c77bf4936a4c8872c674e974a3e9c08d1ad89cf4|make build"
    "dp-api-router,919eda64b28f017cfc147ab607e838209fb65cd9|"
    "dp-cantabular-api-ext,8dedae2b88275e69a7c639163e46b8687741d964|make setup"
    "dp-cantabular-csv-exporter,500a1e4a4646f503f7f055b85c837c6358b6ba65|"
    "dp-cantabular-dimension-api,2872fdc4234d953ec050be6dc4b595c0a16eb260|"
    "dp-cantabular-metadata-exporter,39b239804592ad7668c5f277ea19f83d0f88ecfb|"
    "dp-cantabular-metadata-extractor-api,e3ba4e10950473204316ef6fb0a2de92ed44011e|"
    "dp-cantabular-metadata-service,e96bc56b00e250fc0b973aec8c61bfba2492d1db|make setup"
    "dp-cantabular-server,4b7a6958b98d621a0c802ef5849b22c24a19c46f|../dp-cantabular-metadata-extractor-api/devstack/patchcantabular.sh"
    "dp-cantabular-xlsx-exporter,f3ecb0547cd522abb01d95a980a5d16bfbb5e043|"
    "dp-compose,29f9c6dda9ef4aeb26240c56d9de342eef8eff98|../dp-cantabular-metadata-extractor-api/devstack/fixminio.sh"
    "dp-data-tools,86891bad6ab850fb76f9c252c5924fce7142b977|"
    "dp-dataset-api,7439ea1ddc32f1d8a40b288caf56387cb4fcccfb|"
    "dp-download-service,286dfefa44ae48d584fe83555674dec0408b571e|"
    "dp-filter-api,8d8e086da62fe137a554f1ce42dffdb5177331f2|"
    "dp-frontend-dataset-controller,9c317a2a5e582611c6817c333ddd6f0e85fbb245|make generate-prod"
    "dp-frontend-router,e4fb20f610b968319bf807a30c07e84dcc27c8e4|make assets"
    "dp-import-api,af2e41b18dd7193fc536b05a1efc0750d5893d2f|"
    "dp-import-cantabular-dataset,309fe5086203c6a73b5a8a314abec18c712357c6|"
    "dp-import-cantabular-dimension-options,d43216f8bfd65d7755d92f5a1446055de38d0084|"
    "dp-publishing-dataset-controller,027b4a070e538c5e73b39e1d858328a7d135828e|"
    "dp-recipe-api,44c5571c7c9c1fd8e12f60b121a6700e84563d28|"
    "dp-topic-api,9f16fa3ed83c8c05a09f3dde39c847b3b38b71a6|"
    "florence,653bece98594b1f76d3e80238cbf7da4afa74516|make build"
    "zebedee,b72fad73eeaeee792d22effc05fca874c4891ff6|make build"
)
#    "dp-file-downloader,e512a28e32f7686a6afcbd6929f145a198aa55ca|"
# "dp-cantabular-server,4b7a6958b98d621a0c802ef5849b22c24a19c46f|make setup"
#    florence

# current directory
DIR="$PWD"
DP_COMPOSE_DIR="$DIR/dp-compose"
ACTION=$1

############# override .dat (temp?)

if ! [[ -f "dp_synth_config_1.dat" ]]; then
    echo "need to copy dp_synth_config_1.dat to $PWD"
    exit 1
fi

##################### FUNCTIONS ##########################

getvalues() {
    service=$1
    repo="${service%,*}"
    values="${service#*,}"
    sha="${values%|*}"
    cmd="${values#*|}"
}

logSuccess() {
    echo -e "$GREEN ${1} $RESET"
}

splash() {
    echo "Start Cantabular Services (SCS)"
    echo ""
    echo "usual workflow 'goodclone','setup'"
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
    for service in "${SERVICES[@]}"; do
        getvalues "$service"
        if [[ $repo == "dp-cantabular-metadata-extractor-api" ]]; then
            echo "skipping $repo"
            continue
        fi
        git clone git@github.com:ONSdigital/"${repo}".git 2> /dev/null
        logSuccess "Cloned $repo"
    done
}

goodCloneServices() {
    cd "$DIR" || exit
    for service in "${SERVICES[@]}"; do
        getvalues "$service"
        if [[ $repo == "dp-cantabular-metadata-extractor-api" ]]; then
            echo "skipping $repo"
            continue
        fi
        git clone git@github.com:ONSdigital/"${repo}".git 2> /dev/null
        cd "$repo" || exit
        git checkout -q "$sha"
        cd ..
        logSuccess "Cloned $repo"
        echo
    done
}

rmServices() {
    cd "$DIR" || exit
    doChown
    chmod +w -R .
    for service in "${SERVICES[@]}"; do
        getvalues "$service"
        if [[ $repo == "dp-cantabular-metadata-extractor-api" ]]; then
            echo "skipping $repo"
            continue
        fi
        rm -rf "$repo"
    done
}


quickCloneServices() {
    cd "$DIR" || exit
    for service in "${SERVICES[@]}"; do
        getvalues "$service"
        if [[ $repo == "dp-cantabular-metadata-extractor-api" ]]; then
            echo "skipping $repo"
            continue
        fi
        git clone --depth 1 git@github.com:ONSdigital/"${repo}".git 2> /dev/null
        logSuccess "Cloned $repo"
    done
}


pull() {
    for repo in */ ; do
        if [[ $repo == "dp-cantabular-metadata-extractor-api" ]]; then
            echo "skipping $repo"
            continue
        fi
        cd "$DIR/$repo" || exit
        git pull 2> /dev/null
        logSuccess "'$repo' updated"
    done
}


initDB() {
    echo "Importing Recipes & Dataset documents..."
    DP_CANTABULAR_IMPORT_DIR="$DP_COMPOSE_DIR/cantabular-import"
    cd "$DP_CANTABULAR_IMPORT_DIR" || exit
    make init-db
    logSuccess "Importing Recipes & Dataset documents... Done."
    echo "Importing Topics"
    cd "$DIR"
    cd dp-data-tools || exit
    mongo topics-tools/gen-topics-database/mongo-init-scripts/topics-init.js
    cd ..
    cd dp-topic-api/db-scripts/insert-census-topics || exit
    mongo localhost:27017/topics insert-census-topics.js
    logSuccess "Topics"
}

setupServices () {
    cd "$DIR" || exit
    for service in "${SERVICES[@]}"; do
        getvalues "$service"
        if [[ $cmd != "" ]]; then
            cd "$repo" || exit
            eval "$cmd"
            cd "$DIR" || exit
        fi
    done

    upServices

    initDB
}

upServices () {
    echo "Starting cantabular-metadata-pub..."
    cd "$DP_COMPOSE_DIR/cantabular-metadata-pub" || exit
    ./compose.sh up -d
    echo "Starting dp cantabular metadata pub... Done."
    pollFlorence
    logSuccess "Florence is available at http://localhost:8081/florence"
    logSuccess "         if 1st time accessing it the credentials are: florence@magicroundabout.ons.gov.uk / Doug4l"
    logSuccess "You may need to clear cookies"
}

downServices () {
    echo "Stopping cantabular-metadata-pub..."
    cd "$DP_COMPOSE_DIR/cantabular-metadata-pub" || exit
    ./compose.sh stop
    echo "Stopping cantabular-metadata-pub... Done"
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
