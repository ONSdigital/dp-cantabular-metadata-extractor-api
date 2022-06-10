#!/bin/ksh
# save working versions for scs.sh
for repo in */ ; do cd $repo; echo -n $repo; git log|head -1; cd .. ; done | sort | sed -e 's/\/commit /,/g'
