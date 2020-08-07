# JSON Server
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://raw.githubusercontent.com/chanioxaris/json-server/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/chanioxaris/json-server?status.svg)](https://godoc.org/github.com/chanioxaris/json-server)
[![Release](https://img.shields.io/github/release/chanioxaris/json-server)](https://github.com/chanioxaris/json-server/releases/latest)

[![codecov](https://codecov.io/gh/chanioxaris/json-server/branch/master/graph/badge.svg)](https://codecov.io/gh/chanioxaris/json-server)
[![goreportcard](https://goreportcard.com/badge/github.com/chanioxaris/json-server)](https://goreportcard.com/report/github.com/chanioxaris/json-server)

Create a dummy REST API from a json file with zero coding in seconds. Helps you to continue
your development process when an API you depend on, doesn't exist or isn't complete yet. 

Inspired from [json-server](https://github.com/typicode/json-server) javascript package.

## Executables
Cross-platform (Windows, Linux, macOS) binaries are available to [download](https://github.com/chanioxaris/json-server/releases/latest) 
for the latest release. Select the executable that match your OS and run json-server without any dependencies, with 
minimum effort.

## Getting started
Get the package

`go get github.com/chanioxaris/json-server`

Create a `db.json` file with your desired data

    {
      "posts": [
        { 
           "id": "1", 
           "title": "json-server", 
           "author": "chanioxaris" 
        }
      ],
       "books": [
         {
           "id": "1",
           "title": "Clean Code",
           "published": 2008,
           "author": "Robert Martin"
         },
         {
           "id": "2",
           "title": "Crime and punishment",
           "published": 1866,
           "author": "Fyodor Dostoevsky"
         }
       ]
    }
    
Start JSON Server

`go run main.go start`

If you navigate to http://localhost:3000/posts/1, you will get

    { 
      "id": "1", 
      "title": "json-server", 
      "author": "chanioxaris" 
    }

## Routes
Based on the previous json file and for each resource, the below routes will be generated

````
GET     /<resource>
GET     /<resource>/:id
POST    /<resource>
PUT     /<resource>/:id
PATCH   /<resource>/:id
DELETE  /<resource>/:id
````

When doing requests, it's good to know that:
- For POST requests any `id` value in the body will be honored, but only if not already taken.
- For POST requests without `id` value in the body, a new one will be generated.
- For PUT requests any `id` value in the body will be ignored, as id values are not mutable.
- For PATCH requests any `id` value in the body will be ignored, as id values are not mutable.

## Parameters
- You can specify an alternative port with the flag `-p` or `--port`. Default value is `3000`.

`go run main.go start -p 4000`

- You can specify an alternative file with the flag `-f` or `--file`. Default value is `db.json`.

`go run main.go start -f example.json`

- You can toggle http request logs with the flag `-l` or `--logs`. Default value is `false`.

`go run main.go start -l`

## Known issues
- For users running **macOS Catalina** and newer versions, apple will prevent binary from run as it hasn't been notarized 
and signed. To overcome this issue, you can [add a security exception](https://support.apple.com/en-us/HT202491) 
for json-server.

## Contributing

Contributions to this project are welcomed. For more information see the [contribution guidelines](CONTRIBUTING.md).

## License

json-server is [MIT licensed](LICENSE).