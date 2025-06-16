#####

# Changelog

All notable changes to "APIProbe 📡" will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Go to [legend](#legend-types-of-changes) for further information about the types of changes.

## [Unreleased]

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

[Unreleased]: https://github.com/sven-seyfert/apiprobe/compare/v0.4.0...HEAD
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
