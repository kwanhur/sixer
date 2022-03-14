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
# apisixer
An Apache APISIX package verifier

## Synopsis

```shell
apisixer verify --dist dist-url

[]+/- [1|0] no-binding
checked list:
1. packages download links are ok.
2. checksums and signatures are ok.
3. license and notice exit.
```

## TODO
- [ ] download materials from dist
- [ ] verify checksum and signature
- [ ] check license and notice
