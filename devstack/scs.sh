#!/usr/bin/env bash 
# STM forked this from an existing copy

#set -e

##################### VARIABLES ##########################

# prompt colours
GREEN="\e[32m"
RESET="\e[0m"

# services
SERVICES="
The-Train,292cabc52328c2c93b1bdc45de45269bee2cda7f
babbage,9e8b02804ccf488690974c53e1a91704aeed550e
dp-api-router,4a775fb3aa62dd005996e471587625f29429fa08
dp-cantabular-api-ext,2c1737cdd19d9cebbfb2369aabca9e9e1f8e0837
dp-cantabular-csv-exporter,6a9f6f26db298adb675861f979319c59de396883
dp-cantabular-dimension-api,2319cd72549440ba70c54bdecb7b75f405b9805f
dp-cantabular-filter-flex-api,a90af05a27eb994345a446baff4d33daf7ce57ba
dp-cantabular-metadata-exporter,b6f93b52405be4b1fad7c81d6b3358e2c0c5a24b
dp-cantabular-metadata-service,1af507befd834908f6d0f5e68e62baff7bea8295
dp-cantabular-server,79d958db5f045811abc153732d10b1fcc9a8492c
dp-cantabular-xlsx-exporter,9053af977fa0bd00e3029dd54ccafc698cc5073c
dp-compose,8515e58f4bd11eed98302b4a0fab90a4a0cbda5c
dp-dataset-api,157c21b6f3de188bced6ff94617664bd5b666e7e
dp-download-service,27013e88007863f62a7bb17c415238cea9c45731
dp-filter-api,217cb8fee0633a4c452c02068557c7596344390e
dp-frontend-dataset-controller,4f0d53dcb0cf34c196cd4382526c0eb7948d1d92
dp-frontend-filter-flex-dataset,eb35f1697db4c0c878cb219b54c008325f223b06
dp-frontend-router,84c394b6239edb11f0439259dd5d2cf14871a9ec
dp-import-api,36c7cae1d6a5bba567ef211eaf408e450d4965d1
dp-import-cantabular-dataset,6427bb91f91654f135662c9bef87cb858e4aecc3
dp-import-cantabular-dimension-options,6336a0b76806cc505a90c4993b2f94922d3671fd
dp-publishing-dataset-controller,77579b6b51562bd4df573e5583dd405b9bc0c461
dp-recipe-api,9947cfca3b844cf58c7dd14347e4e7a49071cbf7
florence,a41a1c4d0e4ac36c28e4346c389ca6ee4d8b2867
zebedee,7ab2dbaface788a76ef1545feadf8accbca74b68
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
