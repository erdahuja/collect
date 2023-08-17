
# Respository

The project has "app" which exposes two http handlers:
1. APP API:

    it contains the collect API (v1). It is divided into groups for forms, questions, response and users (acl) management. Each group has it's own dependency system so that this can be broken down to microservices in the future.
2. Collect Consumer:

    A compliant to Event interface, it is an independent consumer with it's own config and boot which can be hosted seperately to consume domain events. It is pluggable as consumers can be written for their own problem statements, we may also need to add "enhancers" for data.
    It can be used for multiple upstream API. Currently, it works with csv since google sheets now requires app validation approval for high scope credenitals. 
2. Debug API:

    it contains pprof, liveness, readiness probes as well Go's debug tooling
3. Scratch:

    it contains bootstrap code like migration and seed scripts. it maintains migration versions in db. It can be viewed in [schema](business/data/dbschema/sql/schema.sql), [seed](business/data/dbschema/sql/seed.sql)

As a rule of thumb, import graphs are as follows app imports business imports foundation.

    app: http handlers (can be swapped with socket, rpc easily)
    business: core business logic, data and RBAC system
    foundation: building blocks of a web server (swap router without changing anything!)

For RBAC mechanism, we are using jwt tokens with roles (admin, collector and user) saved in claims. By using a middleware for each handler we are doing authentication and authorization.
A cache layer is also added in users db to quick access user roles.

[DB Design](docs/db/dbb_design.pdf) | [API design](docs/api/collect.md)

## Getting started
For installing Go, please follow
[go official guide](https://go.dev/doc/install)
For installing Kafka, please follow
[blog](https://hevodata.com/learn/install-kafka-on-mac) // or provide broker, topic in dev.env and config.json

> A Makefile is also available to run basic commands. Please use the same. It is available in mac by default

The command line versions can now be installed straight from the command line itself;

    1. Open "Terminal" (it is located in Applications/Utilities)
    2. In the terminal window, run the command xcode-select --install
    3. In the windows that pops up, click Install, and agree to the Terms of Service.

## Commands
`make db`: db can be reset using 

`make run`: to run collect-api

`make sheetsconsumer`: run sheets consumer

`make status-debug`/`make status-api`: server status can be checked running or not

## First Step
Two users will be seeded on running migration (already done if you want to skip the step)

    username: admin@example.com
    password: admin

    username: user1@example.com
    password: user1

Please use [api](docs/api/collect.md#end-point-get-token) for getting bearer token. You have to use basic auth to be able to generate token. (authorization header with username/password)

Once you have the token you can use the same to try out different api. for admin/collector apis will reject/accept as per role defined.

## API

Please download postman collection from [here](https://elements.getpostman.com/redirect?entityId=26793134-37605187-5b1a-4cdf-86b7-c82e7878094c&entityType=collection)
or [![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/26793134-37605187-5b1a-4cdf-86b7-c82e7878094c?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D26793134-37605187-5b1a-4cdf-86b7-c82e7878094c%26entityType%3Dcollection%26workspaceId%3Db66c69c7-0141-4bc1-8932-3265f809fd2d)

| API  | Policy |
| ------------- | ------------- |
| Get token  | public  |
| Create user | ruleAdmin  |
| Get all users  | ruleAdmin  |
| Create form  | ruleAdmin  |
| Get forms  | ruleAdmin,ruleCollector  |
| Delete form  |ruleAdmin  |
| Get all ques for a form  | ruleCollector |
| Create response  | ruleUser  |
| Get all responses for a form  | ruleUser,ruleCollector  |
| Create answer for a response  | ruleCollector  |
| Create question  | ruleAdmin  |
| Get question by id  | ruleCollector  |
| Delete question  | ruleAdmin  |
| server status  | public  |

Documentation:

[API](https://github.com/erdahuja/collect/blob/main/docs/api/collect.md)

## DB Design
POSTGres db is used as the problem had man relationships. however further scale we can add no sql/ cache server to specific problems (parts of app) or to build views.
We are using hosted db of postgreSQL

For pgAdmin exploration, credentials are available in a env file (though a secret manager would be ideal)

Sharing test credentials here
```
DB_USERNAME=oeaualrc
DB_PASSWORD=CNyRbFWEfCc03DdI9PCSkDNZ29AXk1HU
DB_NAME=oeaualrc
DB_HOST=tiny.db.elephantsql.com
```

[Design](docs/db/dbb_design.pdf)

### Benchmark
We can benchmark using Golang inbuilt capability. 
1. Make a benchmark _test file, (e.g collect_test.go).
2. Inside the test file, write benchmark functions using the func Benchmark___(b *testing.B) signature. "___" is functiond descriptive for same.
3. the benchmark function will have the code to measure performance for functionality.
4. Use b.N (b - Go's benchmarking infrastucture) in a loop to specify the number of times the benchmarked code should be executed. It can auto adjusts the total iterations based on the execution time to get meaningful results.
5. Command: ```go test -bench``` to see benchmark results.

### Logs
We are using structured logging to gather logs that be processed in a 3rd party tool or ourself. Uber's zap sugaraed logger is being used in main app and injected in all api groups. Elasticsearch filters can easily be applied.
*I wish i could show onion layer architecture that i use for logging and telemtry*

### Monitoring
Zapkin telemetry is used. IT is injected in Go's context. Telemetry uri and probability rate is provided via config. For brevity sake i have put fake url, however telemtry exists throughout.

### Alerts
OpsGenie, coralogix or newrelic can be used to setup alerts. Also in k8s we can add auto scaling.

### 3rd party
API rate limiters can be implemented in consumers. Also, 3rd party apis often go down, so we need backup providers or create fallback flow using in house backend for top priotity use case. I had a very good use case solved of govt adhar api that's often down, in case it doesn't hamper UX at all. 

For sheets use case, we need to use request batching to avoid multiple "connections", "authentication" from each request. [api](https://developers.google.com/sheets/api/reference/rest/v4/spreadsheets/batchUpdate)

### Database
Replication and sharding is being left out for later discussion

## E2E tests - WIP
We have used mockgen to mock database and tests can be found for business/core. All use interface and modular design.

I will add tests (unit or otherwise) on further discussions only. 