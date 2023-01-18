# dp-cantabular-metadata-extractor-api
Supply Cantabular metadata for Florence metadata journey.

### Getting started

* Run `make debug`

### Dependencies

* Requires Cantabular Server running on port 8491 (see [dp-cantabular-server](https://github.com/ONSdigital/dp-cantabular-server))
* Requires Cantabular UI running on port 8080 (see [dp-cantabular-ui](https://github.com/ONSdigital/dp-cantabular-ui))
* Requires Cantabular Metadata running on port 8493 (see [dp-cantabular-metadata](https://github.com/ONSdigital/dp-cantabular-server))

There are the following further dependencies if running with AUTHORISATION_ENABLED=true (see [Running with Authentication Enabled](README.md#running-with-authentication-enabled)):-

* Requires the Permissions API, port forwarded from the relevant AWS environment, running on port 25400
* Requires the Identity API, port forwarded from the relevant AWS environment, running on port 25600

* No further dependencies other than those defined in `go.mod`

### Configuration

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

### Running with Authentication Enabled



### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2022, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

