# Scalable Notification System

This is a project assignment from fastwork company to design the notification system that can scalable on spike traffic.

## Table of Contents

1. [Description](#description)
2. [Installation](#installation)
3. [Usage](#usage)
4. [Running Tests](#running-tests)
5. [Note](#note)

## Description

This project was write in golang and contains multiple services that will working together to give notification system working as expect even facing with spike load. This is also contain the business logic of notification as request.

| Event                              | Noti channels | Noti target |
| ---------------------------------- | ------------- | ----------- |
| Buyer sends chat message to seller | Email, Push   | Seller      |
| Buyer purchases product            | Email, Push   | Seller      |
| Remind to pay for pending order    | Push          | Buyer       |
| Purchased items are shipped        | Push          | Buyer       |

I using structure design of Clean Architecture in go to archive `Independent of Frameworks`, `Testable`, `Independent of Database`, and `Independent of any external agency`. Divided in to 4 layer of each domain.
- Models
- Repository, Adapter
- Usecase
- Handler

### Project Logical view

![Project Diagram](assets/logical-view-v1.png "Project Diagram")

## Installation

### Prerequisites

- docker desktop
  - This project can running on local environment by using docker desktop.
- make
  - To start the project on local environment. This project need make to run preset of command in makefile.
- go v1.21.+ (need if you want to run test)
  - This project implement on go version 1.21.0

### Setup and Start

I prepared the command via makefile. You can execute step by step

```bash
make dev.infra.up
```

this command will initiate many service like kafka, mysql, redpanda(kafkaconsole) for use with notification services.

```bash
make dev.migrate.up
```

this command introduce migration of kafka topics and mysql database to seed some data

```bash
make dev.up
```

or

```bash
make dev.up.watch
```

this command will running all services include server and consumers to perform

## Usage

This project introduce two way of sending notification.

| Usage        | Description                                                                                                                                      | Use case                                                               |
| ------------ | ------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------- |
| Synchronous  | API will process notification immediately. Waiting for notification request sending to the 3rd party then return 200                             | Need to ensure notification request if success                         |
| Asynchronous | API will accept notification request and process later as soon as possible at system best effort. After accepted the request API will return 202 | No need to wait. Appropriate with huge number of notifications request |

### API

| Path                           | Description                          | Success      |
| ------------------------------ | ------------------------------------ | ------------ |
| POST/notification              | Synchronous request                  | 200 OK       |
| POST/notification/asynchronous | Asynchronous request                 | 202 Accepted |
| GET /notification/:reqId       | Retrieve Notification request result | 200 OK       |

by the way I provide [postman](https://www.postman.com/downloads/) collection [here](assets/postman_collection.json) to easy play with this project

## Running Tests
This project has the unit-test for the business logic that implement in use case layer. You can run by below prepared make command.
```bash
make dev.test
```

## Note
This project is just the assignment. Missing many of production idle need such as well security implement, the schema validation, well data schema, and etc. The intension need to introduce the design concept that fit with business requirement such as handle large scale of notification request, show the project structure and code organization, and etc.