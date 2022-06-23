# Misc Dev Support Scripts

## General scripts for Cantabular Import Journey

Create an "ons" directory and copy the shell scripts in this directory into it
and run those copies.

A description of the scripts follows.

* scs.sh

This is a wrapper script I forked after being given it by members of Team B which
now works to check out and run known versions of repos as used by
`dp-compose/cantabular-import`

The versions are specified via lines like
`dp-api-router,4a775fb3aa62dd005996e471587625f29429fa08` 

These *old* versions have been tested to "work" (for certain values of!) in
combination.  If you are working on a particular repo you may need to manually
bump that version via git.

The script will check out directories underneath itself, you probably want to copy it
to the root of the directory hierarchy you want.

This works for the use case for "edit_metadata" journey and I can't guarantee it
works for all cantabular journeys, but it should help.

Usual workflow to provision is

```
$ scs.sh rmdocker # if you have old docker containers clean start
$ scs.sh goodclone
# after this point you might want to checkout branches in particular repos you are working on
$ scs.sh setup
```

Beware services can take a long time to start.  The whole system can take many
mins to fully work.  Use florence to confirm the stack works as expected. 

It's not likely to work first time and debugging (see below) is often needed.

Restarting everything can be done via

```
$ scs.sh down
$ scs.sh up
```

There is a nuclear option `scs.sh rmdocker` to aggressively remove docker instances
etc. if you need to restart from scratch.

* health.sh

This calls `/health` endpoints of the services we need to run and is useful
for debugging.  I've seen CRITICAL errors for some things (e.g. S3 uploaders)
which haven't broken my use case.  Maybe it's best to save a copy of the output
once things work to help with debugging future breakage

* nuke-db.sh

This resets the mongo database to a fresh state

* ver.sh

This is used to help maintain "pinned" versions in `scs.sh` and can be ignored
in normal use.

## General Debugging Suggestions

On macOS, you may need to increase the amount of memory committed to docker to 8G.
On Linux, I found adding extra swap helped.

It's not unusual for the whole stack to take quite a while to debug.  Look at
the output of `health.sh` above and try to identify which services are broken
(not running, not responding etc.).  Try restarting those docker containers and
looking at container logs, e.g.

```
$ cd dp-compose/cantabular-import
$ ./logs dp-cantabular-api-ext
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
$ scs.sh chown
```

If you want to remove directories.

Good luck!

## Specific scripts for "edit_metadata" journey

* cant-recipe.sh

This is used to populate the recipe at the start of the journey, usually after
`nuke-db.sh`.

* edit_meta_jor.side

This is a Selenium IDE (browser plugin) script to run the "edit_metadata"
journey on florence.  Looking at it may be useful for related journeys.
