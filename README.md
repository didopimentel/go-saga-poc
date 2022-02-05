# Go Saga POC

Go Saga POC is a project that implements the Saga pattern for distributed services. It uses gRPC as its main
way of communication between the services.

## Project structure

This project contains has the following structure:

| app - contains the presentation layer for each service (e.g.: delivery, orders, payments APIs)
| domain - contains the domain layer. All microservices have their domain here. 
| extensions: contains everything that may be exported to an external package of its own. In this case, is the Saga extension.
| gateways: contains the logic of communication between services AND between a service and datastores.
| proto: contains the .proto definitions.

It uses Makefile for commands.

## Saga extension

The whole purpose of this project is to make use of the saga extension inside the extensions/ folder. 
The extension consists of two components:

### Saga

Is a wrapper of all the steps with its commands. Each step must have a command and its compensation. The compensation command
is the one that will be executed if upcoming step fails.

### Coordinator

The coordinator is responsible for applying the Saga logic. It exposes a single Execute method that takes a context as
parameter and executes all steps in a procedural manner. Ordering here is crucial.

In order to pass values between commands you need to pass down a SagaContextKey. That context key will be overridden in 
every step.

Rollback is done by defining what is the compensation for each step. It is important to understand that the compensation 
of a step will not be executed if that step was the one that failed. That is because it wouldn't make sense, for example,
to try to undo something that was not done (because it failed). So, for example, if you have the following steps:
1. Create order
2. Create payment
3. Create delivery

And the step 3 fails, then steps 2 and 1 will have their compensation commands executed (in that order).

### Error handling

Errors are stored in two fields that can be retrieved by their getter functions. Compensation is not stopped if any of
them fails.

### Upcoming

- The saga extension will be ejected to its own package.
- Provide a way of defining if compensation must be stopped if any error occurs.