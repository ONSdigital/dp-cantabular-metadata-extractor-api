# Misc Dev Support Scripts

## General scripts for Cantabular Import Journey

* scs.sh

This is a wrapper script I was given (and forked) which now works to checkout
and run known versions of repos as used by `dp-compose/cantabular-import`

* health.sh

This calls /health endpoints of the services we need to run and is useful
for debugging.

* ver.sh

This is used to help maintain "pinned" versions in `scs.sh`

## Specific scripts for "edit_metadata" journey

* cant-recipe.sh

This is used to populate the recipe at the start of the journey.

* edit_meta_jor.side

This is a Selenium IDE (browser plugin) script to run the "edit_metadata"
journey.

