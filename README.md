<!--
  ~ Copyright 2022 kwanhur
  ~
  ~ Licensed under the Apache License, Version 2.0 (the "License");
  ~ you may not use this file except in compliance with the License.
  ~ You may obtain a copy of the License at
  ~
  ~ http://www.apache.org/licenses/LICENSE-2.0
  ~
  ~ Unless required by applicable law or agreed to in writing, software
  ~ distributed under the License is distributed on an "AS IS" BASIS,
  ~ WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  ~ See the License for the specific language governing permissions and
  ~ limitations under the License.
  ~
-->

# Sixer

## Name

An Apache project repository package verifier.

## Badges

[![Software License](https://img.shields.io/badge/license-Apache2.0-brightgreen.svg?style=for-the-badge)](LICENSE)
[![Powered By: KWANHUR](https://img.shields.io/badge/powered%20by-kwanhur-green.svg?style=for-the-badge)](https://github.com/kwanhur)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)

## Synopsis

```shell
./sixer help
An Apache project repository package verifier

Usage:
  sixer [flags]
  sixer [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  dashboard   apisix dashboard package verifier
  help        Help about any command
  verbose     Show sixer verbose information
  version     Show sixer version number

Flags:
  -a, --announcer string   Specify release candidate announcer
  -c, --candidate string   Specify release candidate version,like 0.2.0
  -C, --commit string      Specify release commit id
  -h, --help               help for sixer
  -t, --timeout uint       Specify request link timeout, unit: second
  -V, --verbose            Show sixer verbose information
  -v, --version            Show sixer version number

Use "sixer [command] --help" for more information about a command.
```

### Example

```shell
./sixer dashboard -a "Zeping Bai" -c 2.11.0 -C 2c563dc15c54a8deb3ba08707594d4d15da76b1b 
2022/03/19 16:56:04 github https://github.com/apache/apisix-dashboard/blob/release/2.11/CHANGELOG.md#2110 validate ok ✅
2022/03/19 16:56:05 github https://github.com/apache/apisix-dashboard/commit/2c563dc15c54a8deb3ba08707594d4d15da76b1b validate ok ✅
2022/03/19 16:56:08 dist https://dist.apache.org/repos/dist/dev/apisix/apisix-dashboard-2.11.0 validate ok ✅
2022/03/19 16:56:09 dist https://dist.apache.org/repos/dist/dev/apisix/apisix-dashboard-2.11.0/apache-apisix-dashboard-2.11.0-src.tgz validate ok ✅
2022/03/19 16:56:10 dist https://dist.apache.org/repos/dist/dev/apisix/apisix-dashboard-2.11.0/apache-apisix-dashboard-2.11.0-src.tgz.asc validate ok ✅
2022/03/19 16:56:11 dist https://dist.apache.org/repos/dist/dev/apisix/apisix-dashboard-2.11.0/apache-apisix-dashboard-2.11.0-src.tgz.sha512 validate ok ✅
2022/03/19 16:56:17 dist validate checksum ok ✅
2022/03/19 16:56:17 dist validate signature ok ✅
2022/03/19 16:56:17 LICENSE ok ✅
2022/03/19 16:56:17 NOTICE ok ✅
```

## TODO

- [x] verfiy github links
- [x] download materials from dist
- [x] verify checksum and signature
- [x] check license and notice
- [ ] summary vote text plain

## License

[Apache License 2.0](LICENSE)
