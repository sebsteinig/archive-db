# archive-db
## Installation
clone the repository :
```console
git clone git@github.com:WillemNicolas/archive-db.git
```
### With Docker
make sure you have docker installed on your device, otherwise follow the instruction from the docker documentation page : ```https://docs.docker.com/engine/install/```
#### In production mode
 run the following command :
```console
  bash build-prod.sh [OPTIONS] -v, --version x.x
```
If no version is specified, the application will be run with the source code from the main branch, otherwise it will search for the release branch with the specified version.

#### In development mode
 run the following command :
```console
  bash build-dev.sh [OPTIONS] -b, --branch name
```
If no branch is specified, the application will be run with the source code from the develop branch, otherwise it will search for the branch with the specified branch name.
To see the application logs, run the following command :
   ```console
    bash log.sh
   ```
### Without Docker
 1. install go : ```https://go.dev/doc/install```
 2. install postgres : ```https://www.postgresql.org/download/```
 3. create a database and run the [sql script](database/init.sql)
 4. create a .env file in the [api folder](api) with the following variables :
```console
DATABASE_USER="your db user"
DATABASE_PWD="your db password"
DATABASE_URL="postgres://${DATABASE_USER}:${DATABASE_PWD}@your database path/your database name"
API_KEY="your secret api key"
```
5. run the following command in the [api folder](api):
   ```console
     go mod tidy
   ```
6.  to run the application, execute the following command in the [api folder](api):
   ```console
     go run main.go
   ```
7. to see the application logs, run the following command in the [api folder](api) :
   ```console
     tail -f -n "$n" ./archive_api.log
   ```
## API Routes
Documentation for each route can be found in the page at the route /doc.
### Insert  
  There are 4 routes for insert :
  - insert/{id} : (POST) add variables, experiment and execution information in the database. Accept data in the following structure :
       ```json
      {"request" : {
          "table_nimbus_execution" : {
              "exp_id" : "string",
              "config_name" : "string",
              "created_at" : "string (date following ISO format)",
              "extension" : "string",
              "lossless" : "boolean",
              "nan_value_encoding" : "int",
              "threshold" : "float",
              "chunks_time" : "int",
              "chunks_vertical" : "int",
              "rx" : "float",
              "ry" : "float"
          },
          "table_variable" : []{
                  "name" : "string",
                  "path_ts" : "[]string",
                  "path_mean" : "[]string",
                  "levels" : "int",
                  "timesteps" : "int",
                  "xsize" : "int",
                  "xfirst" : "float",
                  "xinc" : "float",
                  "ysize" : "int",
                  "yfirst" : "float",
                  "yinc" : "float",
                  "metadata" : "json object"
              },
          "exp_metadata" : {
              "exp_id" : "string",
              "labels" : []{"label" : "string", "metadata" : "json object"},
              "metadata" : "json object"
          }
      }}
      ```
  - insert/labels/{id} : (POST) add labels in the database. Accept data in the following structure :
    ```json
    {"labels" : [] {"label" : "string", "metadata" : "json object"}}          
    ```
  - insert/publication : (POST) add publication information in the database. Accept data in the following structure :
    ```json
      {"publications" : []{
        "title" : "string",
        "authors_short" : "string",
        "journal" : "string",
        "owner_name" : "string",
        "owner_email" : "string",
        "abstract" : "string",
        "brief_desc" : "string",
        "authors_full" : "string",
        "year" : "int"
      },
      "exp_ids" : []"string"}
    ```
  - insert/clean : (GET) clean the unused data in the database
### Search  
  There are 3 routes for search :
  - search/looking/ : (GET) search labels starting with the character(s) specified in the "for" parameter.
      - query parameters :
          - "for" : string
  - search/ : (GET) search for experiments starting with the character(s) specified in the "like" parameter.
      - query parameters :
          - "like" : string
          - "with" : string
  - search/publication/ : (GET) search publication based on the title, the author and the journal (at least one of these parameters).
### Select  
  There are 2 routes for select :
  - select/{id} : (GET) select experiment by its id
  - select/collection/ : (GET) select multiple experiments by their ids
      - query parameters :
          - ids : []string following json format
          - config_name : string
          - extension : string
          - lossless : bool
          - nan_value_encoding : int
          - threshold : float64
          - chunks_time : int
          - chunks_vertical : int
          - rx : float64
          - ry : float64
