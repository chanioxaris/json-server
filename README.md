# JSON Server

Create a fake REST API from a json file with zero coding in seconds.

Inspired from the [json-server](https://github.com/typicode/json-server) javascript project.

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
       ],
      "version": 1
    }
    
Start JSON Server

`go run main.go`

If you navigate to http://localhost:8080/posts/1, you will get

    { 
      "id": "1", 
      "title": "json-server", 
      "author": "chanioxaris" 
    }

## Routes
Based on the previous json file, some default routes will be generated

### Plural routes

````
GET     /posts
GET     /posts/{id}
POST    /posts
PUT     /posts/{id}
PATCH   /posts/{id}
DELETE  /posts/{id}
````
### Singular routes

````
GET /version
````

When doing requests, it's good to know that:
- For POST requests any `id` value in the body will be honored, but only if not already taken.
- For POST requests without `id` value in the body, a new one will be generated.
- For PUT requests any `id` value in the body will be ignored, as id values are not mutable.

## Parameters
- You can specify an alternative port with the flag `-p` or `--port`. Default value is `3000`.

`go run main.go -p 4000`

- You can specify an alternative file with the flag `-f` or `--file`. Default value is `db.json`.

`go run main.go -f example.json`


## License

json-server is [MIT licensed](LICENSE).