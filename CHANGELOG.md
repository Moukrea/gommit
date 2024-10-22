#  (2024-10-22)


### Bug Fixes

* missing import ([8208e40](https://github.com/Moukrea/gommit/commit/8208e40342521c824c3ca1c6d61ddf07095566f7))


### Features

* Implement daily autoupdate check with prompt and automatic update ([1f626c2](https://github.com/Moukrea/gommit/commit/1f626c208baae3fb70994d8f910349d92617181a))



## [0.4.1](https://github.com/Moukrea/gommit/compare/0.4.0...0.4.1) (2024-10-22)


### Bug Fixes

* Remove duplicate failure art and show opposite of success message on failure ([34b2b81](https://github.com/Moukrea/gommit/commit/34b2b81518551180347b27a032ab171d84b5b831))



# [0.4.0](https://github.com/Moukrea/gommit/compare/0.3.4...0.4.0) (2024-10-22)


### Bug Fixes

* Add missing filepath import and ensure .gommit directory exists before checking for updates ([6bd9a7f](https://github.com/Moukrea/gommit/commit/6bd9a7fe8f660e4d29b4a4102ec2cd4fc25af45d))
* integration issues ([ee149a3](https://github.com/Moukrea/gommit/commit/ee149a3b9063a5c767d51de957b6e18275eb25bd))


### Features

* add function to ensure Gommit directory exists ([e3acb2e](https://github.com/Moukrea/gommit/commit/e3acb2e50e2fc73645e78d50280cda73e10e602e))



## [0.3.4](https://github.com/Moukrea/gommit/compare/0.3.3...0.3.4) (2024-10-22)


### Bug Fixes

* setup script sed issue ([e47315a](https://github.com/Moukrea/gommit/commit/e47315af7ab00759d5f3c61b5d7804830e85c2e3))



## [0.3.3](https://github.com/Moukrea/gommit/compare/0.3.2...0.3.3) (2024-10-22)


### Bug Fixes

* setup script sed issue ([2d4a70f](https://github.com/Moukrea/gommit/commit/2d4a70f21f032b0be95770371b02ab18724f3a5b))



## [0.3.2](https://github.com/Moukrea/gommit/compare/0.3.1...0.3.2) (2024-10-22)


### Bug Fixes

* improve setup scripts ([bb3296e](https://github.com/Moukrea/gommit/commit/bb3296ee8641d985393feb47ef4be2a551142aec))



## [0.3.1](https://github.com/Moukrea/gommit/compare/0.3.0...0.3.1) (2024-10-22)


### Bug Fixes

* **integration:** add correct files to gitignore upon integration ([be3cb02](https://github.com/Moukrea/gommit/commit/be3cb023a439b67526647c5de51cbee3c97702bb))



# [0.3.0](https://github.com/Moukrea/gommit/compare/0.2.0...0.3.0) (2024-10-22)


### Bug Fixes

* add 'feat!' to allowed commit types ([9b3851b](https://github.com/Moukrea/gommit/commit/9b3851bb866461ef30795bfd45b277bc4b6bb6e1))
* **autoupdate_test:** remove correct prefix from server URL ([c6ae35b](https://github.com/Moukrea/gommit/commit/c6ae35b9d4f6193e73acf188065cd1176a375d2d))
* **autoupdate_test:** update test cases to use realistic version numbers ([ec233c9](https://github.com/Moukrea/gommit/commit/ec233c93a6efd6d440599304e628b81efdd9ca00))
* **autoupdate:** handle 404 error in getLatestRelease ([e46a94f](https://github.com/Moukrea/gommit/commit/e46a94fe0ba45a1f7d1b155c57433e0d8722d049))
* **autoupdate:** handle full server URL in tests ([f9cc1cf](https://github.com/Moukrea/gommit/commit/f9cc1cf07932e8bb86f95692a0491a03adf17e64))
* **autoupdate:** improve error handling in getLatestRelease ([8100f5e](https://github.com/Moukrea/gommit/commit/8100f5eec5f76f1b733e2125e952fabf2fd85a72))
* breaking change integration on existing footer ([e4ab592](https://github.com/Moukrea/gommit/commit/e4ab59233b8442148141ed92a47cd7ebb4c663b4))
* **config:** use gommit.conf.yaml as the configuration file name and location ([560e00b](https://github.com/Moukrea/gommit/commit/560e00b85eced5c30a6d07c3cd676b3862782e7b))
* footer manipulation issues ([1a8dc39](https://github.com/Moukrea/gommit/commit/1a8dc399f6bf3f09708c39006f8606f90a8e4f56))
* footer manipulation issues ([e921c76](https://github.com/Moukrea/gommit/commit/e921c76a6c19ffe8ad4f0661e6fe97a247096481))
* footer manipulation issues ([e4881b0](https://github.com/Moukrea/gommit/commit/e4881b09a073495d919a3af478d7aeb05d84b684))
* footer manipulation issues ([6d00ada](https://github.com/Moukrea/gommit/commit/6d00ada402ae14b8384fda99408b5d82b794ae16))
* handle unexpected status code in getLatestRelease ([08a299f](https://github.com/Moukrea/gommit/commit/08a299ffa2e0fa5d7e54102c136e332f1d00783b))
* remove "feat!" from AllowedTypes ([55feac5](https://github.com/Moukrea/gommit/commit/55feac563b51e22117fc34e0c4be0ca106acafb6))
* **tests:** issue with ascii art testing ([0d2149c](https://github.com/Moukrea/gommit/commit/0d2149cb5ff26a30818b30dbb325a5fa6368c839))
* **tests:** multiple issues ([4c8391a](https://github.com/Moukrea/gommit/commit/4c8391aeecf4e171885157d8683c8daf4749bb0f))
* update documentation and logic for Gommit configuration ([88c9a6d](https://github.com/Moukrea/gommit/commit/88c9a6d908044df00f415d36bd63824c55c307ed))
* **validation:** fail hook on bad commit message ([e0efcfd](https://github.com/Moukrea/gommit/commit/e0efcfd91324319c65bbb850902c0df00e7cc74b))


### Features

* **ci:** unit tests ([90c831c](https://github.com/Moukrea/gommit/commit/90c831cfc44a84fbb2cfe14a5b6936062bfa0fd1))
* **commit-rules:** add type-enum, type-case, type-empty, scope-case, and subject-empty rules ([77ee638](https://github.com/Moukrea/gommit/commit/77ee6386bb945638ac6a74a6f9acf6fe2009b5d5))
* **config:** make rules customizable via config file ([05f9e27](https://github.com/Moukrea/gommit/commit/05f9e2751dbbbd9ba33ef10b2563f03b7ae9f944))



# [0.2.0](https://github.com/Moukrea/gommit/compare/0.1.8...0.2.0) (2024-10-22)


### Bug Fixes

* **autoupdate:** dupplicate function ([689fdd0](https://github.com/Moukrea/gommit/commit/689fdd0ce729e0d5d72106bb1f9dc2b6f5b0342c))


### Features

* **version:** handle build args and set default version ([e3fc511](https://github.com/Moukrea/gommit/commit/e3fc511cb6ba36e04651cb1d74575ff69727f1dc))



## [0.1.8](https://github.com/Moukrea/gommit/compare/0.1.7...0.1.8) (2024-10-21)


### Bug Fixes

* **integration:** wrong gommit download url ([2724b16](https://github.com/Moukrea/gommit/commit/2724b1621ead564f44187ef3f2b955b2e4c6de60))



## [0.1.7](https://github.com/Moukrea/gommit/compare/0.1.6...0.1.7) (2024-10-21)


### Bug Fixes

* **integration:** makefile task name issue ([e79341b](https://github.com/Moukrea/gommit/commit/e79341b1b23f458733cc5b3a22afd1dad5b7276c))



## [0.1.6](https://github.com/Moukrea/gommit/compare/0.1.5...0.1.6) (2024-10-21)


### Bug Fixes

* **integration:** makefile wrong script names ([b901912](https://github.com/Moukrea/gommit/commit/b90191247ef9f1dc31fe4dd5fc35cbab6be3c159))



## [0.1.5](https://github.com/Moukrea/gommit/compare/0.1.4...0.1.5) (2024-10-21)


### Bug Fixes

* **integration:** links issues ([31460fa](https://github.com/Moukrea/gommit/commit/31460fadf66d380b099de4889054735ca9270107))



## [0.1.4](https://github.com/Moukrea/gommit/compare/0.1.3...0.1.4) (2024-10-21)


### Bug Fixes

* **ci:** breaking change regex ([54ec9ab](https://github.com/Moukrea/gommit/commit/54ec9ab7638f654c02ba812589d5b9089728fd13))



## [0.1.3](https://github.com/Moukrea/gommit/compare/0.1.2...0.1.3) (2024-10-20)


### Bug Fixes

* **downloader:** unused import ([35a6b22](https://github.com/Moukrea/gommit/commit/35a6b223d9e769f8c329b44dd12d5f18fa7e9298))



## [0.1.2](https://github.com/Moukrea/gommit/compare/0.1.1...0.1.2) (2024-10-20)


### Bug Fixes

* **core:** wring download url in binaries ([94ed96a](https://github.com/Moukrea/gommit/commit/94ed96ab66d1e177761f806a2fe9e0ed7ca84a4e))



## [0.1.1](https://github.com/Moukrea/gommit/compare/0.1.0...0.1.1) (2024-10-20)


### Bug Fixes

* **ci:** include binaries in releases ([355f475](https://github.com/Moukrea/gommit/commit/355f475a94cc17c113833b0c09f653244d22aea7))



# [0.1.0](https://github.com/Moukrea/gommit/compare/69312527c5637b50c36248d8f00dc005f77f781a...0.1.0) (2024-10-20)


### Bug Fixes

* **ci:** binaries build issues ([2ab394f](https://github.com/Moukrea/gommit/commit/2ab394fbd2c77b21b8d0ec97f0069f95b6c18c2a))
* **ci:** bump node version ([9c125b1](https://github.com/Moukrea/gommit/commit/9c125b1273bd7280cc71964eca68ba5607b014b0))
* **ci:** dry run skipping most jobs ([38043c6](https://github.com/Moukrea/gommit/commit/38043c6565e12e2682d473c31aa65733ca70a71d))
* **ci:** dry run variable access ([15dbcab](https://github.com/Moukrea/gommit/commit/15dbcab7a4a41ab8798524d99387f78a04eb4d07))
* **ci:** golang version and release notes issues ([ac202ad](https://github.com/Moukrea/gommit/commit/ac202ad75e54ce6dc7386bd6f14e58014b427380))
* **ci:** update commit and tag handling ([6ac46b3](https://github.com/Moukrea/gommit/commit/6ac46b3e1e9fdaa045f2576145ba458ad288220b))
* **ci:** verbose dry run ([5f543f5](https://github.com/Moukrea/gommit/commit/5f543f54584cf5b8d227d62fc3eaebf6769fa344))
* complete gommit implementation ([6931252](https://github.com/Moukrea/gommit/commit/69312527c5637b50c36248d8f00dc005f77f781a))


### Features

* **ci:** introduce automated releases ([c673b09](https://github.com/Moukrea/gommit/commit/c673b0989541bf7d9856bd9091beb2ec0bebe643))



