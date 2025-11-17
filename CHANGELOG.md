# Changelog

## [3.0.1](https://github.com/rose-pine/rose-pine-bloom/compare/v3.0.0...v3.0.1) (2025-11-17)


### Bug Fixes

* malformed ldflags ([#82](https://github.com/rose-pine/rose-pine-bloom/issues/82)) ([778cedd](https://github.com/rose-pine/rose-pine-bloom/commit/778ceddcaf4b7d31a2c7710f52a41d2796d79607))

## [3.0.0](https://github.com/rose-pine/rose-pine-bloom/compare/v2.2.0...v3.0.0) (2025-11-17)


### ⚠ BREAKING CHANGES

* remove doc command ([#80](https://github.com/rose-pine/rose-pine-bloom/issues/80))
* require explicit template ([#75](https://github.com/rose-pine/rose-pine-bloom/issues/75))
* use cobra cli ([#72](https://github.com/rose-pine/rose-pine-bloom/issues/72))
* use decimals instead of percentages in `hsl-array` format ([#50](https://github.com/rose-pine/rose-pine-bloom/issues/50))
* Format names have been simplified:
    - `hex-ns` is now `hex --plain`
    - `hsl-function` is now `hsl`
    - `rgb-function` is now `rgb`

### Features

* add version flag ([#54](https://github.com/rose-pine/rose-pine-bloom/issues/54)) ([9bb4046](https://github.com/rose-pine/rose-pine-bloom/commit/9bb40468f1cdbd82c5dbfd539c88fa6873201601))
* consolidate build command with logging ([#77](https://github.com/rose-pine/rose-pine-bloom/issues/77)) ([78497b3](https://github.com/rose-pine/rose-pine-bloom/commit/78497b3b45fa8ab59d77f1c5bda6649fb1caf6a2))
* manage meta files ([#55](https://github.com/rose-pine/rose-pine-bloom/issues/55)) ([8e30040](https://github.com/rose-pine/rose-pine-bloom/commit/8e30040705cef1fe9a4040bc8f5165040039faba))
* nix shell ([#69](https://github.com/rose-pine/rose-pine-bloom/issues/69)) ([a13f7b3](https://github.com/rose-pine/rose-pine-bloom/commit/a13f7b3c5c370a8798a1d328c40fd1bcb6b55882))
* remove doc command ([#80](https://github.com/rose-pine/rose-pine-bloom/issues/80)) ([66bfd71](https://github.com/rose-pine/rose-pine-bloom/commit/66bfd71e86015bd8aa7751bac694900e9fd64d1c))
* remove duplicated formats ([#49](https://github.com/rose-pine/rose-pine-bloom/issues/49)) ([85c2762](https://github.com/rose-pine/rose-pine-bloom/commit/85c27625fc14b8d9d24c22f4544b64be498db082))
* require explicit template ([#75](https://github.com/rose-pine/rose-pine-bloom/issues/75)) ([f64f59c](https://github.com/rose-pine/rose-pine-bloom/commit/f64f59cc50cd8f0e9cc3fca3395dcfbd5c7359a4))
* set version during build ([#68](https://github.com/rose-pine/rose-pine-bloom/issues/68)) ([b52e3d5](https://github.com/rose-pine/rose-pine-bloom/commit/b52e3d5163fa9dacbc858663b93da7778bff2f9d))
* use appropriate install/build commands ([#76](https://github.com/rose-pine/rose-pine-bloom/issues/76)) ([7aab83e](https://github.com/rose-pine/rose-pine-bloom/commit/7aab83e69c4a069cea0d4525fc8150e3ce27983a))
* use cobra cli ([#72](https://github.com/rose-pine/rose-pine-bloom/issues/72)) ([59a1e0f](https://github.com/rose-pine/rose-pine-bloom/commit/59a1e0f52ee3b93edbbed71c54b671764374f1f0))
* use decimals instead of percentages in `hsl-array` format ([#50](https://github.com/rose-pine/rose-pine-bloom/issues/50)) ([48768df](https://github.com/rose-pine/rose-pine-bloom/commit/48768df426fd564d3bdc96f19e1629fe536188fa))


### Bug Fixes

* exit after displaying version with -v flag ([#56](https://github.com/rose-pine/rose-pine-bloom/issues/56)) ([5977561](https://github.com/rose-pine/rose-pine-bloom/commit/59775615ca05e065b6e3a04d0823698debec54a2))
* reuse existing license during init ([#79](https://github.com/rose-pine/rose-pine-bloom/issues/79)) ([6356690](https://github.com/rose-pine/rose-pine-bloom/commit/6356690d1436afa3d7ea8668094ae83076042fbb))


### Performance Improvements

* don't process entire string for every variable ([#73](https://github.com/rose-pine/rose-pine-bloom/issues/73)) ([bf1592a](https://github.com/rose-pine/rose-pine-bloom/commit/bf1592affa38e4dc956cb3ab2eb0d4d346be119b))
* optimize color formatting ([#57](https://github.com/rose-pine/rose-pine-bloom/issues/57)) ([14d94c7](https://github.com/rose-pine/rose-pine-bloom/commit/14d94c7136a2ee11abbe263fc3ea8da72311fb7d))

## [2.2.0](https://github.com/rose-pine/rose-pine-bloom/compare/v2.1.0...v2.2.0) (2025-07-02)


### Features

* add `--plain` flag ([#43](https://github.com/rose-pine/rose-pine-bloom/issues/43)) ([2e36a6c](https://github.com/rose-pine/rose-pine-bloom/commit/2e36a6c1e88f6d5aaea736c9aa52099f38f6a6e8))
* add `ansi` format as alias for `rgb-ansi` ([#47](https://github.com/rose-pine/rose-pine-bloom/issues/47)) ([dc8cdcc](https://github.com/rose-pine/rose-pine-bloom/commit/dc8cdcc7fbb686ab539401c014c4e7523d8ac8ea))
* add `appearance` variable as alias for `type` ([#48](https://github.com/rose-pine/rose-pine-bloom/issues/48)) ([1cb3a5c](https://github.com/rose-pine/rose-pine-bloom/commit/1cb3a5c4e7d76d1cb0991dd9f077406c58e9acbf))
* add `rgb-css` format ([#45](https://github.com/rose-pine/rose-pine-bloom/issues/45)) ([223190b](https://github.com/rose-pine/rose-pine-bloom/commit/223190b51119ba947f5a7c2ed8c2f5f208486ee6))
* improve `create` defaults ([#36](https://github.com/rose-pine/rose-pine-bloom/issues/36)) ([6a2c345](https://github.com/rose-pine/rose-pine-bloom/commit/6a2c3450203fe13f187cd06a18322f5510a30f8b))
* support alpha in `hex` format ([#44](https://github.com/rose-pine/rose-pine-bloom/issues/44)) ([bfaeacd](https://github.com/rose-pine/rose-pine-bloom/commit/bfaeacd81d7fcf4052fb23ca6dcde7a6ba0cb5bf))


### Bug Fixes

* update script to target correct README.md file ([#40](https://github.com/rose-pine/rose-pine-bloom/issues/40)) ([d54a733](https://github.com/rose-pine/rose-pine-bloom/commit/d54a7332dfcda9c3acc450e653c77c30fe25ad51))

## [2.1.0](https://github.com/rose-pine/rose-pine-bloom/compare/v2.0.1...v2.1.0) (2025-06-21)


### Features

* add `hsl-css` format ([#26](https://github.com/rose-pine/rose-pine-bloom/issues/26)) ([bfaaa12](https://github.com/rose-pine/rose-pine-bloom/commit/bfaaa121dba7551ca654a095d4ca5437ecaef260))
* add `onaccent` variable ([#33](https://github.com/rose-pine/rose-pine-bloom/issues/33)) ([89afda9](https://github.com/rose-pine/rose-pine-bloom/commit/89afda971c7d9593bacea42a346b7f457925de63))
* add accent autodetection ([#27](https://github.com/rose-pine/rose-pine-bloom/issues/27)) ([c4d1fa5](https://github.com/rose-pine/rose-pine-bloom/commit/c4d1fa54f256ef5563b7f77439f3eff853843905))
* show help for invalid usage ([#32](https://github.com/rose-pine/rose-pine-bloom/issues/32)) ([9ff0f30](https://github.com/rose-pine/rose-pine-bloom/commit/9ff0f3056a17c892ef26c883d24281f86e7a249f))

## [2.0.1](https://github.com/rose-pine/rose-pine-bloom/compare/v2.0.0...v2.0.1) (2025-06-14)


### Bug Fixes

* add missing `--create` flag ([#24](https://github.com/rose-pine/rose-pine-bloom/issues/24)) ([0ba4e58](https://github.com/rose-pine/rose-pine-bloom/commit/0ba4e587e86db5500a480c61457b36653e42e7c7))

## [2.0.0](https://github.com/rose-pine/rose-pine-bloom/compare/v1.0.0...v2.0.0) (2025-06-14)


### ⚠ BREAKING CHANGES

* replace `<hsl,rgb>-ns` with a more generic `--no-commas`
* replace `--template` with a positional argument
* rename `--strip-spaces` to `--no-spaces` ([#12](https://github.com/rose-pine/rose-pine-bloom/issues/12))

### Features

* add `--no-commas` flag ([#17](https://github.com/rose-pine/rose-pine-bloom/issues/17)) ([7c583ce](https://github.com/rose-pine/rose-pine-bloom/commit/7c583ce87df7f060f7dbe152a9fdff9af1f238d2))
* add automatic template detection ([#16](https://github.com/rose-pine/rose-pine-bloom/issues/16)) ([f8a6712](https://github.com/rose-pine/rose-pine-bloom/commit/f8a67123950564c19977b54ed32392e60162741e))
* create template from existing theme ([#8](https://github.com/rose-pine/rose-pine-bloom/issues/8)) ([97c7eea](https://github.com/rose-pine/rose-pine-bloom/commit/97c7eea77fbc807d286da6ae23de41e17a1b8c0f))
* rename `--strip-spaces` to `--no-spaces` ([#12](https://github.com/rose-pine/rose-pine-bloom/issues/12)) ([25d7c92](https://github.com/rose-pine/rose-pine-bloom/commit/25d7c92a5b0dc108af630b9fd528de1779562bd5))
