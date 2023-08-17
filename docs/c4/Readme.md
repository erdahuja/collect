# Problem statement
The idea is to decouple the survey actions like persistance and it's post processing, since many different clients can come up with different requirements. If we keep changing the backend for every new requirement it will become a maintenance nightmare once the requirement is finished or out dated.
Also, not all clients will share same post processing work load, hence we need to make it such that bringing it down and up is easier.
At database level, since different problems have different query patterns we may want to "enhance" the data and persist it in a query friendly db which can be different from out app db.

# Approaches
## 1. A Monolith survey backend
The idea is to design a monolith and keep developing it to achieve desired results. 

### Pros
1. Only one infra, less moving parts, require fewer dev resources to build
### Cons
1. Downtime while deployment
2. Overtime the project will follow 80:20 rule hence most code will become out dated
3. Different post processing requirements may have different query patterns, scale needs

## 2. Microservices with APIs
The idea is to decouple the survey app and it's downstreams once response is persisted in the db. We can use API to call downstreams which are taking care of processing part.

### Pros
1. Lesser downtime
2. Each server can have it's own tech stack with db based on query patterns
3. Each server can have it's own infra based on req.
### Cons
1. the app needs to have spec to send the request to any dowstreams. The application then awaits the response after ack. As a result, there will always be a 2 way communication, which can add delay in survey request processing. 
2. this delay will be evident for those clients who are not using the particular downstream.

## 3. Microservices with EDD
The idea is to use "Event driven architecture" while following principles of "Domain design". We will decouple our survey app with post processing by creating a Event interface which will comply to our internal domains. It will publish events onto a queue on which different consumers who can comply to "Event" interface can consume and generate their own upstream.
This can be "elastic search" for the search problem, "sheets" for the csv problem and "sql" for the validation problem.

## Pros
1. Each consumer can write it's "enhancer" to update incoming data based on domain and perisist in it's own storage. 
2. the survey app need not to wait for the post processing.
3. Almost No downtime
4. Each consumer can independently scale
## Cons
1. Eventual consistency in complex cases where different consumers process events at different times and are dependent on each other. However our use cases has only one such dependency which can handled.
2. Cost and complexity in Devops, Capex.

### System Context diagram

## Level 1: A System Context diagram provides a starting point, showing how the software system in scope fits into the world around it.

[here](docs/c4/context-diagram.drawio.png)

### Container diagram 
## Level 2: A Container diagram zooms into the software system in scope, showing the high-level technical building blocks.

[here](docs/c4/container-diagram.drawio.png)

### Component diagram 

## Level 3: A Component diagram zooms into an individual container, showing the components inside it.


[here](docs/c4/component-diagram.drawio.png)

## Level 4: UML diagrams (refer low level design and db diagram)