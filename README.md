# cast
small media cast platform for your humble abode

## Prerequisites

| tool                  | description                                           |
| --------------------- | ----------------------------------------------------- |
| golang                | to build the service app                              |
| flutter               | to build the UI ( web already built in ./ui/dist/web) |
| docker (with compose) | to run the services                                   |

## Build & run

Run these commands on the machine that is going to host the platform.

```shell
# build the service executable
$ GOOS=linux make build

# prepare environment file (see .env.example)
$ cat .env
FILE_BROWSER_PORT=8000
SERVICE_HOST=192.168.0.50
SERVICE_PORT=8080
WEB_PORT=9000
MEDIA_LIBRARY_PATH=/path/to/media/directory

# finally, run the platform detached
$ make run

# to get logs
$ curl "http://${SERVICE_HOST}:8080/logs" | jq
{
  "data": [
    "http://192.168.0.50:43403/log-2023-12-17_15-04-25.833.log"
  ]
}
```

## Build UI

The UI is built & stored in the repo because at the moment, as far as I can tell, flutter does not support building on all architectures, including ARM for raspberry pi (the one I am using for deployment). Hopefully this will improve in the future, then the pre-built UI distribution can be removed from the repo.

To build the UI run the following command. Subsequent deployments will use the updated UI from `ui/dist/`.

```shell
$ make build-ui

# redeploy
$ make run
```

## Further improvements

* add support for
  * other devices (only Chromecast supported)
  * other media sources (only local file system supported)
  * other file extensions (only .mp4 and .srt supported)
* add extra device operations (play speed, queueing, etc.)
* better service feedback
  * better errors (with source file/stack trace)
  * device logs in the app
* add device selection API & selection from UI
* fix bugs
  * fix 'no transition' error
* code TODOs
