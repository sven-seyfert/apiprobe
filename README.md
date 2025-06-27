#####

<div align="center">
    <img src="assets/images/logo.png" width="90" />
    <h2 align="center">Welcome to <code>APIProbe</code> 📡</h2>
</div>

[![license](https://img.shields.io/badge/license-MPL--2.0-indianred.svg?style=flat-square&logo=spdx&logoColor=white)](https://github.com/sven-seyfert/apiprobe/blob/main/LICENSE.md)
[![release](https://img.shields.io/github/release/sven-seyfert/apiprobe.svg?color=slateblue&style=flat-square&logo=github)](https://github.com/sven-seyfert/apiprobe/releases/latest)
[![docs](https://img.shields.io/badge/docs-github--pages-steelblue.svg?style=flat-square&logo=markdown&logoColor=white)](https://sven-seyfert.github.io/apiprobe/)
[![go report](https://img.shields.io/badge/report-A+-green.svg?style=flat-square&logo=go&logoColor=white)](https://goreportcard.com/report/github.com/sven-seyfert/apiprobe)
[![go coverage](https://img.shields.io/badge/coverage-Ø_0%25-seagreen.svg?style=flat-square&logo=go&logoColor=white)](https://raw.githack.com/sven-seyfert/apiprobe/main/coverage/coverage.html)
[![go.mod version)](https://img.shields.io/github/go-mod/go-version/sven-seyfert/apiprobe?color=lightskyblue&label=go.mod&style=flat-square&logo=go&logoColor=white)](https://github.com/sven-seyfert/apiprobe/blob/main/go.mod)
[![last commit](https://img.shields.io/github/last-commit/sven-seyfert/apiprobe.svg?color=darkgoldenrod&style=flat-square&logo=github)](https://github.com/sven-seyfert/apiprobe/commits/main)
[![contributors](https://img.shields.io/github/contributors/sven-seyfert/apiprobe.svg?color=darkolivegreen&style=flat-square&logo=github)](https://github.com/sven-seyfert/apiprobe/graphs/contributors)

[Description](#description) | [Features](#features) | [Getting started](#getting-started) | [Configuration](#configuration) | [Behind the scenes](#behind-the-scenes) | [Contributing](#contributing) | [License](#license) | [Acknowledgements](#acknowledgements)

---

## Description

#### *What*

The project **APIProbe** 📡 is a Go-based lightweight CLI tool designed for automated API monitoring, structured request testing and response change detection. It loads JSON-defined API requests, applies test cases, handles secrets securely, diffs responses and sends webhook notifications when changes or errors occur.

#### *Why this*

Compared to other API testing or monitoring tools like Postman for example, **APIProbe** 📡 is more developer focused and especially for the automated, data driven (test cases) use cases designed. It's designed to run in CI/CD workflow.

#### *Stability notice*

At the moment the project is in a stable initial state. Means it's on a basic level approach done. Further more advanced features are planned.

## Features

- **Structured API definitions**:<br>
  Define and load multiple API requests based on JSON files (JSON objects).

- **Test case support**:<br>
  Define multiple test cases per request to cover various scenarios (data driven approach).

- **Secrets management**:<br>
  Securely store secrets in an encrypted database instead of plain credentials/secrets in the JSON definition files (SQLite).

- **Response diffing**:<br>
  Detect changes through a before and after comparison.

- **Webhook notifications**:<br>
  Send summary reports or error alerts to collaboration tools (like WebEx, Slack, MS Teams).

- **Custom Logging**:<br>
  Log to console and log file with multiple log levels.

- **Flexible filtering**:<br>
  Filter by ID or tags, generate new IDs, insert new secrets into database.

## Getting started

🏃‍♂️ [Preconditions](#preconditions) | [Installation](#installation) | [Usage](#usage)

### Preconditions

- Go 1.20+ installed ([download](https://golang.org/dl/)) or simply run the executable instead.
- Ensure `$GOPATH/bin` is in your `PATH`.
- Dependency binaries `curl` and `jq` in `./lib/`.
- SQLite available (preinstalled on the most OS and systems).

### Installation

1. Clone the repository:

    ``` bash
    git clone https://github.com/sven-seyfert/apiprobe.git
    cd apiprobe
    ```

2. Ensure dependencies are available:

   `curl` and `jq` in the `./lib/` folder or adjust paths in `./internal/exec/curl.go` and `./internal/exec/jq.go`.

3. Run or build the program:

    ``` get
    # load
    go mod tidy
    go mod download
    ```

    ``` bash
    # run program
    go run main.go

    # or build and run executable
    go build
    ./apiprobe
    ```

    Or see [Makefile](https://github.com/sven-seyfert/apiprobe/blob/main/Makefile) commands.

### Usage

🏃‍♂️ [Global Flags](#global-flags) | [Examples](#examples)

#### *Global Flags*

| Flags                     | Description                                                                                                                    |
| ---                       | ---                                                                                                                            |
| `--id "<hex hash>"`       | Run only the request matching this ID.                                                                                         |
| `--tags "reqres, booker"` | Run all requests containing any of the comma-separated tags.                                                                   |
| `--new-id`                | Generates and returns a new random hex ID for use in JSON definitions.                                                         |
| `--add-secret "<value>"`  | Securely stores secrets in SQLite database. Returns a placeholder like "\<secret-b29ff12b50\>"<br>for use in JSON definitions. |

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

- **Generate new ID**:

    ``` bash
    go run main.go --new-id
    # or by executable (faster)
    ./apiprobe.exe --new-id
    ```

- **Add new secret**:

    ``` bash
    go run main.go --add-secret "myApiKey123"
    # or by executable (faster)
    ./apiprobe.exe --add-secret "myApiKey123"
    ```

    For more instructions, see section [secret management](#secret-management) below.

## Configuration

🏃‍♂️ [config.json](#configjson) | [JSON definitions](#json-definitions) | [Secret management](#secret-management)

### config.json

Setup your webhook URL for WebEx, Slack, MS Teams etc. At the moment only WebEx is available (more to be developed).

#### *heartbeat*

Define the interval (in hours) how often you want to get a heartbeat message. This is useful when you don't receive much failures or changes with you API requests and still want to know is the program running and healthy.

#### *notification*

You can use a placeholder like `<secret-f0f0f0f0f0>` to avoid any plaintext URL section or other secret tokens. To do so, add your URL section to the database using `--add-secret "<value>"` and replace the generated returned placeholder with your original URL section. For more instructions, see section [secret management](#secret-management) below.

### JSON definitions

Define your APIs in JSON files under `./data/input/`. Each file contains an array of objects following the schema:

#### *Minimal definition*

``` json
[
    {
        "id": "0f1e2d3c4b",
        "isAuthRequest": false,
        "preRequestId": "",
        "request": {
            "description": "Short description of the request (purpose)",
            "method": "GET|POST",
            "url": "https://api.example.com",
            "endpoint": "/api/path",
            "basicAuth": "",
            "headers": [],
            "params": [],
            "postBody": ""
        },
        "testCases": [
            {
                "name": "",
                "paramsData": "",
                "postBodyData": ""
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
        "isAuthRequest": false,
        "preRequestId": "",
        "request": {
            "description": "Short description of the request (purpose)",
            "method": "GET|POST",
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
            "postBody": "{\"Username\": \"John Doe\", \"Password\": \"<secret-b4c3d2e1f0>\"}"
        },
        "testCases": [
            {
                "name": "Test with Marry Doe",
                "paramsData": "",
                "postBodyData": "{\"Username\": \"Marry Doe\", \"Password\": \"<secret-ff00ee11cc>\"}"
            },
            {
                "name": "Test with Julia Ismo",
                "paramsData": "",
                "postBodyData": "{\"Username\": \"Julia Ismo\", \"Password\": \"<secret-cc11ee00ff>\"}"
            },
            {
                "name": "Test with John Doe and different animalId",
                "paramsData": "animalId=4567",
                "postBodyData": ""
            }
        ],
        "tags": [
            "animals",
            "cars",
            "env-prod"
        ],
        "jq": "."
    }
]
```

#### *Explanation*

Mandatory = (M)<br>
Mandatory for POST = (P)

| JSON key                   | JSON value description                                                                                                                                                          | Default value                               |
| --                         | ---                                                                                                                                                                             | ---                                         |
| **id** (M)                 | Unique 10 character hex hash. Use `--new-id` to generate.                                                                                                                       |                                             |
| **isAuthRequest** (M)      | Is the request type an auth request (login or similar)? In case of true, the request will be handled differently.                                                               | false                                       |
| **preRequestId**           | Define the request which should be executed before the actual request (because e.g. the<br>preconditional auth request provides a token which is necessary for the actual one). | "" (empty string)                           |
| **request**                | JSON node for all request related values.                                                                                                                                       |                                             |
| **request.description**    | Endpoint description (purpose).                                                                                                                                                 | "" (empty string)                           |
| **request.method** (M)     | HTTP Method; currently only GET and POST requests are supported.                                                                                                                |                                             |
| **request.url** (M)        | Interface (API) URL                                                                                                                                                             |                                             |
| **request.endpoint** (M)   | Request endpoint.                                                                                                                                                               |                                             |
| **request.basicAuth**      | User and password for a basic authentification; format \<user\>:\<password\>.                                                                                                   | "" (empty string)                           |
| **request.headers**        | Request header list (one or n headers).                                                                                                                                         | [] (empty string array)                     |
| **request.params**         | URL query parameter list (one or n params); no ? or & needed, only the raw query parameter(s).                                                                                  | [] (empty string array)                     |
| **request.postBody** (P)   | JSON message body (payload) for POST requests; data have to be JSON valid (escaping " ==> \").                                                                                  | "" (empty string)                           |
| **testCases**              | Data driven test data list (one or n test data entries); these variations apply to query params or post body. See [minimal definition](#minimal-definition).                    |                                             |
| **testCases.name**         | Define the name of your test case.                                                                                                                                              | "" (empty string)                           |
| **testCases.paramsData**   | Define which query parameter should be applied (replaced) in request.params for the test cases. See [advanced definition](#advanced-definition).                                | "" (empty string)                           |
| **testCases.postBodyData** | Define post body data that will be applied (replaced) in request.postBody for the test cases. See [advanced definition](#advanced-definition).                                  | "" (empty string)                           |
| **tags**                   | Representation of the topic, of a application, environment etc.                                                                                                                 | [] (empty string array)                     |
| **jq**                     | JSON query syntax; prettify JSON response (default ".").                                                                                                                        | "." (dot is the fallback if "" is provided) |

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

## Behind the scenes

🏃‍♂️ [Project layout](#project-layout) | [How it works](#how-it-works) | [Logging, Reporting](#logging-reporting)

### Project layout

Most important parts (directories and files):

``` text
apiprobe/
├── assets/
│   └── images/         # Images, screenshots
├── config/
│   └── config.json     # User defined config entries (like notification settings)
├── data/
│   ├── input/          # JSON request definitions organized by service and environment
│   └── output/         # Auto-generated responses (snapshots)
├── db/
│   ├── seed.csv        # Initial secrets data
│   └── store.db        # SQLite database
├── internal/           # Go packages
├── lib/                # Dependency binaries (curl & jq)
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
3. **Filtering**: Apply `--id` or `--tags` CLI flag filters.
4. **Secrets**: Replace `<secret-...>` placeholders with actual secrets.
5. **Execution**:
   - Build curl arguments and run HTTP requests.
   - Capture status codes and response bodies.
   - Filter response through `jq`.
6. **Diffing**:
   - Compute SHA256 of formatted response.
   - Compare with existing snapshot in `./data/output`.
   - Update file and record change if different.
7. **Reporting**:
   - Increment counters for errors and changes.
   - Optionally write `./logs/report.json`.
   - Send WebEx webhook summary.

### Logging, Reporting

- **Console & file logging**: All logs to console and to file, like `./logs/2025-06/18/2025-06-18-12-58-54.938.log`.
- **Report file**: JSON report at `./logs/report.json` when errors/changes occur.
- **Webhook**: Automatic notifications to WebEx (WebEx only at the moment). Later also to Slack, MS Teams etc.

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
- Thanks to the authors, maintainers and contributors of the various projects
  - [golang](https://github.com/golang/go) by the Go team at Google; License: [BSD-3-Clause](https://github.com/golang/go?tab=BSD-3-Clause-1-ov-file#readme)
  - [cURL](https://github.com/curl/curl) by Daniel Stenberg; License: [MIT](https://github.com/sven-seyfert/apiprobe/blob/main/lib/curl-license.txt)
  - [jq](https://github.com/jqlang/jq) by Stephen Dolan; License: [MIT](https://github.com/sven-seyfert/apiprobe/blob/main/lib/jq-license.txt)
  - [SQLite](https://www.sqlite.org/copyright.html) by Richard Hipp; License: [Public Domain](https://www.sqlite.org/copyright.html)

##

[To the top](#description)
