# dp-cantabular-metadata-extractor-api

Supply Cantabular metadata for Florence metadata journey.

## Getting started

* Run `make debug`

## Dependencies

* Requires Cantabular Server running on port 8491 (see [dp-cantabular-server](https://github.com/ONSdigital/dp-cantabular-server))
* Requires Cantabular UI running on port 8080 (see [dp-cantabular-ui](https://github.com/ONSdigital/dp-cantabular-ui))
* Requires Cantabular Metadata running on port 8493 (see [dp-cantabular-metadata](https://github.com/ONSdigital/dp-cantabular-server))

There are also the following further dependencies if running with AUTHORISATION_ENABLED=true (see [Running with Authorisation Enabled](README.md#running-with-authorisation-enabled)) :-

* Requires the Permissions API, port forwarded from the relevant AWS environment, running on port 25400
* Requires the Identity API, port forwarded from the relevant AWS environment, running on port 25600

* No further dependencies other than those defined in `go.mod`

## Configuration

| Environment variable                      | Default                  | Description
| ----------------------------------------- | ------------------------ | -----------
| BIND_ADDR                                 | localhost:28300          | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT                 | 5s                       | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL                      | 30s                      | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT              | 90s                      | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)
| CANTABULAR_METADATA_URL                   | http://localhost:8493    | Host and port for `dp-cantabular-metadata-service`
| AUTHORISATION_ENABLED                     | false                    | dp-authorisation V2 enabled
| JWT_VERIFICATION_PUBLIC_KEYS              | [view here](https://github.com/ONSdigital/dp-authorisation/blob/main/authorisation/config.go#L20)                | JWT verification public keys
| PERMISSIONS_API_URL                       | http://localhost:25400   | Permissions API URL
| PERMISSIONS_CACHE_UPDATE_INTERVAL         | 60s                      | Permisssions cache update interval
| PERMISSIONS_MAX_CACHE_TIME                | 300s                     | Permissions max cache time
| ZEBEDEE_URL                               | http://localhost:8082    | Zebedee URL
| IDENTITY_WEB_KEY_SET_URL                  | http://localhost:25600   | Identity web key set URL
| AUTHORISATION_IDENTITY_CLIENT_MAX_RETRIES | 2                        | Identity client max retries

## Running with Authorisation Enabled

* The following 3 prerequisites need to be run, either directly or in docker:

1. cantabular server
2. cantabular ui
3. cantabular metadata

* Log in to the relevant AWS environment NB. the authentication is done by AWS Cognito in either sandbox, staging, or prod, depending on which environment you log into.
* Use consul to get the IP address and port for both the permissions api and identity api services.
* Port forward to the permissions api like this:

```shell
dp ssh <env> <IP address> -p 25400:<IP address>:<port>
```

E.g. dp ssh sandbox 10.30.139.75 -p 25400:10.30.139.75:22958

* Port forward to the identity api like this:

```shell
dp ssh <env> <IP address> -p 25600:<IP address>:<port>
```

E.g. dp ssh sandbox 10.30.138.21 -p 25600:10.30.138.21:28100

* Make sure that you have a Florence account with admin permissions in the AWS Cognito user pool e.g. the one named sandbox-florence-users. NB. You can find this in the AWS Management Console. If not then ask a Florence admin to create an admin account for you (via Florence).
* Get a JWT token, from the identity api, by sending a POST request to the following endpoint:

```hp
http://localhost:25600/v1/tokens
```

NB. It will require a JSON request body as follows: {"email":"your email","password":"your password"}. The JWT token will be in the response Header and will look something like this (for example):

```shell
Bearer eyJraWQiOiJqeFlva3pnVER5UVVNb1VTM0c0ODNoa0VjY3hFSklKdCtHVjAraHVSRUpBPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiI5OTg4YTljNy04MzYzLTRmM2EtOWRkYy0xNzU3MjAzMzBlNDUiLCJjb2duaXRvOmdyb3VwcyI6WyJyb2xlLWFkbWluIiwicm9sZS1wdWJsaXNoZXIiXSwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMi5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTJfV1NEOUVjQXN3IiwiY2xpZW50X2lkIjoiNGV2bDkxZzR0czVpc211ZGhyY2JiNGRhb2MiLCJvcmlnaW5fanRpIjoiYjllYmNkNDktMjg0My00ZDBhLWFlZjctZDE4OGQyMDRmODljIiwiZXZlbnRfaWQiOiI0ZTJlMTUyNi1iMTczLTQxMzEtYjgzYi03OWE0MWM2MDA5OTEiLCJ0b2tlbl91c2UiOiJhY2Nlc3MiLCJzY29wZSI6ImF3cy5jb2duaXRvLnNpZ25pbi51c2VyLmFkbWluIiwiYXV0aF90aW1lIjoxNjc0MDM3OTU1LCJleHAiOjE2NzQwMzg4NTQsImlhdCI6MTY3NDAzNzk1NSwianRpIjoiMzU3MGM2NjYtMzU4Yi00OTdkLWExZWYtMWUwNTg4MzRmYzM0IiwidXNlcm5hbWUiOiJmZWNkYTg5NS0xNTJhLTQ5MWQtYjBmYi0yNTAwMzhlNjRlZGMifQ.MJ7cDY8B7LsdWelCWHw-eLpb-fEBr7NRWwftkj5fbizuIjzdn7shzsV8qetAZtsjRaURjIZwugd3f637zMj0WV76e3Sj3L7QUW-KNQjiKqbYr5RtoeJ91fUv8UaU-o7-74fmhl2Y_D22QQHVngVMKtj74GowA1TA0AoCfR2qml6B5zrgtsizJth1ySPyZHorVkyo-qA4JT_ZJg4x7QbEDJYW43zKD5JASwZcP6KmXl19YfZEvPTf4y7taYiomNU3ro73hKzwuO61wz9KJVqo9JIhSnJ7Lb-Fc86C7BwVYBPHsucxTE5pXaOHY-zKpySe_PV1u6gBnZSPP7-CjuCivQ
```

* export AUTHORISATION_ENABLED=true
* Run the cantabular metatdata exporter api on port 28300

```shell
make debug
```

* Send a GET request to the following endpoint:

```hp
http://localhost:28300/cantabular-metadata/dataset/<name of cantabular blob>/lang/<language> 
```

NB. It will require the JWT Token, including the word 'Bearer', in a Header named Authorization. For example:

```shell
curl http://localhost:28300/cantabular-metadata/dataset/RM001/lang/en -H 'Authorization: Bearer eyJraWQiOiJqeFlva3pnVER5UVVNb1VTM0c0ODNoa0VjY3hFSklKdCtHVjAraHVSRUpBPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiI5OTg4YTljNy04MzYzLTRmM2EtOWRkYy0xNzU3MjAzMzBlNDUiLCJjb2duaXRvOmdyb3VwcyI6WyJyb2xlLWFkbWluIiwicm9sZS1wdWJsaXNoZXIiXSwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMi5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTJfV1NEOUVjQXN3IiwiY2xpZW50X2lkIjoiNGV2bDkxZzR0czVpc211ZGhyY2JiNGRhb2MiLCJvcmlnaW5fanRpIjoiYjllYmNkNDktMjg0My00ZDBhLWFlZjctZDE4OGQyMDRmODljIiwiZXZlbnRfaWQiOiI0ZTJlMTUyNi1iMTczLTQxMzEtYjgzYi03OWE0MWM2MDA5OTEiLCJ0b2tlbl91c2UiOiJhY2Nlc3MiLCJzY29wZSI6ImF3cy5jb2duaXRvLnNpZ25pbi51c2VyLmFkbWluIiwiYXV0aF90aW1lIjoxNjc0MDM3OTU1LCJleHAiOjE2NzQwMzg4NTQsImlhdCI6MTY3NDAzNzk1NSwianRpIjoiMzU3MGM2NjYtMzU4Yi00OTdkLWExZWYtMWUwNTg4MzRmYzM0IiwidXNlcm5hbWUiOiJmZWNkYTg5NS0xNTJhLTQ5MWQtYjBmYi0yNTAwMzhlNjRlZGMifQ.MJ7cDY8B7LsdWelCWHw-eLpb-fEBr7NRWwftkj5fbizuIjzdn7shzsV8qetAZtsjRaURjIZwugd3f637zMj0WV76e3Sj3L7QUW-KNQjiKqbYr5RtoeJ91fUv8UaU-o7-74fmhl2Y_D22QQHVngVMKtj74GowA1TA0AoCfR2qml6B5zrgtsizJth1ySPyZHorVkyo-qA4JT_ZJg4x7QbEDJYW43zKD5JASwZcP6KmXl19YfZEvPTf4y7taYiomNU3ro73hKzwuO61wz9KJVqo9JIhSnJ7Lb-Fc86C7BwVYBPHsucxTE5pXaOHY-zKpySe_PV1u6gBnZSPP7-CjuCivQ'
```

* Metadata should be returned successfully in JSON format.

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

Copyright © 2022, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
