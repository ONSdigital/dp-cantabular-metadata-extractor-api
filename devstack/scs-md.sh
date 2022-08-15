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
    "babbage,9addfc57ab7db67d1a2ac77b526bd6016a5eebb5|make build"
    "dp-api-router,479c4c05d9d506993915d6cee0af356f7ed525be|"
    "dp-cantabular-api-ext,7dafd3c8cd8832b5644008a8248d57e5ee2924dc|make setup"
    "dp-cantabular-dimension-api,2872fdc4234d953ec050be6dc4b595c0a16eb260|"
    "dp-cantabular-metadata-extractor-api,2809afb1228611da7cec107ae2e48483283b0270|"
    "dp-cantabular-metadata-service,d28e66ded9489a27f2e17c91ddb5897d9be7ff9c|make setup"
    "dp-cantabular-server,3a5a2e83762152b3e45df3dabb37ef5ad9bde244|make setup"
    "dp-compose,41f6429098632ee6738bb5b06ac69e89df0332c3|"
    "dp-dataset-api,fd48f0d07455363ce0a26273a8fa51f29d643c64|"
    "dp-download-service,286dfefa44ae48d584fe83555674dec0408b571e|"
    "dp-frontend-dataset-controller,9af0dc34764eb740eb9841abb5238aeddb0d3d1f|make generate-prod"
    "dp-frontend-router,e325614d0b2fa5269e41bd8dc67fe15eaa42c00c|make assets"
    "dp-import-api,42363ef883f3de178e12258c386bbcf248f73dad|"
    "dp-import-cantabular-dataset,12d261dad5a33c44f16a1a065dfedfa15c9235c9|"
    "dp-import-cantabular-dimension-options,8a8c4984e9126f1b3c142bbca3484f33f75747c0|"
    "dp-publishing-dataset-controller,027b4a070e538c5e73b39e1d858328a7d135828e|"
    "dp-recipe-api,bc72f556a36395b195b34cb94480d759b287c649|"
    "florence,b086ec4e0a078942daad612582c95611a13ba465|make build"
    "zebedee,b72fad73eeaeee792d22effc05fca874c4891ff6|make build"
    "The-Train,9147c1fd7c3158059ec8e5ef7354c8aaaaf582a4|make build"
)

# current directory
DIR="$PWD"
DP_COMPOSE_DIR="$DIR/dp-compose"
ACTION=$1

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
        git clone git@github.com:ONSdigital/"${repo}".git 2> /dev/null
        logSuccess "Cloned $repo"
    done
}

goodCloneServices() {
    cd "$DIR" || exit
    for service in "${SERVICES[@]}"; do
        getvalues "$service"
        git clone git@github.com:ONSdigital/"${repo}".git 2> /dev/null
        cd "$repo" || exit
        git reset --hard "$sha"
        cd ..
        logSuccess "Cloned $repo"
        echo
    done
    bumpOurs
}

bumpOurs() {
    cd "$DIR" || exit
    cd dp-cantabular-metadata-extractor-api
    git checkout feature/devstack-minimal-2021
    git pull
    cd ..
    # TODO florence etc.
}

rmServices() {
    cd "$DIR" || exit
    doChown
    chmod +w -R .
    for service in "${SERVICES[@]}"; do
        getvalues "$service"
        rm -rf "$repo"
    done
}


quickCloneServices() {
    cd "$DIR" || exit
    for service in "${SERVICES[@]}"; do
        getvalues "$service"
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
    DP_CANTABULAR_IMPORT_DIR="$DP_COMPOSE_DIR/cantabular-import"
    cd "$DP_CANTABULAR_IMPORT_DIR" || exit
    make init-db
    logSuccess "Importing Recipes & Dataset documents... Done."
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
