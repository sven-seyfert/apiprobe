#####

<div align="center">
    <img src="assets/images/logo.png" width="90" />
    <h2 align="center">Welcome to <code>APIProbe</code> 📡</h2>
</div>

[![license](https://img.shields.io/badge/license-MPL--2.0-indianred.svg?style=flat-square&logo=spdx&logoColor=white)](https://github.com/sven-seyfert/apiprobe/blob/main/LICENSE.md)
[![release](https://img.shields.io/github/release/sven-seyfert/apiprobe.svg?color=slateblue&style=flat-square&logo=github)](https://github.com/sven-seyfert/apiprobe/releases/latest)
[![docs](https://img.shields.io/badge/docs-github--pages-steelblue.svg?style=flat-square&logo=markdown&logoColor=white)](https://sven-seyfert.github.io/apiprobe/)
[![go report](https://img.shields.io/badge/report-A+-green.svg?style=flat-square&logo=go&logoColor=white)](https://goreportcard.com/report/github.com/sven-seyfert/apiprobe)
[![go coverage](https://img.shields.io/badge/coverage-Ø_20%25-seagreen.svg?style=flat-square&logo=go&logoColor=white)](https://raw.githack.com/sven-seyfert/apiprobe/main/coverage/coverage.html)
[![go.mod version)](https://img.shields.io/github/go-mod/go-version/sven-seyfert/apiprobe?color=lightskyblue&label=go.mod&style=flat-square&logo=go&logoColor=white)](https://github.com/sven-seyfert/apiprobe/blob/main/go.mod)
[![last commit](https://img.shields.io/github/last-commit/sven-seyfert/apiprobe.svg?color=darkgoldenrod&style=flat-square&logo=github)](https://github.com/sven-seyfert/apiprobe/commits/main)
[![contributors](https://img.shields.io/github/contributors/sven-seyfert/apiprobe.svg?color=darkolivegreen&style=flat-square&logo=github)](https://github.com/sven-seyfert/apiprobe/graphs/contributors)

[Description](#description) | [Features](#features) | [Getting started](#getting-started) | [Configuration](#configuration) | [Authentication](#authentication) | [Behind the scenes](#behind-the-scenes) | [Contributing](#contributing) | [License](#license) | [Acknowledgements](#acknowledgements)

---

## Description

#### 🥇 *What*

The project **APIProbe** 📡 is a Go-based lightweight CLI tool designed for automated API monitoring, structured request testing and response change detection. It loads JSON-defined API requests, applies test cases, handles secrets securely, diffs responses and sends webhook notifications when changes or errors occur.

#### 🥈 *Why this*

Unlike GUI-based tools such as Postman, **APIProbe** 📡 is built with developers in mind and optimized for fully automated, data-driven workflows. You can invoke it interactively for quick ad-hoc checks on your local machine or integrate it seamlessly into your CI/CD pipelines for continuous monitoring on remote machines.

#### 🥉 *Stability notice*

Currently in a stable initial state — core features implemented; more advanced capabilities planned.

## Features

- **Structured API definitions**:<br>
  Define and load multiple API requests based on JSON files.

- **Test case support**:<br>
  Define multiple test cases per request to cover various scenarios (data driven approach).

- **Secrets management**:<br>
  Securely store secrets in an encrypted database instead of plain credentials/secrets in the JSON definition files (SQLite).

- **Authentication token handling**:<br>
  Send auth requests, store returned tokens and automatically inject them into dependent requests via `<auth-token>` placeholder.

- **Response diffing**:<br>
  Detect changes through a before and after comparison.

- **Webhook notifications**:<br>
  Send summary reports or error alerts to collaboration tools (like WebEx, MS Teams).

- **Custom Logging**:<br>
  Log to console and log file with multiple log levels.

- **Flexible filtering**:<br>
  Filter by ID or tags, generate new IDs, insert new secrets into database.

## Getting started

🏃‍♂️ [Preconditions](#preconditions) | [Installation](#installation) | [Usage](#usage)

### Preconditions

- Go 1.20+ installed ([download](https://golang.org/dl/)) or simply run the executable instead.
- Ensure `$GOPATH/bin` is in your `PATH`.
- Dependency (binary) `curl` in `./lib/`.
- SQLite available (preinstalled on the most OS and systems).

### Installation

1. Clone the repository:

    ``` bash
    git clone https://github.com/sven-seyfert/apiprobe.git
    cd apiprobe
    ```

2. Ensure dependencies are available:

   `curl` in the `./lib/` folder or adjust the path in `./internal/exec/curl.go`.

3. Run or build the program:

    ``` bash
    # load
    go mod tidy
    go mod download
    ```

    ``` bash
    # run program
    go run main.go

    # or build and run executable
    go build
    ./apiprobe.exe
    ```

    Or see [Makefile](https://github.com/sven-seyfert/apiprobe/blob/main/Makefile) commands.

### Usage

🏃‍♂️ [Global Flags](#global-flags) | [Examples](#examples) | [Remote execution](#remote-execution)

#### *Global Flags*

| Flags                                | Description                                                                                                                    |
| ---                                  | ---                                                                                                                            |
| `--help`                             | Show all flags (switches) and their explanations. Shows also the program version.                                              |
| `--name "Environment: PROD"`         | Set custom name for the test run (for the execution). Shown in final notification.                                             |
| `--id "<hex hash>"`                  | Run only the request matching this ID.                                                                                         |
| `--tags "animals, cars"`             | Run all requests containing any of the comma-separated tags.                                                                   |
| `--exclude "ff00fceb61, bb11abc987"` | Do not run any request that contains the ID of the comma-separated ID list.                                                    |
| `--new-id`                           | Generates and returns a new random hex ID for use in JSON definitions.                                                         |
| `--new-file`                         | Generates a new JSON definition template file. Then enter the request values/data and done.                                    |
| `--add-secret "<value>"`             | Securely stores secrets in SQLite database. Returns a placeholder like "\<secret-b29ff12b50\>"<br>for use in JSON definitions. |

#### *Examples*

- **Run all API requests**:

    ``` bash
    go run main.go
    # or by executable (faster)
    ./apiprobe.exe
    ```

- **Filter and run API requests by ID**:

    ``` bash
    go run main.go --id "ff00fceb61"
    # or by executable (faster)
    ./apiprobe.exe --id "ff00fceb61"
    ```

- **Filter and run API requests by tags**:

    ``` bash
    go run main.go --tags "reqres, booker, env-prod"
    # or by executable (faster)
    ./apiprobe.exe --tags "reqres, booker, env-prod"
    ```

- **Exclude API requests from run by ID**:

    ``` bash
    go run main.go --exclude "ff00fceb61, bb11abc987"
    # or by executable (faster)
    ./apiprobe.exe --exclude "ff00fceb61, bb11abc987"
    ```

    ``` bash
    # combination example:
    # run every request with tag "<tag-name>" except request with specific <ID> and name the test run "MyFirstRun"
    go run main.go --tags "env-prod" --exclude "bb11abc987" --name "MyFirstRun"
    ```

- **Generate new ID**:

    ``` bash
    go run main.go --new-id
    # or by executable (faster)
    ./apiprobe.exe --new-id
    ```

- **Generate new JSON definition template file**:

    ``` bash
    go run main.go --new-file
    # or by executable (faster)
    ./apiprobe.exe --new-file
    ```

- **Add new secret**:

    ``` bash
    go run main.go --add-secret "myApiKey123"
    # or by executable (faster)
    ./apiprobe.exe --add-secret "myApiKey123"
    ```

    For more instructions, see section [secret management](#secret-management) below.

#### *Remote execution*

You can run the CLI regularly via various schedulers or task runners.

> Windows Task Scheduler

A sample XML definition is provided under `./remote/windows-tasks-scheduler.xml`.<br>
Use it to register a scheduled task that invokes `apiprobe.exe` at your desired interval.<br>
For example, to schedule a daily run at 2 AM, import the XML and adjust the `<Triggers>` section accordingly.

## Configuration

🏃‍♂️ [apiprobe.json](#apiprobejson) | [JSON definitions](#json-definitions) | [Secret management](#secret-management)

### apiprobe.json

Setup your webhook URL for WebEx, MS Teams etc. At the moment only WebEx is available (more to be developed).

#### *debugMode*

Activate or deactivate debug mode. This will print the cURL format representation of the request to the console. You then can simply test your request via cURL directly.

#### *heartbeat*

Define the interval (in hours) how often a heartbeat message should be sent. This is useful when you don't receive many failures or changes with you API requests and still want to know is the program running and healthy.

#### *notification*

You can use a placeholder like `<secret-f0f0f0f0f0>` to avoid any plaintext URL section or other secret tokens. To do so, add your URL section to the database using `--add-secret "<value>"` and replace the generated returned placeholder with your original URL section. For more instructions, see section [secret management](#secret-management) below.

### JSON definitions

Define your APIs in JSON files under `./data/input/`. Each file contains an array of objects following the schema:

#### *Minimal definition*

``` json
[
    {
        "id": "0f1e2d3c4b",
        "isActive": true,
        "isAuthRequest": false,
        "preRequestId": "",
        "request": {
            "description": "Short description of the request (purpose)",
            "method": "GET|POST|PUT",
            "url": "https://api.example.com",
            "endpoint": "/api/path",
            "basicAuth": "",
            "headers": [],
            "params": [],
            "postBody": {},
            "name": ""
        },
        "testCases": [
            {
                "name": "",
                "paramsData": "",
                "postBodyData": {}
            }
        ],
        "tags": [
            "env-prod"
        ],
        "jq": ""
    }
]
```

#### *Advanced definition*

``` json
[
    {
        "id": "0f1e2d3c4b",
        "isActive": true,
        "isAuthRequest": false,
        "preRequestId": "",
        "request": {
            "description": "Short description of the request (purpose)",
            "method": "GET|POST|PUT",
            "url": "https://api.example.com",
            "endpoint": "/api/path",
            "basicAuth": "<secret-b4c3d2e1f0>",
            "headers": [
                "Content-Type: application/json"
            ],
            "params": [
                "animalId=1337",
                "pageSize=25",
                "page=3"
            ],
            "postBody": {
                "Username": "John Doe",
                "Password": "<secret-b4c3d2e1f0>"
            },
            "name": ""
        },
        "testCases": [
            {
                "name": "Test with Marry Doe",
                "paramsData": "",
                "postBodyData": {
                    "Username": "Marry Doe",
                    "Password": "<secret-ff00ee11cc>"
                }
            },
            {
                "name": "Test with Julia Ismo",
                "paramsData": "",
                "postBodyData": {
                    "Username": "Julia Ismo",
                    "Password": "<secret-cc11ee00ff>"
                }
            },
            {
                "name": "Test with John Doe and different animalId",
                "paramsData": "animalId=4567",
                "postBodyData": {}
            }
        ],
        "tags": [
            "animals",
            "cars",
            "env-prod"
        ],
        "jq": ".data = (.data | sort_by(.type))"
    }
]
```

#### *Explanation*

Mandatory = (M)<br>
Mandatory for POST request = (P)

| JSON key                   | JSON value description                                                                                                                                                                               | Default value                               |
| --                         | ---                                                                                                                                                                                                  | ---                                         |
| **id** (M)                 | Unique 10 character hex hash. Use `--new-id` to generate.                                                                                                                                            |                                             |
| **isActive** (M)           | Toggle the request execution by this boolean flag. In case the endpoint still exists but is temporary inactive, simply set 'false' and this requests will not be processed.                          | true                                        |
| **isAuthRequest** (M)      | Marks this as an authentication request (e.g. login). When true, the tool will make the resulting token available to subsequent requests.                                                            | false                                       |
| **preRequestId**           | ID of the preconditional request to run before this one. The response payload (e.g. token) of that pre-request will automatically be made available to this request’s headers or body if referenced. | "" (empty string)                           |
| **request**                | JSON node for all request related values.                                                                                                                                                            |                                             |
| **request.description**    | Endpoint description (purpose).                                                                                                                                                                      | "" (empty string)                           |
| **request.method** (M)     | HTTP Method; currently only GET, POST and PUT requests are supported.                                                                                                                                |                                             |
| **request.url** (M)        | Interface (API) URL                                                                                                                                                                                  |                                             |
| **request.endpoint** (M)   | Request endpoint.                                                                                                                                                                                    |                                             |
| **request.basicAuth**      | User and password for a basic authentification; format \<user\>:\<password\>.                                                                                                                        | "" (empty string)                           |
| **request.headers**        | Request header list (one or n headers).                                                                                                                                                              | [] (empty string array)                     |
| **request.params**         | URL query parameter list (one or n params); no ? or & needed, only the raw query parameter(s).                                                                                                       | [] (empty string array)                     |
| **request.postBody** (P)   | JSON message body (payload) for POST requests. Custom JSON object.                                                                                                                                   | {} (empty JSON object)                      |
| **request.name**           | Define the name of the first test case.                                                                                                                                                              | "" (empty string)                           |
| **testCases**              | Data driven test data list (one or n test data entries); these variations apply to query params or post body. See [minimal definition](#minimal-definition).                                         |                                             |
| **testCases.name**         | Define the name of your test case.                                                                                                                                                                   | "" (empty string)                           |
| **testCases.paramsData**   | Define which query parameter should be applied (replaced) in request.params for the test cases. See [advanced definition](#advanced-definition).                                                     | "" (empty string)                           |
| **testCases.postBodyData** | Define post body data that will be applied (replaced) in request.postBody for the test cases. See [advanced definition](#advanced-definition).                                                       | {} (empty JSON object)                      |
| **tags**                   | Representation of the topic, of a application, environment etc.                                                                                                                                      | [] (empty string array)                     |
| **jq**                     | JSON query syntax; prettify JSON response (default ".").                                                                                                                                             | "." (dot is the fallback if "" is provided) |

### Secret management

1. Insert a new secret:

    ``` bash
    go run main.go --add-secret "superSecretValue"
    # or by executable (faster)
    ./apiprobe.exe --add-secret "superSecretValue"
    ```

    You'll receive a message with a placeholder like:

    ```
    Use this placeholder "<secret-ab12cd34ef>" in your JSON file instead of the actual secret value.
    ```

2. In your JSON, replace the real value with this placeholder:

    ```json
    "headers": ["Authorization: Bearer <secret-ab12cd34ef>"]
    ```

    Secrets are securely stored in the SQLite database `./db/store.db`.

## Authentication

This section details how authentication token requests are handled.

1. **Define Auth Requests**: In your JSON definitions set `"isAuthRequest": true` and include the endpoint to obtain your token.
2. **Token Extraction**: Auth responses are parsed (via `jq`) and added to the Token Store under the auth request ID.
3. **Injecting Tokens**: For any request that depends on an auth request, set `"preRequestId": "<auth-request-id>"` **and include the header placeholder**:

   ```json
   "headers": [
     "Authorization: Bearer <auth-token>"
   ]
   ```

   The tool replaces `<auth-token>` with the actual token from the Token Store before executing the request.

4. **Usage**: Ensure your JSON definitions reference `<auth-token>` exactly, so that the CLI can locate and replace it.

## Behind the scenes

🏃‍♂️ [Project layout](#project-layout) | [How it works](#how-it-works) | [Logging, Reporting](#logging-reporting)

### Project layout

Most important parts (directories and files):

``` text
apiprobe/
├── assets/
│   └── images/         # Images, screenshots
├── config/
│   └── apiprobe.json   # User defined config entries (like notification settings)
├── data/
│   ├── input/          # JSON request definitions organized by service and environment
│   └── output/         # Auto-generated responses (snapshots)
├── db/
│   ├── seed.csv        # Initial secrets data
│   └── store.db        # SQLite database
├── internal/           # Go packages
├── lib/                # Dependency binary (curl)
├── logs/               # Execution logs (auto-generated)
├── remote/             # Windows Task Scheduler templates
├── CHANGELOG.md        # Version history
├── LICENSE.md          # MPL-2.0 License
├── main.go             # CLI entrypoint
└── Makefile            # Build & run helpers
```

### How it works

1. **Initialization**: Logger setup, DB connection, CLI flags setup and config load. Also seed default data insertion.
2. **Loading**: Recursively parse JSON files (API request definitions) into `APIRequest` objects.
3. **Filtering**: Apply `--exclude`, `--id` and `--tags` CLI flag filters.
4. **Prepending**: Dependent pre-requests will be merged (prepended) to the list of requests.
5. **Secrets**: Replace `<secret-...>` placeholders with actual secrets from the database.
6. **Authentication**: If an API definition has `isAuthRequest: true`, the response token is stored in an in-memory Token Store keyed by the request ID. For any subsequent requests with `preRequestId`, the `<auth-token>` placeholder in headers is replaced with the stored token before execution.
7. **Execution**:
   - Build cURL arguments and run HTTP request.
   - Capture status code and response body.
   - Filter response body through `jq`.
8. **Diffing**:
   - Compute SHA256 of formatted response.
   - Compare with existing snapshot file in `./data/output`.
   - Update file and record change if different.
9. **Reporting**:
   - Increment counters for errors and changes.
   - Depending on counter results write `./logs/report.json`<br>
     or with suffix `./logs/report-test.json`, `./logs/report-prod.json` depending on `--name` flag content.
   - Send WebEx webhook summary.

### Logging, Reporting

- **Console & file logging**: All logs to console and to file, like `./logs/2025-06/18/2025-06-18-12-58-54.938.log`.
- **Report file**: JSON report at `./logs/report.json` or `./logs/report-test.json` (see above) when errors/changes occur.
- **Webhook**: Automatic notifications to WebEx (WebEx only at the moment). Later also to MS Teams etc.

## Contributing

1. Fork repository.
2. Create feature branch: `git checkout -b feature/my-new-feature`.
3. Commit changes: `git commit -m "Added: My new feature."`.
4. Push to branch: `git push origin feature/my-new-feature`.
5. Open a pull request (PR).

**Please ensure:**

- You added function comments for new functions (\*.go).
- Code passes golangci-lint (`golangci-lint run ./...`).
- You added documentation for new features (README.md).

## License

Copyright (c) 2025 Sven Seyfert (SOLVE-SMART)<br>
Distributed under the MPL-2.0 License. See [LICENSE](https://github.com/sven-seyfert/apiprobe/blob/main/LICENSE.md) for more information.

## Acknowledgements

- Opportunity by [GitHub](https://github.com)
- Badges by [Shields](https://shields.io) and [SimpleIcons](https://simpleicons.org)
- Thanks to the authors, maintainers and contributors of the various projects and products
  - [golang](https://github.com/golang/go) by the Go team at Google; License: [BSD-3-Clause](https://github.com/golang/go/blob/master/LICENSE)
  - [cURL](https://github.com/curl/curl) by Daniel Stenberg; License: [MIT](lib/curl-license.txt)
  - [jq](https://github.com/jqlang/jq) by Stephen Dolan; License: [MIT](https://github.com/jqlang/jq?tab=License-1-ov-file)
  - [gojq](https://github.com/itchyny/gojq) by itchyny; License: [MIT](https://github.com/itchyny/gojq/blob/main/LICENSE)
  - [SQLite](https://www.sqlite.org/copyright.html) by Richard Hipp; License: [Public Domain](https://www.sqlite.org/copyright.html)
  - [go-sqlite](https://github.com/zombiezen/go-sqlite) by Roxy Light; License: [ISC](https://github.com/zombiezen/go-sqlite/blob/main/LICENSE)

##

[To the top](#description)
