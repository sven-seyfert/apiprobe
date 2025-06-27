#####

# Changelog

All notable changes to "APIProbe 📡" will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Go to [legend](#legend-types-of-changes) for further information about the types of changes.

## [Unreleased]

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

[Unreleased]: https://github.com/sven-seyfert/apiprobe/compare/v0.7.0...HEAD
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
