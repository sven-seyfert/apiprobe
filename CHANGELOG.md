#####

# Changelog

All notable changes to "APIProbe 📡" will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Go to [legend](#legend-types-of-changes) for further information about the types of changes.

## [Unreleased]

## [0.14.0] - 2025-08-26

### Added

- Code coverage (further improvements to follow). [baaa834](https://github.com/sven-seyfert/apiprobe/commit/baaa8349b95d63c61f38032be515485133bb0399)

### Changed

- Increase cURL timeouts to improve execution robustness (some APIs are slower than others). [6e15fe8](https://github.com/sven-seyfert/apiprobe/commit/6e15fe893feee414bd1a8c1e628cb448ada45991)

### Documented

- Add missing function comment in loader.go file. [99dfe3c](https://github.com/sven-seyfert/apiprobe/commit/99dfe3ca46eafa4383f19b1e127cc37682065b74)
- Update README.md file. [8b052b0](https://github.com/sven-seyfert/apiprobe/commit/8b052b0fd0089737874a76ab55631f61441c2537)

### Refactored

- Trivial renamings. [255cfe3](https://github.com/sven-seyfert/apiprobe/commit/255cfe3425c299a7e0dd685c68dc9c4b3c70ba87)
- Add blackbox unit tests for jq.go file. [06230d7](https://github.com/sven-seyfert/apiprobe/commit/06230d7778493c7d4a4d5f6f96e0d24ee8c0ea9a)

## [0.13.0] - 2025-08-20

### Changed

- Replace external JQ dependency (binary) usage by 'gojq' library usage. [2a331da](https://github.com/sven-seyfert/apiprobe/commit/2a331da1ca253ff5a4f1dd5a7fc749e2190a5de4)
- Adjust and extract notification functionality into separate go package. [615c413](https://github.com/sven-seyfert/apiprobe/commit/615c413c22d206531946a2726e482ce444a62d6e)
- Dependency update (and go mod tidy). [aba60f2](https://github.com/sven-seyfert/apiprobe/commit/aba60f249186f8bd02c000ee3dc1b3515d9a18d9)
- Output with sorted JSON nodes (because of the usage of 'gojq' library). [8ffe0f8](https://github.com/sven-seyfert/apiprobe/commit/8ffe0f8d2858cdc01fc9d29085fd353bd1b7613e)

### Documented

- Update README.md file. [3980f18](https://github.com/sven-seyfert/apiprobe/commit/3980f18a1d850779c2d17bee647a6e5f39ee11f0)

### Removed

- External dependency (binary) jq.exe. [702a846](https://github.com/sven-seyfert/apiprobe/commit/702a8465c7ab6618e8fd1d0bcb6a6ca6f9875c17)

## [0.12.0] - 2025-08-19

### Changed

- URL encoding for request params (always) and POST body (special case). [edeb1ee](https://github.com/sven-seyfert/apiprobe/commit/edeb1eedf761a18d6270d653f809cf37577f0ade)
- Set function 'ContainsSubstring' to case-insensitive. [d3af659](https://github.com/sven-seyfert/apiprobe/commit/d3af65945f55f3902b893798c3f90f9a27cedac9)
- Use JSON object as POST body data (and test case data) instead of escaped JSON string. [9b6a3e2](https://github.com/sven-seyfert/apiprobe/commit/9b6a3e2a0483c70ce048a3a037b14696b051424d)
- Handle empty POST body data correctly. [f5482dd](https://github.com/sven-seyfert/apiprobe/commit/f5482dd42cdd11b1a3f1a48db3418e41406c5dfe)
- JSON keys postBody and postBodyData are now JSON objects instead of strings. [f1c2ebf](https://github.com/sven-seyfert/apiprobe/commit/f1c2ebfcdd73cfbd7bbc6e8abccf9426a9875aa5)

### Documented

- Update README.md file. [73f9359](https://github.com/sven-seyfert/apiprobe/commit/73f935938ebc7be305209cc49be5bf8a1805a6a3)

### Refactored

- No extra variable (unnecessary). [685e67c](https://github.com/sven-seyfert/apiprobe/commit/685e67cf813a47285448cf45a7aed5b0af71de03)

### Styled

- Apply golangci-lint suggestions. [15265c6](https://github.com/sven-seyfert/apiprobe/commit/15265c65ee834c3300afd5575276993dd47c45e8)

## [0.11.0] - 2025-08-13

### Added

- CLI option (flag) '--name' to add a custom test run name. [d373e66](https://github.com/sven-seyfert/apiprobe/commit/d373e6651f5a73666b716ea3472cf0fb62acc113)

### Changed

- Use last n characters of string instead of first n characters (in log file output). [64951e4](https://github.com/sven-seyfert/apiprobe/commit/64951e4982c634bbf0c15cf2ef62b5a7290c5cfc)
- Apply CLI option '--name' in batch file example. [0cec95c](https://github.com/sven-seyfert/apiprobe/commit/0cec95c98c4b063ee099dc07beb9935e44977568)

### Documented

- Extend remote execution example (batch file example). [8e30d57](https://github.com/sven-seyfert/apiprobe/commit/8e30d5718e955765c3273c22c200fab049f8c9b2)
- Update README.md file. [5a087da](https://github.com/sven-seyfert/apiprobe/commit/5a087da43d4da7eeee1c9c94fa4c36059be67ea7)

### Styled

- Fix typo in batch file. [d072ca5](https://github.com/sven-seyfert/apiprobe/commit/d072ca58fbdacf03850f5e4c7d2178dc407fbe5a)

## [0.10.0] - 2025-08-07

### Added

- Token store for auth token management (for requests with depending tokens). [ebeea3d](https://github.com/sven-seyfert/apiprobe/commit/ebeea3d7eb40fdc4c07147c2daea197ecc27dd4d)

### Changed

- Dependency update (tidy). [fb52b7d](https://github.com/sven-seyfert/apiprobe/commit/fb52b7d832a9aaa73c0a9819d81aeb6606c28905)

### Documented

- Update README.md file. [a9f9a56](https://github.com/sven-seyfert/apiprobe/commit/a9f9a56d7a55bc4a1ac1122ac4cbbf062e7441a2)
- Add missing function comment. [4383296](https://github.com/sven-seyfert/apiprobe/commit/4383296c941e72fa79f078e2e3d94ea1a4f65619)
- Add authentication section to README.md file. [6b3d42d](https://github.com/sven-seyfert/apiprobe/commit/6b3d42df5feda96ded56b20e5a887a957b8c46e8)

### Refactored

- Apply new golangci-lint settings. [ad854d9](https://github.com/sven-seyfert/apiprobe/commit/ad854d9690fe7d5c1f45fd2bbb5943fc67300176)

### Styled

- Append whitespace to string. [cd7abdb](https://github.com/sven-seyfert/apiprobe/commit/cd7abdbe7d843b4f3ca3d720f37dbed1e1da0058)

## [0.9.0] - 2025-08-06

### Added

- Functionality to execute depending pre-requests first then actual request. [c18434e](https://github.com/sven-seyfert/apiprobe/commit/c18434e5990c1efc558fa42e1cdc4b033ff5a88b)

### Changed

- Tiny adjustment on notification message text. [34716fa](https://github.com/sven-seyfert/apiprobe/commit/34716faf4599d9aaa021038bc2ca86ca2f1735bf)
- Improve Makefile to get cross-platform compilation. [0abe9a9](https://github.com/sven-seyfert/apiprobe/commit/0abe9a90677c9e08968af87b61c57645e3d334eb)

### Documented

- Improve license section readability by code block usage. [ab14777](https://github.com/sven-seyfert/apiprobe/commit/ab1477702f69ff42ab9d0916d50558dc63ba887c)
- Update README.md file. [6dce375](https://github.com/sven-seyfert/apiprobe/commit/6dce375f07d1f70cbb58345875f28a5b0332eb2c)

### Refactored

- Extract filter function in separate 'loader' package. [d509228](https://github.com/sven-seyfert/apiprobe/commit/d509228677ac7bdbb5b3c007c4abb45cdf0786f5)
- Readability improvements by variable renaming. [011a4bf](https://github.com/sven-seyfert/apiprobe/commit/011a4bfbb4c49fccf8cc7fa0c6fe03a9748b2886)

## [0.8.0] - 2025-08-01

### Added

- Exclude feature as command line argument (exclude requests to be run). [6daf963](https://github.com/sven-seyfert/apiprobe/commit/6daf963c08624b9c11f476c1e6573f42689e35c7)

### Changed

- Few minor adjustments and improvements. [07c3bad](https://github.com/sven-seyfert/apiprobe/commit/07c3bad742c6774e06da4863bb251026c5f1c55d)

### Documented

- Improve function comments (more precise and return value description). [69aa925](https://github.com/sven-seyfert/apiprobe/commit/69aa9255bf9458af8a2080bc00b6c3b3c2f8fa24)
- Update README.md file. [67c4be2](https://github.com/sven-seyfert/apiprobe/commit/67c4be2cc1d7e5b3052446577f321c267ef2e7aa)
- Project version bump. [c3fd3c0](https://github.com/sven-seyfert/apiprobe/commit/c3fd3c0dc4950283080910aeab35dc3c963edf69)

### Fixed

- Wrong notification text syntax. [2aab332](https://github.com/sven-seyfert/apiprobe/commit/2aab332a15c36bcf4fa5a5bb8e820899621da988)

### Refactored

- Simply renamings for better readability. [bcff334](https://github.com/sven-seyfert/apiprobe/commit/bcff334d73f81f4d2a0d0178e1d14e6838ad9429)

## [0.7.0] - 2025-06-27

### Changed

- File report.json will be removed each time the requests execution end with no failures or changes. [30fa830](https://github.com/sven-seyfert/apiprobe/commit/30fa8305f0611624e7329533fea2597bdff4f4c4)
- Restructure go structs (APIRequests, more granular). [6299af9](https://github.com/sven-seyfert/apiprobe/commit/6299af9582aad71a831e2ac8f5aa4f68d2b31ff6)
- Apply new go structs values to the relevant locations. [2602ff0](https://github.com/sven-seyfert/apiprobe/commit/2602ff0c2789afe8ca173770c44bb9715106e48c)
- Update test case secret replacement and test case handling. [9e0c135](https://github.com/sven-seyfert/apiprobe/commit/9e0c1359fd0d39e87643504afe1a4096e42265e8)
- Restructure input files (JSON definition files). [c1a4a8d](https://github.com/sven-seyfert/apiprobe/commit/c1a4a8d3cd9a9d1ae7dbca2af9c4078775a4fe30)

### Documented

- Update README.md file. [732a334](https://github.com/sven-seyfert/apiprobe/commit/732a3346b35271f0f4d576fd8257295a7d206e6a)

### Refactored

- Readability improvements. [17e04fe](https://github.com/sven-seyfert/apiprobe/commit/17e04fe46174633e6d0d93d2ef0531bfac381fb8)
- Trivial adjustment (renaming). [c0f99d2](https://github.com/sven-seyfert/apiprobe/commit/c0f99d2d4e0beeb980af082b43fef387b6659bcc)

### Styled

- Fix typo. [156e8ee](https://github.com/sven-seyfert/apiprobe/commit/156e8ee3e7f9ae0db13b41d7a6b1333ca7e20ed0)

## [0.6.1] - 2025-06-20

### Fixed

- Filtering problem when jq filter is empty. [18dd02d](https://github.com/sven-seyfert/apiprobe/commit/18dd02db327f4ad291175eb5f6f52ae5fd42f9aa)

## [0.6.0] - 2025-06-20

### Added

- Heartbeat message notification (configurable interval) instead of run notification on each execution. [3d399ca](https://github.com/sven-seyfert/apiprobe/commit/3d399ca155cd5b0519370d575f9891aae2e1e111)

### Changed

- Extend notification message by hostname. [2c72842](https://github.com/sven-seyfert/apiprobe/commit/2c72842a2279fa16773ec57318e182445b31a6b4)

### Documented

- Update README.md file. [1f6398f](https://github.com/sven-seyfert/apiprobe/commit/1f6398fdc3cb52404e4784d37fde5b352a4679d1)

### Styled

- Trivial code style adjustment. [97704c4](https://github.com/sven-seyfert/apiprobe/commit/97704c41a8d03e61c5296b6a1a238f9ab97eb2bc)

## [0.5.0] - 2025-06-19

### Changed

- Dependency update (SQLite). [172c8cb](https://github.com/sven-seyfert/apiprobe/commit/172c8cb66cba8c66913294f80cf1a8567459c279)
- Log files are stored more granular in month and day folders. [bce77cd](https://github.com/sven-seyfert/apiprobe/commit/bce77cd65e113c2d8f1b3a7324432e0d62b9d101)

## [0.4.0] - 2025-06-16

### Added

- Extracted config approach by config.json file and code to read json values instead of hard coded configuration. [6d93d53](https://github.com/sven-seyfert/apiprobe/commit/6d93d53902c09c88adef14fd0f858fda446aec62)

### Changed

- Replace hard coded WebEx webhook values by config.json entries. [0060fb1](https://github.com/sven-seyfert/apiprobe/commit/0060fb14e105d6a3f53737703d624d14499e2ff8)
- Better exception handling for webhook json key. [418cddf](https://github.com/sven-seyfert/apiprobe/commit/418cddfe1c8c674a4b55a817615d536c2136089d)

### Documented

- Update project logo image. [790bab8](https://github.com/sven-seyfert/apiprobe/commit/790bab81095df023bf3722659943d654ecddf9cf)
- Renew function comments. [ca4f984](https://github.com/sven-seyfert/apiprobe/commit/ca4f9847ef893206b0f286a3a769b6150e091ccf)
- Add notify section (config.json action) to configuration section in the README.md file. [6581b81](https://github.com/sven-seyfert/apiprobe/commit/6581b81f458ecac7ed7aa253734b5d2a21a6f616)
- Update README.md file. [9b6c63f](https://github.com/sven-seyfert/apiprobe/commit/9b6c63f2fd4de7004981fd3e0fc476a850c31276)
- Project version bump. [432b17e](https://github.com/sven-seyfert/apiprobe/commit/432b17e2938bf5d36f2117603ffb10292469274c)

### Fixed

- Missing status code return value. [d613143](https://github.com/sven-seyfert/apiprobe/commit/d6131430b2e14396f2827e40ae1a9e5757016485)

### Refactored

- Tidy dependencies (modules). [8da3941](https://github.com/sven-seyfert/apiprobe/commit/8da3941fdaf63652f160e19e97e0eca26cae4c65)
- Extract program version constant to separate package. [bdc6757](https://github.com/sven-seyfert/apiprobe/commit/bdc67572e436ec6944266da109dee43fe51c7127)
- Reorder program flow regarding config loading. [b88f4a5](https://github.com/sven-seyfert/apiprobe/commit/b88f4a5a5b5c01e10382acc41c43e4fda6bdb5bd)
- Rename WebEx secretSpace to space. [4f39fec](https://github.com/sven-seyfert/apiprobe/commit/4f39fecc0077d6b6d61cfc3dee50d26c5209f8f0)

## [0.3.0] - 2025-06-13

### Added

- Initial commit [GitHub]. [ad47458](https://github.com/sven-seyfert/apiprobe/commit/ad474586e1dd285130aed31aee3596d6014837b0)
- First batch of files of the stable approach. [f885b9f](https://github.com/sven-seyfert/apiprobe/commit/f885b9f0a070418e0b50a63e5eef6c3fefc25528)
- Second batch of files of the stable approach. [0c1c09c](https://github.com/sven-seyfert/apiprobe/commit/0c1c09c8d8191e2026f2e14ac1136ec711a53803)
- Example request input and output result for reqres.in API endpoint users. [a362657](https://github.com/sven-seyfert/apiprobe/commit/a3626573711d2e479a3cffe76dd73a94764859b5)

### Documented

- Update README.md file. [41e88aa](https://github.com/sven-seyfert/apiprobe/commit/41e88aa0147098a0727020f5b45850c9328af4b3)

[Unreleased]: https://github.com/sven-seyfert/apiprobe/compare/v0.14.0...HEAD
[0.14.0]: https://github.com/sven-seyfert/apiprobe/compare/v0.13.0...v0.14.0
[0.13.0]: https://github.com/sven-seyfert/apiprobe/compare/v0.12.0...v0.13.0
[0.12.0]: https://github.com/sven-seyfert/apiprobe/compare/v0.11.0...v0.12.0
[0.11.0]: https://github.com/sven-seyfert/apiprobe/compare/v0.10.0...v0.11.0
[0.10.0]: https://github.com/sven-seyfert/apiprobe/compare/v0.9.0...v0.10.0
[0.9.0]: https://github.com/sven-seyfert/apiprobe/compare/v0.8.0...v0.9.0
[0.8.0]: https://github.com/sven-seyfert/apiprobe/compare/v0.7.0...v0.8.0
[0.7.0]: https://github.com/sven-seyfert/apiprobe/compare/v0.6.1...v0.7.0
[0.6.1]: https://github.com/sven-seyfert/apiprobe/compare/v0.6.0...v0.6.1
[0.6.0]: https://github.com/sven-seyfert/apiprobe/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/sven-seyfert/apiprobe/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/sven-seyfert/apiprobe/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/sven-seyfert/apiprobe/releases/tag/v0.3.0

---

### Legend - Types of changes

- `Added` for new features.
- `Changed` for changes in existing functionality.
- `Deprecated` for soon-to-be removed features.
- `Documented` for documentation only changes.
- `Fixed` for any bug fixes.
- `Refactored` for changes that neither fixes a bug nor adds a feature.
- `Removed` for now removed features.
- `Security` in case of vulnerabilities.
- `Styled` for changes like whitespaces, formatting, missing semicolons etc.

##

[To the top](#changelog)
