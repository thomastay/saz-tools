## [0.0.11](https://github.com/prantlf/saz-tools/compare/v0.0.10...v0.0.11) (2020-05-24)

### Bug Fixes

* Recognize zero-length responses, when the server just closed the connection ([c674936](https://github.com/prantlf/saz-tools/commit/c67493603708552ff5f95db8820587b82ce898aa))
* Recover a previous session in the table properly ([349c83f](https://github.com/prantlf/saz-tools/commit/349c83f4f1cd65672fd4e1c5c9d1c7597046ba57))
* Remove man pages when uninstalling the NPM package ([7721ce6](https://github.com/prantlf/saz-tools/commit/7721ce6d12c0a8d73c127dbc7b4d31507481d2cd))

### Features

* Add a snap package for Ubuntu ([ebec33f](https://github.com/prantlf/saz-tools/commit/ebec33ff85704629eef6db3fd291220743acab23))
* Introduce a button for redrawing the table if the additionaly shown columns got tool wide ([46c1371](https://github.com/prantlf/saz-tools/commit/46c1371a83c164836cef47466080898def663ed9))

## [0.0.10](https://github.com/prantlf/saz-tools/compare/v0.0.9...v0.0.10) (2020-05-22)

### Bug Fixes

* Remove platform binaries when uninstalling the NPM module ([ff8d5e2](https://github.com/prantlf/saz-tools/commit/ff8d5e219d5e4733a157f302b109292efcd4680a))

### Features

* Add command-line parameter to print the version of the tools ([fe9f7df](https://github.com/prantlf/saz-tools/commit/fe9f7df4debbe61a96afc88ec894d0ae9c33d6c1))
* Add man pages ([9dd7c31](https://github.com/prantlf/saz-tools/commit/9dd7c31457c5282a4f95bb15743d5a6a0e764ea0))
* Hide the help overlay on hitting any key ([cbed9bf](https://github.com/prantlf/saz-tools/commit/cbed9bfcfa571ac3be3fd7a8de75d45df964355e))
* Use botstrap to show tooltips ([0daf5dc](https://github.com/prantlf/saz-tools/commit/0daf5dcaf8b8a08753e2e386a5bca32ccd2dada9))
* Use icons on buttons, make colourng of network sessions optional ([654584f](https://github.com/prantlf/saz-tools/commit/654584fea4ea9fb3eddb518ee8ebcf9863bdd810))

# [0.0.9](https://github.com/prantlf/saz-tools/compare/v0.0.8...v0.0.9) (2020-05-21)

## Features

* Offer a help overlay on the first page opening ([1e71100](https://github.com/prantlf/saz-tools/1e71100b2dcdabaa3e319d66923de46c265c2bcd))
* Distribute binaries using Homebrew ([1e71100](https://github.com/prantlf/saz-tools/876dc4bed3cbbbf87741e0a6ab5f64ee1f7fee2f))
* Distribute binaries using NPM ([1e71100](https://github.com/prantlf/saz-tools/24dde848167ee94828c8a0813c4873e5a0c8ad05))

# [0.0.8](https://github.com/prantlf/saz-tools/compare/v0.0.7...v0.0.8) (2020-05-20)

## Bug Fixes

* Correct computation of aggregated column stats ([da30713](https://github.com/prantlf/saz-tools/da30713688aa92358d79318e2881d6cfbad67a6a))

# [0.0.7](https://github.com/prantlf/saz-tools/compare/v0.0.6...v0.0.7) (2020-05-18)

## Features

* Colour network sessions, add help and detailed error handling, split sources ([128bcb5](https://github.com/prantlf/saz-tools/128bcb51c12272870959ff7678777fe718d49e10))

# [0.0.6](https://github.com/prantlf/saz-tools/compare/v0.0.5...v0.0.6) (2020-05-09)

## Features

* Allow opening and downloading of session details, request body and response body ([2a10baa](https://github.com/prantlf/saz-tools/2a10baaf831a4c80068f95b9609fb90481810c5))

Other changes included:

* Remove the scroller plugin in favour of the page scrolling.
* Replace the multi-part form with the direct .saz file in the REST API.
* Support the HEAD method in the REST API to check the presence of a cached .saz file.
* Highlight the syntax in the output of HTML responses instead of letting them execute by the browser.

# [0.0.5](https://github.com/prantlf/saz-tools/compare/v0.0.4...v0.0.5) (2020-05-08)

## Features

* Show network session details ([4f45dda](https://github.com/prantlf/saz-tools/4f45ddad8a9f2277371a615e8b19390b15e3f5fa))

## Bug Fixes

* Remove the prefix saz from packages in pkg/ ([be663a6](https://github.com/prantlf/saz-tools/be663a6d379c96f618142704698d008844348781))

# [0.0.4](https://github.com/prantlf/saz-tools/compare/v0.0.3...v0.0.4) (2020-05-08)

## Features

* Offer the export to Excel ([15f48d3](https://github.com/prantlf/saz-tools/15f48d34cc1c99ba86098dba1ca81f709091ff07))
* Support drag and drop, accept multiple files ([22214c7](https://github.com/prantlf/saz-tools/22214c7c32c37fac9dc3feea4620b696b1ae697b))

## Bug Fixes

* Do not distribute the build tool move-generated-comments ([648f23c](https://github.com/prantlf/saz-tools/648f23c4d917e5915907511db9d0b18176464f82))

# [0.0.3](https://github.com/prantlf/saz-tools/compare/v0.0.2...v0.0.3) (2020-05-07)

## Features

* Cache SAZ files on the server, store previous SAZ files on the client ([4a163ff](https://github.com/prantlf/saz-tools/4a163ff2a262b5ed664792e8412a31c64de0b041))

## Bug Fixes

* Fix leaking of properties from previously parsed sessions to next ones ([d98918b](https://github.com/prantlf/saz-tools/d98918b23365949c4a01d7c6ca03f667b6fc348d))

# [0.0.2](https://github.com/prantlf/saz-tools/compare/v0.0.1...v0.0.2) (2020-05-07)

## Documentation

* Add GoDoc documentation ([3311850](https://github.com/prantlf/saz-tools/331185019877e370cdb7ba69e5a640212a02d551))

## 0.0.1 (2020-05-06)

Initial release of `sazdump` and `sazserve`.
