<h1 align="center" style="color:#667EEA">
Conman
</h1>
<p align="center">
  <strong>Just yet another configuration store, but simpler.</strong>
</p>

Conman is as standalone http server written in GoLang that stores the configuration in a JSON file. It is ideal for small use cases and can be spun up in most uncomplicated way.

### Features

- Watch for file changes at regular intervals
- Watching of file can be configured using HTTP request
- Get & Set configuration using HTTP request
- Basic Authentication

### Run Locally (Not for production)

- Clone the repository

```bash
  git clone https://github.com/orted-org/conman.git
```

Since there is no dependency to the server, you can simply run the server by

```bash
  go run *.go
```

This will run there server on port `4000`, set a default API_SECRET as `secret@api` create a file named `temp.json` in the current directory and would not watch for file changes.

### Run in Production

- Clone the repository

```bash
  git clone https://github.com/orted-org/conman.git
```

The following environment variables are required:

- `CONMAN_FILENAME` JSON file to get/set configuration (Default: `temp.json`)
- `CONMAN_WATCH_DURATION` Watch duration for the file changes (Default: `-1`, indicating not to watch for file changes)
- `CONMAN_API_SECRET` API secret for basic authentication of the endpoints (Default: `secret@api`)
- `CONMAN_PORT` Port to listen the server (Default: `4000`)

Wit the following environment variables, you can simply run the server by

```bash
  go run *.go
```

Additionally if you want to use Docker, a sample `Dockerfile` & `docker-compose.yml` have been included. Make necessary changes and

```bash
  docker-compose -f docker-compose.yml up --build
```

This will run the server on the port specified in `docker-compose.yml` file or fallback to default port `4000`.

### API reference

#### Every endpoint must contain the Authorization header in request with in the following form `Basic base64-api-secret`

Example:

```
Authorization: Basic c2VjcmV0QGFwaQ==
```

#### Get the basic stats

```http
  GET /stats
```

```json
{
  "status": 200,
  "data": {
    "file_name": "./temp.json",
    "watch_duration": 0
  }
}
```

#### Get all the config

```http
  GET /
```

```json
{
  "status": 200,
  "data": {
    "foo": {
      "bar": 30
    },
    "redis": "13.98.343.90"
  }
}
```

#### Get particular config

```http
  GET /foo
```

```json
{
  "status": 200,
  "data": {
    "foo": "bar"
  }
}
```

#### Set config

```http
  POST /
```

| Parameter | Type     | Description                                |
| :-------- | :------- | :----------------------------------------- |
| `hard`    | `string` | **true/false**. If set config to file also |

#### Set watch duration for file changes

```http
  PUT /watch
```

| Parameter  | Type     | Description                                                        |
| :--------- | :------- | :----------------------------------------------------------------- |
| `duration` | `number` | **required**. Time in seconds. -1 to disable watching file changes |

### Contributing

Contributions are always welcome.
Please feel free to make a PR or let me know via Email.

### License

Licensed under the [MIT](https://choosealicense.com/licenses/mit/) license.
