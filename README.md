# UltimateMicroservice 🏢
[![Go](https://img.shields.io/badge/Go-73d6ec)](https://go.dev/)
[![MySQL](https://img.shields.io/badge/MySQL-eb7c01)](https://www.mysql.com/)
[![Gin](https://img.shields.io/badge/Gin%20Framework-8A2BE2)](https://gin-gonic.com/)
[![Google UUID](https://img.shields.io/badge/Google%20UUID-e30419)](https://github.com/google/uuid)
[![JWT](https://img.shields.io/badge/JWT-d63aff)](https://github.com/google/uuid)
[![Docker](https://img.shields.io/badge/Docker-1d63ed)](https://www.docker.com/)
<br/>

## General
A simple microservice exposing a RESTfull API that handles companies information.

## Table of Contents 🗂️
- [Functionality](#functionality)
- [Directories' Structure](#directories)
- [Installation Process](#installation)

## Functionality ⚙️ <a name="functionality"></a>
The RESTful API exposes the following endpoints (i.e. resources) at port 8080:

|Endpoint|HTTP Method|Usage|Comment|
|---|---|---|---|
|host:8080/company/:id|GET|Fetch the requested company resource based on company id||
|host:8080/company|POST|Insert a new Company resource|User authentication required|
|host:8080/company/:id|DELETE|Delete the specified company resource|User authentication required|
|host:8080/company/:id|PATCH|Modify the requested company resource|User authentication required|
|host:8080/signin|POST|User Autherntication Utility: Login process||
|host:8080/signup|POST|User Autherntication Utility: User creation process||


## Directories' Structure 📂<a name="directories"></a>

As far as directory management is concerned, it was decided to follow the structure mentioned below:

```
UltimateMicroservice
├── database
│   ├── queries
│   │   └── databaseCreation.sql
│   └── Dockerfile
├── microservice
│   ├── cmd
│   │   └── app
│   │       ├── dataManagement
│   │       │   └── dataManagement.go
│   │       ├── model
│   │       │   └── model.go
│   │       └── main.go
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── compose.yaml
└── README.md
```

* The __database__ folder includes the SQL query for the database initialliazation.
* The __microservice__ folder contains the source code of our application, which according to Go standards is stored inside the __cmd/app__ folder. Inside this folder the __main__ package contains the web server, the __model__ package offers a basic data model of the application, while the __dataManagement__ package implements the business logic of this application.

## Installation 📦<a name="installation"></a>

One may use two different ways in order to run this simple application. The first one is to run both the web server services and the database service locally, while the second one is to utilize the power of Docker containers and Docker compose.

### Option I: Local Execution

In this particular case, the installation of the go compiler (version 1.20 or later) is needed for the host machine to run this application, alongside the set up of mySQL server. Further instructions, suitable for each architecture and operating system, can be found at the related websites ([MySQL](https://dev.mysql.com/doc/refman/8.4/en/installing.html) & [Go](https://go.dev/dl/)).

After those steps, the process is straght forward.
Connect to the database and run the __databaseCreation.sql__

```
cd UltimateMicroservice
mysql -u root -p <your_password>
source database/queries/databaseCreation.sql
```

Navigate to the microservice directory, build the project and run.

```
cd UltimateMicroservice/microservice/cmd/app
go mod download
go build .
./app
```

### Option II: Containerization

In this option, the only thing that the user needs to install is the Docker engine & the Docker Compose tool, which is described thoroughly in the related website ([here](https://docs.docker.com/compose/install/)).

The advantages of this approach is the fact that docker will handle the building process of the appliation for us and in fact in a isolated manner.

The first time running the application, we need to run just the following command:

```
cd UltimateMicroservice
docker compose up --build
```
After that, the _--buid_ flag can be omitted. It is important to note that while the compose.yaml file sets a volume for the databse in order to achieve persistance, __the databaseCreation.sql will drop the previously created database and create a new one__. If this is not the desirable behaviour, the only thing that needs to happen is the removal of the _--init-file /docker-entrypoint-initdb.d/databaseCreation.sql_ flag inside the compose.yaml file (located in the parent folder of this directory).

### ⚠️ Important Notice ⚠️

None of the above solutions wil **NOT** work without the creation of a **.evt** file, containing the secrets of the application.

```
cd UltimateMicroservice
touch .env
```

Its internal structure should look like this:

```
DBENDPOINT=<localhost | db> # The endpoint for the database.
MYSQL_DATABASE= "goSchema"
MYSQL_USER= <user for DB access>
MYSQL_PASSWORD= <password of the same user>
JWT_SECRET_KEY=<JWT Key of the Web Server>
```
Localhost should be used in case of bare metal execution, while the db tag should be specified during docker usage (Otherwise the web app container will receive a connection refuse from the database).