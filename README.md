# events-router

<!-- Badger start badges -->
[![Status of the build](https://badger.spt-engprod-pro.mpi-internal.com/badge/travis/Yapo/events-router)](https://travis.mpi-internal.com/Yapo/events-router)
[![Testing Coverage](https://badger.spt-engprod-pro.mpi-internal.com/badge/coverage/Yapo/events-router)](https://reports.spt-engprod-pro.mpi-internal.com/#/Yapo/events-router?branch=master&type=push&daterange&daterange)
[![Style/Linting issues](https://badger.spt-engprod-pro.mpi-internal.com/badge/issues/Yapo/events-router)](https://reports.spt-engprod-pro.mpi-internal.com/#/Yapo/events-router?branch=master&type=push&daterange&daterange)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/flaky_tests/Yapo/events-router)](https://databulous.spt-engprod-pro.mpi-internal.com/test/flaky/Yapo/events-router)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/quality_index/Yapo/events-router)](https://databulous.spt-engprod-pro.mpi-internal.com/quality/repo/Yapo/events-router)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/engprod/Yapo/events-router)](https://github.mpi-internal.com/spt-engprod/badger)
<!-- Badger end badges -->
Goal: Assign to each blocket event a related kafka topic.

Confluence: [Events+Router](https://confluence.mpi-internal.com/display/YAPO/Events+Router)
### Diagram

![flow](https://confluence.mpi-internal.com/rest/gliffy/1.0/embeddedDiagrams/6a5a0e37-82a3-4711-baf4-3dbee2ded96b.png)

## How to run events-router

* Create the dir: `~/go/src/github.mpi-internal.com/Yapo`

* Set the go path: `export GOPATH=~/go` or add the line on your file `.bash_rc`

* Clone this repo:

  ```
  $ cd ~/go/src/github.mpi-internal.com/Yapo
  $ git clone git@github.mpi-internal.com:Yapo/events-router.git
  ```

* On the top dir execute the make instruction to clean and start:

  ```
  $ cd events-router
  $ make start
  ```

* To get a list of available commands:

  ```
  $ make help
  Targets:
    test                 Run tests and generate quality reports
    cover                Run tests and output coverage reports
    coverhtml            Run tests and open report on default web browser
    checkstyle           Run gometalinter and output report as text
    setup                Install golang system level dependencies
    build                Compile the code
    run                  Execute the service
    start                Compile and start the service
    fix-format           Run gofmt to reindent source
    info                 Display basic service info
    docker-build         Create docker image based on docker/dockerfile
    docker-publish       Push docker image to containers.mpi-internal.com
    docker-attach        Attach to this service's currently running docker container output stream
    docker-compose-up    Start all required docker containers for this service
    docker-compose-down  Stop all running docker containers for this service
    help                 This help message
  ```

* If you change the code:

  ```
  $ make start
  ```

* How to run the tests

  ```
  $ make [cover|coverhtml]
  ```

* How to check format

  ```
  $ make checkstyle
  ```

## Endpoints
### GET  /api/v1/healthcheck
Reports whether the service is up and ready to respond.

> When implementing a new service, you MUST keep this endpoint
and update it so it replies according to your service status!

#### Request
No request parameters

#### Response
* Status: Ok message, representing service health

```javascript
200 OK
{
	"Status": "OK"
}
```

## Contact
dev@schibsted.cl
