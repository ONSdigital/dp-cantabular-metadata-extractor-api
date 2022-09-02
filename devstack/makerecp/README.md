# makerecp

`makerecp` as its name suggests will create the JSON payload for the Recipe API

All below commands need running from this directory.

It can also check whether the expected data exists in the metadata server.

A list of "working" (in the sense of the metadata journey working to preview) 
ONS Dataset IDs is included and can be displayed by:

`go run . -list`

This list worked against `cantabm_v10-1-0__ar2776-c21ew_metadata-v1-0_cantab_20220812-1_20220817-1`
with the sed modifications in `dp-cantabular-metadata-service/Makefile` at 20220902.


## Sending a req to the Recipe API

The simplest use is to use the wrapper script `cant-recipe2.sh`.

Use an ID from the "working list" with IDs starting with TS ("Topic Summary")
to be preferred since at the time of writing that was the focus for data
quality.

eg. `$ ./cant-recipe2.sh TS002` will use curl to POST a request with the correct headers.

If you need to confirm using the correct UUID in a request like 

`curl "http://localhost:8081/recipes/258cd19a-58c5-4e12-b2c9-b18b6376a0a3"`

will confirm by querying the Recipe API (via Florence)

# Generating the JSON body

If you wish to use Postman or similar just use

`go run . -id TS002` to generate the JSON payload

On a Mac piping like

`go run . -id TS002|pbcopy` 

will put the payload into the Copy/Paste buffer (or use `xsel` or similar on
Linux)

You can set the correct (?) alias/name of the recipe rather than using the
hard-coded "Testing for metadata demo v3" via

`go run . -id TS002 -setalias`

(note this breaks the selenium test)

# Metadata Server Data validation

`go run . -id TS002 -check` 

will confirm the expected dataset id and metadata dimensions (non-geographic)
are in use by the metadata-server.

`go run . -checkdims "ltla,sex"`

will confirm the presence of those dimensions in the main cantabular server.

(This is useful for finding working recipes given the current limitations of
synthetic data)

# Using against sandbox

Set up port forwarding the ext-api and metadata servers and use the
`extapihost` and `host` flags.
