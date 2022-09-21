# Misc Dev Support Scripts

## General scripts for Cantabular Import Journey

Create a directory named something like `ons` as a ONS root directory and
checkout the correct version (usually develop) of this repo under it.

Type `make` which will create symlinks for the shell scripts under this
ONS root directory

XXX You will also need a copy of `dp_synth_config_1.dat` under this directory.
This file is used to create the cantabular server database (not the cantabular
metadata server).

A description of these scripts follows.

* scs-md.sh

This is a wrapper script I forked after being given it by members of Team B which
now works to check out and run known versions of repos as used by
`dp-compose/cantabular-metadata-pub`

The versions are specified via lines like
`dp-api-router,4a775fb3aa62dd005996e471587625f29429fa08|` 

These *old* versions have been tested to "work" (for certain values of!) in
combination.  If you are working on a particular repo you may need to manually
bump that version via git.

The script will check out directories underneath itself, you probably want to copy it
to the root of the directory hierarchy you want.

This works for the use case for "edit_metadata" journey and I can't guarantee it
works for all cantabular journeys, but it should help.

Usual workflow to provision is

```
$ scs-md.sh rmdocker # if you have old docker containers clean start
$ scs-md.sh goodclone
# after this point you might want to checkout branches in particular repos you are working on
$ scs-md.sh setup
```

Beware services can take a long time to start.  The whole system can take many
mins to fully work.  Use florence to confirm the stack works as expected. 

It's not likely to work first time and debugging (see below) is often needed.

Hopefully recent versions of `scs-md.sh` should provision a more minimal and
stable stack.

Restarting everything can be done via

```
$ scs-md.sh down
$ scs-md.sh up
```

There is a nuclear option `scs-md.sh rmdocker` to aggressively remove docker instances
etc. if you need to restart from scratch.  This will destroy all traces of docker data
not just this stack.

* health.sh

This calls `/health` endpoints of the services we need to run and is useful
for debugging.  I've seen CRITICAL errors for some things (e.g. S3 uploaders)
which haven't broken my use case.  Maybe it's best to save a copy of the output
once things work to help with debugging future breakage

* nuke-db.sh

This resets the mongo database to a fresh state

* ver.sh

This is used to help maintain "pinned" versions in `scs-md.sh` and can be ignored
in normal use.

## General Debugging Suggestions

On macOS, you may need to increase the amount of memory committed to docker to 8G.
On Linux, I found adding extra swap helped.

Note on a low resource system it might be necessarily to run all this after a
fresh reboot before you started browsers, slack and your IDE.

It's not unusual for the whole stack to take quite a while to debug.  Look at
the output of `health.sh` above and try to identify which services are broken
(not running, not responding etc.).  Although this script can only contact
ports exposed on localhost and not all the docker internal ones. Try restarting
those docker containers and looking at container logs, e.g.

```
$ docker logs -f cantabular-metadata-pub_dp-cantabular-metadata-extractor-api_1
```

Note `health.sh` can't check the health of things exclusively accessible from
docker.

You may also need to `docker rmi` a particular image and rebuild it from scratch
(`mvn clean` type targets might be needed for Java builds etc.)  Sometimes
running `make debug` to run outside docker gives insight into issues.

Once working due to the number of services the whole stack isn't that stable
even in 16G of memory.  I would avoid updating it too frequently unless it's
really needed for your work.

You may have to fix up permissions with 

```
$ scs-md.sh chown
```

If you want to remove directories.

Good luck!

## Specific scripts for "edit_metadata" journey

* cant-recipe.sh

This is used to populate the recipe at the start of the journey, usually after
`nuke-db.sh`.

You will need to set your florence web password like

```
export FLORENCE_WEB_PW=XXXXX
```

* edit_meta_jor.side

This is a Selenium IDE (browser plugin) script to run the "edit_metadata"
journey on florence.  Looking at it may be useful for related journeys.

PUBLISHING 

Published dataset landing page is at eg.
http://localhost:20000/datasets/TS009/editions/2021/versions/1

The four files are visible at http://localhost:9002/buckets/public-bucket/browse

