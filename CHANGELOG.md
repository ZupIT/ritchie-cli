# Changelog

## [Unreleased](https://github.com/zupit/ritchie-cli/tree/HEAD)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/2.0.5...HEAD)

**Closed issues:**

- Frees use of formulas that do not require the commons repository [\#495](https://github.com/ZupIT/ritchie-cli/issues/495)

## [2.0.5](https://github.com/zupit/ritchie-cli/tree/2.0.5) (2020-09-08)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/2.0.4...2.0.5)

**Closed issues:**

- \[Feature\] - rit clean temp [\#476](https://github.com/ZupIT/ritchie-cli/issues/476)
- \[FEATURE\] Don't execute locally if docker not installed when using flag [\#437](https://github.com/ZupIT/ritchie-cli/issues/437)
- Reduce qnt of parameters in Run method in node [\#432](https://github.com/ZupIT/ritchie-cli/issues/432)
- \[FEATURE\] Suggest to download programming languages \(locally\) if needed [\#427](https://github.com/ZupIT/ritchie-cli/issues/427)
- \[FEATURE\] Functional tests for upgrade [\#419](https://github.com/ZupIT/ritchie-cli/issues/419)
- \[FEATURE\] Functional tests for tutorial [\#418](https://github.com/ZupIT/ritchie-cli/issues/418)
- \[FEATURE\] Functional tests for build formula [\#417](https://github.com/ZupIT/ritchie-cli/issues/417)
- \[FEATURE\] Functional tests for build [\#416](https://github.com/ZupIT/ritchie-cli/issues/416)
- \[FEATURE\] Functional tests for set repo-priority [\#412](https://github.com/ZupIT/ritchie-cli/issues/412)
- \[FEATURE\] Functional tests for delete repo [\#409](https://github.com/ZupIT/ritchie-cli/issues/409)
- \[FEATURE\] Functional tests for delete context [\#408](https://github.com/ZupIT/ritchie-cli/issues/408)
- \[FEATURE\] Functional tests for completion powershell [\#407](https://github.com/ZupIT/ritchie-cli/issues/407)
- \[FEATURE\] Functional tests for completion fish [\#406](https://github.com/ZupIT/ritchie-cli/issues/406)
- Create provider and credential automatically according to the executed formula [\#385](https://github.com/ZupIT/ritchie-cli/issues/385)
- \[FEATURE\] - Check config.json before executing the formula and inform if the user needs to pass credentials. [\#384](https://github.com/ZupIT/ritchie-cli/issues/384)
- \[BUG\] - Add the version of the repository to the docker image generated [\#380](https://github.com/ZupIT/ritchie-cli/issues/380)
- Improved CI to avoid 'killed signal' error [\#497](https://github.com/ZupIT/ritchie-cli/issues/497)
- \[BUG\] STDIN & DOCKER flag don't work together [\#494](https://github.com/ZupIT/ritchie-cli/issues/494)
- Rit list credential - E-mail format [\#483](https://github.com/ZupIT/ritchie-cli/issues/483)
- Metrics acceptance question when rit is updated [\#463](https://github.com/ZupIT/ritchie-cli/issues/463)
- Enable rit execution inside the docker. [\#455](https://github.com/ZupIT/ritchie-cli/issues/455)
- Ask for default execution type \(docker / local\) with RIT INIT. [\#443](https://github.com/ZupIT/ritchie-cli/issues/443)
- When build fail locally, suggest user to use --docker flag [\#439](https://github.com/ZupIT/ritchie-cli/issues/439)
- \[FEATURE\] Functional tests for create formula [\#414](https://github.com/ZupIT/ritchie-cli/issues/414)
- \[BUG\] STDIN doesn't work with input type PASSWORD [\#404](https://github.com/ZupIT/ritchie-cli/issues/404)
- Not adding community repo leads to unpredictable behavior [\#402](https://github.com/ZupIT/ritchie-cli/issues/402)
- \[BUG\] Add repo command may present an incomprehensible error [\#397](https://github.com/ZupIT/ritchie-cli/issues/397)
- Conditional options on config.json [\#390](https://github.com/ZupIT/ritchie-cli/issues/390)
- \[FEATURE\] - Allow formulas without dockerfile and only run locally [\#386](https://github.com/ZupIT/ritchie-cli/issues/386)
- \[BUG\] - Formulas with the same directory structure are causing an error in the build and execution [\#383](https://github.com/ZupIT/ritchie-cli/issues/383)
- Check http status when request stable version [\#373](https://github.com/ZupIT/ritchie-cli/issues/373)

**Merged pull requests:**

- Improved CI to avoid 'killed signal' error [\#503](https://github.com/ZupIT/ritchie-cli/pull/503) ([DittrichLucas](https://github.com/DittrichLucas))
- Remove command lock without initialization [\#496](https://github.com/ZupIT/ritchie-cli/pull/496) ([brunats](https://github.com/brunats))
- Add stdin validation for docker execution [\#493](https://github.com/ZupIT/ritchie-cli/pull/493) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))
- add metrics question on rit upgrade command [\#487](https://github.com/ZupIT/ritchie-cli/pull/487) ([victor-schumacher](https://github.com/victor-schumacher))
-  Duplication of letters in the command list credentials [\#484](https://github.com/ZupIT/ritchie-cli/pull/484) ([victor-schumacher](https://github.com/victor-schumacher))
- Add the 'password' input type to stdin [\#482](https://github.com/ZupIT/ritchie-cli/pull/482) ([kaduartur](https://github.com/kaduartur))
- Ask for default execution type \(docker / local\) with RIT INIT [\#480](https://github.com/ZupIT/ritchie-cli/pull/480) ([kaduartur](https://github.com/kaduartur))
- Added step to question about formula found [\#478](https://github.com/ZupIT/ritchie-cli/pull/478) ([brunats](https://github.com/brunats))
- Add functional test for tutorial [\#465](https://github.com/ZupIT/ritchie-cli/pull/465) ([DittrichLucas](https://github.com/DittrichLucas))
- Release 2.0.4 merge [\#461](https://github.com/ZupIT/ritchie-cli/pull/461) ([zup-ci](https://github.com/zup-ci))
- upgrade when version server don't respond with 200 [\#456](https://github.com/ZupIT/ritchie-cli/pull/456) ([victor-schumacher](https://github.com/victor-schumacher))
- Add functional test for completion powershell [\#486](https://github.com/ZupIT/ritchie-cli/pull/486) ([DittrichLucas](https://github.com/DittrichLucas))
- Add functional test for completion fish [\#479](https://github.com/ZupIT/ritchie-cli/pull/479) ([DittrichLucas](https://github.com/DittrichLucas))
- Add functional test for upgrade [\#469](https://github.com/ZupIT/ritchie-cli/pull/469) ([DittrichLucas](https://github.com/DittrichLucas))
- Adding EMAIL to GITLAB credential [\#467](https://github.com/ZupIT/ritchie-cli/pull/467) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))
- Fix incomprehensible error in command rit add repo [\#466](https://github.com/ZupIT/ritchie-cli/pull/466) ([brunats](https://github.com/brunats))
- Add functional test for build [\#458](https://github.com/ZupIT/ritchie-cli/pull/458) ([DittrichLucas](https://github.com/DittrichLucas))
- Adding conditional steps to config prompt [\#454](https://github.com/ZupIT/ritchie-cli/pull/454) ([henriquemoraeszup](https://github.com/henriquemoraeszup))
- Asking for credential during formula runtime [\#423](https://github.com/ZupIT/ritchie-cli/pull/423) ([henriquemoraes8](https://github.com/henriquemoraes8))
- Mount volume of the .rit folder inside the container [\#403](https://github.com/ZupIT/ritchie-cli/pull/403) ([fabianofernandeszup](https://github.com/fabianofernandeszup))

## [2.0.4](https://github.com/zupit/ritchie-cli/tree/2.0.4) (2020-08-21)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/2.0.3...2.0.4)

**Closed issues:**

- \[FEATURE\] Functional tests for update [\#415](https://github.com/ZupIT/ritchie-cli/issues/415)
- Functional tests for create [\#413](https://github.com/ZupIT/ritchie-cli/issues/413)
- \[FEATURE\] Functional tests for list credential [\#411](https://github.com/ZupIT/ritchie-cli/issues/411)
- \[FEATURE\] Functional tests for list [\#410](https://github.com/ZupIT/ritchie-cli/issues/410)
- \[BUG\] Creating python formulas with dependencies [\#223](https://github.com/ZupIT/ritchie-cli/issues/223)
- \[BUG\]Error on login when password contains special character [\#168](https://github.com/ZupIT/ritchie-cli/issues/168)
- \[BUG\] Conflict between V1 and V2 releases [\#445](https://github.com/ZupIT/ritchie-cli/issues/445)
- Update and improve issues templates [\#421](https://github.com/ZupIT/ritchie-cli/issues/421)
- \[FEATURE\] Add entry by stdin in the init command [\#396](https://github.com/ZupIT/ritchie-cli/issues/396)
- \[FEATURE\] Create new command "rit metrics" [\#374](https://github.com/ZupIT/ritchie-cli/issues/374)
- \[FEATURE\] Process and send metrics to server [\#370](https://github.com/ZupIT/ritchie-cli/issues/370)
- \[FEATURE\] Send formula commands metrics [\#369](https://github.com/ZupIT/ritchie-cli/issues/369)
- \[FEATURE\] Collect commands metrics [\#368](https://github.com/ZupIT/ritchie-cli/issues/368)
- \[FEATURE\] Ask on init if user want to send usage metrics [\#367](https://github.com/ZupIT/ritchie-cli/issues/367)

**Merged pull requests:**

- Add functional test for update [\#457](https://github.com/ZupIT/ritchie-cli/pull/457) ([DittrichLucas](https://github.com/DittrichLucas))
- Standardized metrics command [\#452](https://github.com/ZupIT/ritchie-cli/pull/452) ([kaduartur](https://github.com/kaduartur))
- Add functional test for create [\#448](https://github.com/ZupIT/ritchie-cli/pull/448) ([DittrichLucas](https://github.com/DittrichLucas))
- Add functional test for list credential [\#442](https://github.com/ZupIT/ritchie-cli/pull/442) ([DittrichLucas](https://github.com/DittrichLucas))
- Updating Horus pipeline job [\#434](https://github.com/ZupIT/ritchie-cli/pull/434) ([nathannascimentozup](https://github.com/nathannascimentozup))
- Add functional test for list [\#430](https://github.com/ZupIT/ritchie-cli/pull/430) ([DittrichLucas](https://github.com/DittrichLucas))
- Add entry via stdin to the init command [\#426](https://github.com/ZupIT/ritchie-cli/pull/426) ([brunats](https://github.com/brunats))
- \[ENHANCEMENT\] Hello world formula command in README [\#401](https://github.com/ZupIT/ritchie-cli/pull/401) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))
- Improves inputListCustomMock to use a custom function [\#400](https://github.com/ZupIT/ritchie-cli/pull/400) ([brunats](https://github.com/brunats))
- \[FEATURE\] Add EMAIL to provider variables [\#395](https://github.com/ZupIT/ritchie-cli/pull/395) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))
- Release 2.0.3 merge [\#377](https://github.com/ZupIT/ritchie-cli/pull/377) ([zup-ci](https://github.com/zup-ci))
- Fix tutorial test [\#460](https://github.com/ZupIT/ritchie-cli/pull/460) ([kaduartur](https://github.com/kaduartur))
- \[FIX\] Change pipeline lint [\#453](https://github.com/ZupIT/ritchie-cli/pull/453) ([ernelio](https://github.com/ernelio))
- Formatted metric text in init command [\#451](https://github.com/ZupIT/ritchie-cli/pull/451) ([kaduartur](https://github.com/kaduartur))
- Add collect metric to the main file [\#450](https://github.com/ZupIT/ritchie-cli/pull/450) ([kaduartur](https://github.com/kaduartur))
- Collect Product Metrics [\#449](https://github.com/ZupIT/ritchie-cli/pull/449) ([kaduartur](https://github.com/kaduartur))
- Fix version.sh expression [\#446](https://github.com/ZupIT/ritchie-cli/pull/446) ([kaduartur](https://github.com/kaduartur))
- Create HTTP and gRPC client to metrics [\#441](https://github.com/ZupIT/ritchie-cli/pull/441) ([kaduartur](https://github.com/kaduartur))
- Update tutorial on stdin and add tests [\#435](https://github.com/ZupIT/ritchie-cli/pull/435) ([brunats](https://github.com/brunats))
- Collector function for metrics [\#429](https://github.com/ZupIT/ritchie-cli/pull/429) ([victor-schumacher](https://github.com/victor-schumacher))
- Update issues templates [\#424](https://github.com/ZupIT/ritchie-cli/pull/424) ([brunats](https://github.com/brunats))
- Improves issues templates [\#422](https://github.com/ZupIT/ritchie-cli/pull/422) ([brunasilvazup](https://github.com/brunasilvazup))
- Adds question about metrics in init [\#420](https://github.com/ZupIT/ritchie-cli/pull/420) ([brunats](https://github.com/brunats))
- Metrics check [\#405](https://github.com/ZupIT/ritchie-cli/pull/405) ([victor-schumacher](https://github.com/victor-schumacher))
- \[FEATURE\] add metrics command [\#394](https://github.com/ZupIT/ritchie-cli/pull/394) ([victor-schumacher](https://github.com/victor-schumacher))
- \[FEATURE\] Create implementation to UserIdGenerator [\#379](https://github.com/ZupIT/ritchie-cli/pull/379) ([kaduartur](https://github.com/kaduartur))
- \[ENHANCEMENT\] build formula tutorial [\#378](https://github.com/ZupIT/ritchie-cli/pull/378) ([victor-schumacher](https://github.com/victor-schumacher))
- \[FEATURE\] Create interface for metrics [\#375](https://github.com/ZupIT/ritchie-cli/pull/375) ([kaduartur](https://github.com/kaduartur))

## [2.0.3](https://github.com/zupit/ritchie-cli/tree/2.0.3) (2020-08-07)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0...2.0.3)

**Closed issues:**

- \[BUG\] Create formula with specific words [\#249](https://github.com/ZupIT/ritchie-cli/issues/249)
- \[BUG\] Using rit build formula appends to application.bat instead of replacing [\#238](https://github.com/ZupIT/ritchie-cli/issues/238)

**Merged pull requests:**

- Release 2.0.2 merge [\#365](https://github.com/ZupIT/ritchie-cli/pull/365) ([zup-ci](https://github.com/zup-ci))
- \[FEATURE\] Add set path .msi [\#372](https://github.com/ZupIT/ritchie-cli/pull/372) ([ernelio](https://github.com/ernelio))

## [1.0.0](https://github.com/zupit/ritchie-cli/tree/1.0.0) (2020-08-06)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/2.0.2...1.0.0)

**Merged pull requests:**

- Release 1.0.0 beta.23 [\#339](https://github.com/ZupIT/ritchie-cli/pull/339) ([sandokandias](https://github.com/sandokandias))
- \[Fix\] Credential AWS Provider Default [\#313](https://github.com/ZupIT/ritchie-cli/pull/313) ([fabianofernandeszup](https://github.com/fabianofernandeszup))
- \[FIX\] stable version URL [\#371](https://github.com/ZupIT/ritchie-cli/pull/371) ([kaduartur](https://github.com/kaduartur))

## [2.0.2](https://github.com/zupit/ritchie-cli/tree/2.0.2) (2020-08-06)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/2.0.1...2.0.2)

**Merged pull requests:**

- fix wix [\#364](https://github.com/ZupIT/ritchie-cli/pull/364) ([ernelio](https://github.com/ernelio))
- \[FIX\] Fix Makefile is\_release [\#363](https://github.com/ZupIT/ritchie-cli/pull/363) ([ernelio](https://github.com/ernelio))
- Release 2.0.1 merge [\#362](https://github.com/ZupIT/ritchie-cli/pull/362) ([zup-ci](https://github.com/zup-ci))
- fix message about Ritchie Legacy-1.0.0 [\#360](https://github.com/ZupIT/ritchie-cli/pull/360) ([rodrigomedeirosf](https://github.com/rodrigomedeirosf))

## [2.0.1](https://github.com/zupit/ritchie-cli/tree/2.0.1) (2020-08-06)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/2.0.0...2.0.1)

**Merged pull requests:**

- \[FIX\] Add repo with the commons name [\#361](https://github.com/ZupIT/ritchie-cli/pull/361) ([kaduartur](https://github.com/kaduartur))
- \[Feature\] Add hooks in packagings [\#359](https://github.com/ZupIT/ritchie-cli/pull/359) ([ernelio](https://github.com/ernelio))
- Release 2.0.0 merge [\#358](https://github.com/ZupIT/ritchie-cli/pull/358) ([zup-ci](https://github.com/zup-ci))

## [2.0.0](https://github.com/zupit/ritchie-cli/tree/2.0.0) (2020-08-05)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.24...2.0.0)

**Closed issues:**

- \[BUG\] lint doesn't fail pipeline when find issue [\#347](https://github.com/ZupIT/ritchie-cli/issues/347)
- \[FEATURE\] Adding hello world template in KOTLIN \(rit create formula\) [\#220](https://github.com/ZupIT/ritchie-cli/issues/220)
- \[FEATURE\] Adding hello world template in SWIFT \(rit create formula\) [\#188](https://github.com/ZupIT/ritchie-cli/issues/188)
- \[FEATURE\] Adding hello world template in PERL \(rit create formula\) [\#187](https://github.com/ZupIT/ritchie-cli/issues/187)
- \[FEATURE\] Adding hello world template in RUST \(rit create formula\) [\#181](https://github.com/ZupIT/ritchie-cli/issues/181)
- \[FEATURE\] Adding hello world template in C\# \(rit create formula\) [\#179](https://github.com/ZupIT/ritchie-cli/issues/179)
- \[FEATURE\] rit list credential command [\#178](https://github.com/ZupIT/ritchie-cli/issues/178)

**Merged pull requests:**

- \[Fix\] lint [\#357](https://github.com/ZupIT/ritchie-cli/pull/357) ([kaduartur](https://github.com/kaduartur))
- \[FIX\] init warning msg [\#355](https://github.com/ZupIT/ritchie-cli/pull/355) ([kaduartur](https://github.com/kaduartur))
- Fix notice [\#354](https://github.com/ZupIT/ritchie-cli/pull/354) ([brunats](https://github.com/brunats))
- Added information about legacy version 1.0.0 [\#353](https://github.com/ZupIT/ritchie-cli/pull/353) ([brunats](https://github.com/brunats))
- remove files [\#352](https://github.com/ZupIT/ritchie-cli/pull/352) ([ernelio](https://github.com/ernelio))
- fix license [\#351](https://github.com/ZupIT/ritchie-cli/pull/351) ([ernelio](https://github.com/ernelio))
- Support gitlab provider on command add repo [\#350](https://github.com/ZupIT/ritchie-cli/pull/350) ([kaduartur](https://github.com/kaduartur))
- fix license [\#349](https://github.com/ZupIT/ritchie-cli/pull/349) ([ernelio](https://github.com/ernelio))
- \[FIX\] Fix make build stderr [\#348](https://github.com/ZupIT/ritchie-cli/pull/348) ([kaduartur](https://github.com/kaduartur))
- Fix linter [\#346](https://github.com/ZupIT/ritchie-cli/pull/346) ([victor-schumacher](https://github.com/victor-schumacher))
- Credential through file [\#345](https://github.com/ZupIT/ritchie-cli/pull/345) ([victor-schumacher](https://github.com/victor-schumacher))
- \[FEATURE\] Beautification for windows [\#344](https://github.com/ZupIT/ritchie-cli/pull/344) ([kaduartur](https://github.com/kaduartur))
- \[FEATURE\] Update go-cli-spinner lib [\#340](https://github.com/ZupIT/ritchie-cli/pull/340) ([kaduartur](https://github.com/kaduartur))
- \[Fix\] build local [\#338](https://github.com/ZupIT/ritchie-cli/pull/338) ([kaduartur](https://github.com/kaduartur))
- List repo url [\#337](https://github.com/ZupIT/ritchie-cli/pull/337) ([viniciussousazup](https://github.com/viniciussousazup))
- Fix work space that have space [\#336](https://github.com/ZupIT/ritchie-cli/pull/336) ([viniciussousazup](https://github.com/viniciussousazup))
- change errKeyNotFoundTemplate [\#335](https://github.com/ZupIT/ritchie-cli/pull/335) ([viniciussousazup](https://github.com/viniciussousazup))
- \[FEATURE\] Path for win bin download [\#333](https://github.com/ZupIT/ritchie-cli/pull/333) ([ernelio](https://github.com/ernelio))
- Change build formulas behavior [\#332](https://github.com/ZupIT/ritchie-cli/pull/332) ([kaduartur](https://github.com/kaduartur))
- Change VerifyNewVersion [\#331](https://github.com/ZupIT/ritchie-cli/pull/331) ([viniciussousazup](https://github.com/viniciussousazup))
- Verbose flag [\#330](https://github.com/ZupIT/ritchie-cli/pull/330) ([viniciussousazup](https://github.com/viniciussousazup))
- Improment PR guide lines [\#328](https://github.com/ZupIT/ritchie-cli/pull/328) ([sandokandias](https://github.com/sandokandias))
- Adding ansible credentials and bugfix for 2.0 [\#327](https://github.com/ZupIT/ritchie-cli/pull/327) ([henriquemoraes8](https://github.com/henriquemoraes8))
- remove team and single [\#326](https://github.com/ZupIT/ritchie-cli/pull/326) ([ernelio](https://github.com/ernelio))
- Adding licensing to everything [\#325](https://github.com/ZupIT/ritchie-cli/pull/325) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- Fix/list credential [\#323](https://github.com/ZupIT/ritchie-cli/pull/323) ([victor-schumacher](https://github.com/victor-schumacher))
- Verify Credencial of config.json [\#322](https://github.com/ZupIT/ritchie-cli/pull/322) ([viniciussousazup](https://github.com/viniciussousazup))
- Change Set credencial [\#321](https://github.com/ZupIT/ritchie-cli/pull/321) ([viniciussousazup](https://github.com/viniciussousazup))
- Improvement readme [\#318](https://github.com/ZupIT/ritchie-cli/pull/318) ([sandokandias](https://github.com/sandokandias))
- Fix/build\_docker [\#317](https://github.com/ZupIT/ritchie-cli/pull/317) ([kaduartur](https://github.com/kaduartur))
- Metadata [\#312](https://github.com/ZupIT/ritchie-cli/pull/312) ([viniciussousazup](https://github.com/viniciussousazup))
- Updated help messages [\#311](https://github.com/ZupIT/ritchie-cli/pull/311) ([brunats](https://github.com/brunats))
- Run Horus [\#310](https://github.com/ZupIT/ritchie-cli/pull/310) ([Leonardo-Beda-ZUP](https://github.com/Leonardo-Beda-ZUP))
- Change Upgrade and CommonsUrl [\#309](https://github.com/ZupIT/ritchie-cli/pull/309) ([viniciussousazup](https://github.com/viniciussousazup))
- Change functional tests [\#308](https://github.com/ZupIT/ritchie-cli/pull/308) ([viniciussousazup](https://github.com/viniciussousazup))
- Improvement/repo priority setter test [\#306](https://github.com/ZupIT/ritchie-cli/pull/306) ([miguelhbrito](https://github.com/miguelhbrito))
- Improvement/repo test [\#305](https://github.com/ZupIT/ritchie-cli/pull/305) ([viniciussousazup](https://github.com/viniciussousazup))
- Create journey tips [\#304](https://github.com/ZupIT/ritchie-cli/pull/304) ([brunats](https://github.com/brunats))
- Feature/build formula [\#302](https://github.com/ZupIT/ritchie-cli/pull/302) ([kaduartur](https://github.com/kaduartur))
- List credentials command [\#295](https://github.com/ZupIT/ritchie-cli/pull/295) ([victor-schumacher](https://github.com/victor-schumacher))

## [1.0.0-beta.24](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.24) (2020-08-03)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.23...1.0.0-beta.24)

**Closed issues:**

- Create journey tips [\#289](https://github.com/ZupIT/ritchie-cli/issues/289)

**Merged pull requests:**

- \[FIX\] Fix build formula on Windows [\#307](https://github.com/ZupIT/ritchie-cli/pull/307) ([kaduartur](https://github.com/kaduartur))
- Vulnerability SSL pipeline blocking ignored [\#239](https://github.com/ZupIT/ritchie-cli/pull/239) ([Leonardo-Beda-ZUP](https://github.com/Leonardo-Beda-ZUP))

## [1.0.0-beta.23](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.23) (2020-07-27)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.22...1.0.0-beta.23)

**Merged pull requests:**

- Release 1.0.0-beta.22 merge [\#303](https://github.com/ZupIT/ritchie-cli/pull/303) ([zup-ci](https://github.com/zup-ci))

## [1.0.0-beta.22](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.22) (2020-07-20)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.21...1.0.0-beta.22)

**Merged pull requests:**

- added clean autocompletion [\#301](https://github.com/ZupIT/ritchie-cli/pull/301) ([marcosgmgm](https://github.com/marcosgmgm))
- Release 1.0.0-beta.21 merge [\#299](https://github.com/ZupIT/ritchie-cli/pull/299) ([zup-ci](https://github.com/zup-ci))
- Fix tmp bin dir pattern [\#298](https://github.com/ZupIT/ritchie-cli/pull/298) ([viniciussousazup](https://github.com/viniciussousazup))
- \[Suggest\] Pull review guidelines according to last meeting [\#292](https://github.com/ZupIT/ritchie-cli/pull/292) ([henriquemoraes8](https://github.com/henriquemoraes8))
- Clean formulas command [\#288](https://github.com/ZupIT/ritchie-cli/pull/288) ([henriquemoraes8](https://github.com/henriquemoraes8))

## [1.0.0-beta.21](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.21) (2020-07-20)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.20...1.0.0-beta.21)

**Closed issues:**

- \[BUG\] make test-local doesn't work with space in PATH obtained with PWD command [\#183](https://github.com/ZupIT/ritchie-cli/issues/183)
- \[BUG\] Error when adding kubeconfig too big [\#173](https://github.com/ZupIT/ritchie-cli/issues/173)
- \[BUG\]Local Formula [\#148](https://github.com/ZupIT/ritchie-cli/issues/148)
- \[FEATURE\] Resize website layout [\#140](https://github.com/ZupIT/ritchie-cli/issues/140)
- \[FEATURE\] Adding a new input type : SELECTOR [\#50](https://github.com/ZupIT/ritchie-cli/issues/50)

**Merged pull requests:**

- change path to stable version [\#297](https://github.com/ZupIT/ritchie-cli/pull/297) ([marcosgmgm](https://github.com/marcosgmgm))
- adding powershell and fish and modifying helpers [\#294](https://github.com/ZupIT/ritchie-cli/pull/294) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- \[FIX\] Windows Installer no longer requires admin privilege [\#291](https://github.com/ZupIT/ritchie-cli/pull/291) ([douglasvinter](https://github.com/douglasvinter))
- Add context and -it arg to docker [\#290](https://github.com/ZupIT/ritchie-cli/pull/290) ([JoaoDanielRufino](https://github.com/JoaoDanielRufino))
- Add files entry to `set credential` command [\#287](https://github.com/ZupIT/ritchie-cli/pull/287) ([marcoscostazup](https://github.com/marcoscostazup))
- \[FIX\] Update vendor and remove glide [\#285](https://github.com/ZupIT/ritchie-cli/pull/285) ([ernelio](https://github.com/ernelio))
- \[FIX\]Change legacy-version [\#281](https://github.com/ZupIT/ritchie-cli/pull/281) ([ernelio](https://github.com/ernelio))
- \[FIX\] Node Dockerfile template [\#280](https://github.com/ZupIT/ritchie-cli/pull/280) ([henriquemoraes8](https://github.com/henriquemoraes8))
- \[DEPRECATION\] Removed `rit clean repo` command [\#278](https://github.com/ZupIT/ritchie-cli/pull/278) ([marcoscostazup](https://github.com/marcoscostazup))
- adding legatsy pipeline to allow running old code [\#275](https://github.com/ZupIT/ritchie-cli/pull/275) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- Release 1.0.0-beta.20 merge [\#274](https://github.com/ZupIT/ritchie-cli/pull/274) ([zup-ci](https://github.com/zup-ci))
- Improvement/rit completion [\#273](https://github.com/ZupIT/ritchie-cli/pull/273) ([viniciussousazup](https://github.com/viniciussousazup))
- \[fix\] del in uninstaller hook win [\#267](https://github.com/ZupIT/ritchie-cli/pull/267) ([ernelio](https://github.com/ernelio))
- Add signature in single [\#265](https://github.com/ZupIT/ritchie-cli/pull/265) ([ernelio](https://github.com/ernelio))
- Release 1.0.0-beta.19 merge [\#259](https://github.com/ZupIT/ritchie-cli/pull/259) ([zup-ci](https://github.com/zup-ci))
- FEATURE - Improving information about help [\#256](https://github.com/ZupIT/ritchie-cli/pull/256) ([brunats](https://github.com/brunats))
- \[Feature\] Adding ruby language support [\#252](https://github.com/ZupIT/ritchie-cli/pull/252) ([henriquemoraes8](https://github.com/henriquemoraes8))
- Feature/verbose flag [\#250](https://github.com/ZupIT/ritchie-cli/pull/250) ([antonioolf](https://github.com/antonioolf))
- Fix formula path separator based on os [\#240](https://github.com/ZupIT/ritchie-cli/pull/240) ([felipemdrs](https://github.com/felipemdrs))

## [1.0.0-beta.20](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.20) (2020-07-13)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.19...1.0.0-beta.20)

**Closed issues:**

- \[FEATURE\] Adding hello world template in RUBY \(rit create formula\) [\#186](https://github.com/ZupIT/ritchie-cli/issues/186)
- \[FEATURE\] Enhancement of the rit set credential command on single version [\#177](https://github.com/ZupIT/ritchie-cli/issues/177)

## [1.0.0-beta.19](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.19) (2020-07-08)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.18...1.0.0-beta.19)

**Closed issues:**

- \[FEATURE\] Check the possibility of using vendor [\#47](https://github.com/ZupIT/ritchie-cli/issues/47)

**Merged pull requests:**

- \[fix\] Signature bin windows [\#251](https://github.com/ZupIT/ritchie-cli/pull/251) ([ernelio](https://github.com/ernelio))
- Feature/sign bin [\#248](https://github.com/ZupIT/ritchie-cli/pull/248) ([ernelio](https://github.com/ernelio))
- \[FIX\] run go mod vendor [\#247](https://github.com/ZupIT/ritchie-cli/pull/247) ([kaduartur](https://github.com/kaduartur))
- \[FEATURE\] Add glide and vendor in project [\#244](https://github.com/ZupIT/ritchie-cli/pull/244) ([ernelio](https://github.com/ernelio))
- Release 1.0.0-beta.18 merge [\#243](https://github.com/ZupIT/ritchie-cli/pull/243) ([zup-ci](https://github.com/zup-ci))
- \[FEATURE\] List and add on set credentials [\#241](https://github.com/ZupIT/ritchie-cli/pull/241) ([victor-schumacher](https://github.com/victor-schumacher))
- \[ENHANCEMENT\] simplify PR template [\#237](https://github.com/ZupIT/ritchie-cli/pull/237) ([sandokandias](https://github.com/sandokandias))
- \[FEATURE\] Change Login and Init behavior [\#236](https://github.com/ZupIT/ritchie-cli/pull/236) ([marcoscostazup](https://github.com/marcoscostazup))
- \[FEATURE\] survey prompt [\#233](https://github.com/ZupIT/ritchie-cli/pull/233) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- \[FEATURE\] Build formulas on Windows [\#232](https://github.com/ZupIT/ritchie-cli/pull/232) ([kaduartur](https://github.com/kaduartur))
- modifying final stuff [\#229](https://github.com/ZupIT/ritchie-cli/pull/229) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- Feature/surveyprompt [\#225](https://github.com/ZupIT/ritchie-cli/pull/225) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- Feature/surveyprompt [\#210](https://github.com/ZupIT/ritchie-cli/pull/210) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))

## [1.0.0-beta.18](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.18) (2020-07-06)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.17...1.0.0-beta.18)

**Merged pull requests:**

- fix lint timeout [\#235](https://github.com/ZupIT/ritchie-cli/pull/235) ([ernelio](https://github.com/ernelio))
- \[fix\] lint killed [\#234](https://github.com/ZupIT/ritchie-cli/pull/234) ([ernelio](https://github.com/ernelio))
- Release 1.0.0-beta.17 merge [\#231](https://github.com/ZupIT/ritchie-cli/pull/231) ([zup-ci](https://github.com/zup-ci))
- \[FIX\] improvements rit upgrade [\#215](https://github.com/ZupIT/ritchie-cli/pull/215) ([marcoscostazup](https://github.com/marcoscostazup))

## [1.0.0-beta.17](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.17) (2020-06-30)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.15...1.0.0-beta.17)

**Closed issues:**

- \[FEATURE\] Add to command `rit create formula` the build formula command [\#214](https://github.com/ZupIT/ritchie-cli/issues/214)
- \[BUG\] Windows colors [\#196](https://github.com/ZupIT/ritchie-cli/issues/196)
- \[FEATURE\] Adding hello world template in PHP \(rit create formula\) [\#185](https://github.com/ZupIT/ritchie-cli/issues/185)
- \[BUG\] Create formula error [\#171](https://github.com/ZupIT/ritchie-cli/issues/171)
- \[BUG\]Module Axios not found. [\#149](https://github.com/ZupIT/ritchie-cli/issues/149)

**Merged pull requests:**

- \[fix\] Fix stable qa [\#226](https://github.com/ZupIT/ritchie-cli/pull/226) ([ernelio](https://github.com/ernelio))
- \[feature\] Add integration test in core commands team. [\#224](https://github.com/ZupIT/ritchie-cli/pull/224) ([ernelio](https://github.com/ernelio))
- \[FEATURE\] Docker ritchie-server [\#219](https://github.com/ZupIT/ritchie-cli/pull/219) ([marcosgmgm](https://github.com/marcosgmgm))
- \[FIX\] release creator [\#218](https://github.com/ZupIT/ritchie-cli/pull/218) ([ernelio](https://github.com/ernelio))
- fix dialer when err [\#217](https://github.com/ZupIT/ritchie-cli/pull/217) ([viniciussousazup](https://github.com/viniciussousazup))
- \[FEATURE\] Adding hello world template in PHP \(Issue \#185\) [\#216](https://github.com/ZupIT/ritchie-cli/pull/216) ([antonioolf](https://github.com/antonioolf))
- Feature/cli team tests [\#213](https://github.com/ZupIT/ritchie-cli/pull/213) ([dmbarra](https://github.com/dmbarra))
- \[FEATURE\] Create formula and build [\#212](https://github.com/ZupIT/ritchie-cli/pull/212) ([kaduartur](https://github.com/kaduartur))
- Fix delivery release packaging [\#211](https://github.com/ZupIT/ritchie-cli/pull/211) ([ernelio](https://github.com/ernelio))
- Release 1.0.0-beta.16 merge [\#208](https://github.com/ZupIT/ritchie-cli/pull/208) ([zup-ci](https://github.com/zup-ci))
- Release 1.0.0-beta.15 merge [\#207](https://github.com/ZupIT/ritchie-cli/pull/207) ([zup-ci](https://github.com/zup-ci))
- Release 1.0.0-beta.14 merge [\#206](https://github.com/ZupIT/ritchie-cli/pull/206) ([zup-ci](https://github.com/zup-ci))
- Fix/colors windows [\#200](https://github.com/ZupIT/ritchie-cli/pull/200) ([victor-schumacher](https://github.com/victor-schumacher))
- \[FEATURE\] Change login and init to get totp [\#190](https://github.com/ZupIT/ritchie-cli/pull/190) ([viniciussousazup](https://github.com/viniciussousazup))

## [1.0.0-beta.15](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.15) (2020-06-24)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.14...1.0.0-beta.15)

## [1.0.0-beta.14](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.14) (2020-06-24)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.16...1.0.0-beta.14)

## [1.0.0-beta.16](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.16) (2020-06-24)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.13...1.0.0-beta.16)

**Closed issues:**

- \[FEATURE\] rit build formula command [\#176](https://github.com/ZupIT/ritchie-cli/issues/176)
- \[BUG\] PWD environment variable error passed by formula runner [\#155](https://github.com/ZupIT/ritchie-cli/issues/155)
- \[BUG\] RIT Team: set server with invalid URL [\#91](https://github.com/ZupIT/ritchie-cli/issues/91)
- \[FEATURE\] Run formulas inside a docker container [\#80](https://github.com/ZupIT/ritchie-cli/issues/80)
- \[FEATURE\] Create a beautiful screen after login and logout [\#22](https://github.com/ZupIT/ritchie-cli/issues/22)

**Merged pull requests:**

- Revert "change pgk of init constant" [\#205](https://github.com/ZupIT/ritchie-cli/pull/205) ([marcosgmgm](https://github.com/marcosgmgm))
- Revert "Bug/fix horus" [\#204](https://github.com/ZupIT/ritchie-cli/pull/204) ([marcosgmgm](https://github.com/marcosgmgm))
- change pgk of init constant [\#203](https://github.com/ZupIT/ritchie-cli/pull/203) ([viniciussousazup](https://github.com/viniciussousazup))
- changes [\#202](https://github.com/ZupIT/ritchie-cli/pull/202) ([ernelio](https://github.com/ernelio))
- fix remove horus [\#201](https://github.com/ZupIT/ritchie-cli/pull/201) ([ernelio](https://github.com/ernelio))
- Bug/fix horus [\#199](https://github.com/ZupIT/ritchie-cli/pull/199) ([viniciussousazup](https://github.com/viniciussousazup))
- fix generate msi win [\#198](https://github.com/ZupIT/ritchie-cli/pull/198) ([ernelio](https://github.com/ernelio))
- \[FEATURE\] Supporting password fields in formula config [\#197](https://github.com/ZupIT/ritchie-cli/pull/197) ([sandokandias](https://github.com/sandokandias))
- fix generate version in windows [\#195](https://github.com/ZupIT/ritchie-cli/pull/195) ([ernelio](https://github.com/ernelio))
- \[bug\] fix completion zsh,bash [\#194](https://github.com/ZupIT/ritchie-cli/pull/194) ([viniciussousazup](https://github.com/viniciussousazup))
- Fix/version msi release [\#193](https://github.com/ZupIT/ritchie-cli/pull/193) ([ernelio](https://github.com/ernelio))
- fix msi release [\#192](https://github.com/ZupIT/ritchie-cli/pull/192) ([ernelio](https://github.com/ernelio))
- \[FEATURE\] Pinning ssl [\#191](https://github.com/ZupIT/ritchie-cli/pull/191) ([marcosgmgm](https://github.com/marcosgmgm))
- \[FIX\] Packaging changelog [\#184](https://github.com/ZupIT/ritchie-cli/pull/184) ([ernelio](https://github.com/ernelio))
- \[Feature\] Improvements circleci [\#182](https://github.com/ZupIT/ritchie-cli/pull/182) ([ernelio](https://github.com/ernelio))
- \[FEATURE\] Build formula with Ritchie-cli [\#180](https://github.com/ZupIT/ritchie-cli/pull/180) ([kaduartur](https://github.com/kaduartur))
- cleanup, refactor and improvements regarding beta and nightly version [\#175](https://github.com/ZupIT/ritchie-cli/pull/175) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- Feature/unix tests for init [\#174](https://github.com/ZupIT/ritchie-cli/pull/174) ([dmbarra](https://github.com/dmbarra))
- Update ritchie-bot-config.yml [\#172](https://github.com/ZupIT/ritchie-cli/pull/172) ([viniciussousazup](https://github.com/viniciussousazup))
- Feature/packaging [\#169](https://github.com/ZupIT/ritchie-cli/pull/169) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- change run-tests and scripts to windows [\#167](https://github.com/ZupIT/ritchie-cli/pull/167) ([viniciussousazup](https://github.com/viniciussousazup))
- \[FIX\] trailing slashes removed [\#166](https://github.com/ZupIT/ritchie-cli/pull/166) ([marcoscostazup](https://github.com/marcoscostazup))
- Feature/rit upgrade [\#165](https://github.com/ZupIT/ritchie-cli/pull/165) ([viniciussousazup](https://github.com/viniciussousazup))
- Fix smoke test for release [\#164](https://github.com/ZupIT/ritchie-cli/pull/164) ([dmbarra](https://github.com/dmbarra))
- Feature/beta pipeline [\#163](https://github.com/ZupIT/ritchie-cli/pull/163) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- Feature/colors [\#162](https://github.com/ZupIT/ritchie-cli/pull/162) ([victor-schumacher](https://github.com/victor-schumacher))
- \[FIX\] Formulas messages [\#161](https://github.com/ZupIT/ritchie-cli/pull/161) ([kaduartur](https://github.com/kaduartur))
- removed flags [\#160](https://github.com/ZupIT/ritchie-cli/pull/160) ([victor-schumacher](https://github.com/victor-schumacher))
- Temporarily disable windows jobs [\#159](https://github.com/ZupIT/ritchie-cli/pull/159) ([dmbarra](https://github.com/dmbarra))
- Feature/rework for workflows [\#158](https://github.com/ZupIT/ritchie-cli/pull/158) ([dmbarra](https://github.com/dmbarra))
- \[FIX\] pwd shell [\#157](https://github.com/ZupIT/ritchie-cli/pull/157) ([marcosgmgm](https://github.com/marcosgmgm))
- Fix/tpl [\#156](https://github.com/ZupIT/ritchie-cli/pull/156) ([victor-schumacher](https://github.com/victor-schumacher))
- \[FEATURE\] Login prompt [\#154](https://github.com/ZupIT/ritchie-cli/pull/154) ([marcosgmgm](https://github.com/marcosgmgm))
- Feature/windows functional tests [\#153](https://github.com/ZupIT/ritchie-cli/pull/153) ([dmbarra](https://github.com/dmbarra))
- \[FIX\] Adding trapdoor to error while removing branch [\#152](https://github.com/ZupIT/ritchie-cli/pull/152) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- Release 1.0.0-beta.13 merge [\#147](https://github.com/ZupIT/ritchie-cli/pull/147) ([zup-ci](https://github.com/zup-ci))
- \[Feature\] color messages [\#146](https://github.com/ZupIT/ritchie-cli/pull/146) ([victor-schumacher](https://github.com/victor-schumacher))
- \[FEATURE\] run formula docker [\#113](https://github.com/ZupIT/ritchie-cli/pull/113) ([kaduartur](https://github.com/kaduartur))

## [1.0.0-beta.13](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.13) (2020-06-03)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.12...1.0.0-beta.13)

**Merged pull requests:**

- changed folder name [\#145](https://github.com/ZupIT/ritchie-cli/pull/145) ([victor-schumacher](https://github.com/victor-schumacher))
- Release 1.0.0-beta.12 merge [\#144](https://github.com/ZupIT/ritchie-cli/pull/144) ([zup-ci](https://github.com/zup-ci))
- \[FIX\] update local repo name [\#142](https://github.com/ZupIT/ritchie-cli/pull/142) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))

## [1.0.0-beta.12](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.12) (2020-06-03)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.11...1.0.0-beta.12)

**Merged pull requests:**

- \[FIX\] changed go to compiled [\#143](https://github.com/ZupIT/ritchie-cli/pull/143) ([victor-schumacher](https://github.com/victor-schumacher))
- Release 1.0.0-beta.11 merge [\#139](https://github.com/ZupIT/ritchie-cli/pull/139) ([zup-ci](https://github.com/zup-ci))

## [1.0.0-beta.11](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.11) (2020-06-02)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.10...1.0.0-beta.11)

**Closed issues:**

- \[BUG\] Slash in new formula's command generate a path error [\#108](https://github.com/ZupIT/ritchie-cli/issues/108)
- \[FEATURE\] Update RIT CREATE FORMULA command [\#107](https://github.com/ZupIT/ritchie-cli/issues/107)

**Merged pull requests:**

- \[FIX\] env runner formula [\#138](https://github.com/ZupIT/ritchie-cli/pull/138) ([marcosgmgm](https://github.com/marcosgmgm))
- \[Fix\] Remove special char in prompt [\#137](https://github.com/ZupIT/ritchie-cli/pull/137) ([kaduartur](https://github.com/kaduartur))
- \[FEATURE\] update templates [\#136](https://github.com/ZupIT/ritchie-cli/pull/136) ([victor-schumacher](https://github.com/victor-schumacher))
- \[FIX\] Rename field qtd to qty [\#135](https://github.com/ZupIT/ritchie-cli/pull/135) ([kaduartur](https://github.com/kaduartur))
- \[fix\] change name Passphrase [\#134](https://github.com/ZupIT/ritchie-cli/pull/134) ([ernelio](https://github.com/ernelio))
- Feature/improvements security [\#133](https://github.com/ZupIT/ritchie-cli/pull/133) ([ernelio](https://github.com/ernelio))
- fixed golang creator [\#132](https://github.com/ZupIT/ritchie-cli/pull/132) ([victor-schumacher](https://github.com/victor-schumacher))
- Feature/nightly [\#127](https://github.com/ZupIT/ritchie-cli/pull/127) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- Feature/stdin functional tests [\#126](https://github.com/ZupIT/ritchie-cli/pull/126) ([dmbarra](https://github.com/dmbarra))
- \[FEATURE\] Improvements circleci [\#125](https://github.com/ZupIT/ritchie-cli/pull/125) ([ernelio](https://github.com/ernelio))
- \[FEATURE\] cmd init [\#123](https://github.com/ZupIT/ritchie-cli/pull/123) ([sandokandias](https://github.com/sandokandias))
- Release 1.0.0-beta.10 merge [\#122](https://github.com/ZupIT/ritchie-cli/pull/122) ([zup-ci](https://github.com/zup-ci))
- \[FEATURE\] Formula creator improvements [\#104](https://github.com/ZupIT/ritchie-cli/pull/104) ([victor-schumacher](https://github.com/victor-schumacher))

## [1.0.0-beta.10](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.10) (2020-05-27)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.9...1.0.0-beta.10)

**Merged pull requests:**

- \[FEATURE\] adding link to changelog inside the release description [\#121](https://github.com/ZupIT/ritchie-cli/pull/121) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- Release 1.0.0-beta.9 merge [\#120](https://github.com/ZupIT/ritchie-cli/pull/120) ([zup-ci](https://github.com/zup-ci))

## [1.0.0-beta.9](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.9) (2020-05-27)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.8...1.0.0-beta.9)

**Closed issues:**

- \[BUG\] Changelog generation not working properly [\#105](https://github.com/ZupIT/ritchie-cli/issues/105)

**Merged pull requests:**

- testing in production the old way [\#119](https://github.com/ZupIT/ritchie-cli/pull/119) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))

## [1.0.0-beta.8](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.8) (2020-05-27)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.7...1.0.0-beta.8)

**Closed issues:**

- \[BUG\] Getting error when running rit set server in Windows [\#83](https://github.com/ZupIT/ritchie-cli/issues/83)
- \[FEATURE\] Command for publishing a new formula [\#33](https://github.com/ZupIT/ritchie-cli/issues/33)

**Merged pull requests:**

- Revert "Improve Pipeline with security test" [\#115](https://github.com/ZupIT/ritchie-cli/pull/115) ([ernelio](https://github.com/ernelio))
- Improve Pipeline with security test [\#114](https://github.com/ZupIT/ritchie-cli/pull/114) ([flavioanellozup](https://github.com/flavioanellozup))
- \[FIX\] Links documentation after update beta-7 [\#112](https://github.com/ZupIT/ritchie-cli/pull/112) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))
- \[FIX\] bin name linux [\#111](https://github.com/ZupIT/ritchie-cli/pull/111) ([marcosgmgm](https://github.com/marcosgmgm))
- \[FEATURE\] Update slack token [\#109](https://github.com/ZupIT/ritchie-cli/pull/109) ([kaduartur](https://github.com/kaduartur))
- \[FIX\] changelog [\#106](https://github.com/ZupIT/ritchie-cli/pull/106) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- Feature/functional test team [\#102](https://github.com/ZupIT/ritchie-cli/pull/102) ([dmbarra](https://github.com/dmbarra))
- Release 1.0.0-beta.7 merge [\#101](https://github.com/ZupIT/ritchie-cli/pull/101) ([zup-ci](https://github.com/zup-ci))

## [1.0.0-beta.7](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.7) (2020-05-20)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.6...1.0.0-beta.7)

**Closed issues:**

- \[BUG\] the formula should know where the user call it. [\#99](https://github.com/ZupIT/ritchie-cli/issues/99)
- \[FEATURE\] STDIN Implementation on Ritchie-CLI [\#96](https://github.com/ZupIT/ritchie-cli/issues/96)
- \[FEATURE\] - Access control for formula [\#82](https://github.com/ZupIT/ritchie-cli/issues/82)
- \[FEATURE\] Global credential by organization [\#77](https://github.com/ZupIT/ritchie-cli/issues/77)
- \[BUG\] Requirements golint [\#42](https://github.com/ZupIT/ritchie-cli/issues/42)
- \[FEATURE\] Allow running any language [\#19](https://github.com/ZupIT/ritchie-cli/issues/19)
- \[FEATURE\] Add commons repository by default in Single version [\#16](https://github.com/ZupIT/ritchie-cli/issues/16)

**Merged pull requests:**

- \[FIX\] Failed formulas use pwd [\#100](https://github.com/ZupIT/ritchie-cli/pull/100) ([marcosgmgm](https://github.com/marcosgmgm))
- \[FEATURE\] STDIN Implementation [\#97](https://github.com/ZupIT/ritchie-cli/pull/97) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))
- Feature/access control for formulas [\#95](https://github.com/ZupIT/ritchie-cli/pull/95) ([marcosgmgm](https://github.com/marcosgmgm))
- add victor [\#94](https://github.com/ZupIT/ritchie-cli/pull/94) ([victor-schumacher](https://github.com/victor-schumacher))
- \[FEATURE\] Adding changelog to pipeline [\#93](https://github.com/ZupIT/ritchie-cli/pull/93) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- Feature/dockerfile create f [\#90](https://github.com/ZupIT/ritchie-cli/pull/90) ([victor-schumacher](https://github.com/victor-schumacher))
- Feature/functional tests [\#89](https://github.com/ZupIT/ritchie-cli/pull/89) ([dmbarra](https://github.com/dmbarra))
- Feature/ritchie-bo\_-config [\#85](https://github.com/ZupIT/ritchie-cli/pull/85) ([kaduartur](https://github.com/kaduartur))
- \[FIX\] Set server error [\#84](https://github.com/ZupIT/ritchie-cli/pull/84) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))
- \[FIX\] updating gitbook urls [\#81](https://github.com/ZupIT/ritchie-cli/pull/81) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))
- \[FEATURE\] Add credential for the organizations  [\#79](https://github.com/ZupIT/ritchie-cli/pull/79) ([kaduartur](https://github.com/kaduartur))
- change to . js [\#78](https://github.com/ZupIT/ritchie-cli/pull/78) ([victor-schumacher](https://github.com/victor-schumacher))
- \[FEATURE\] Move the fileutil to stream package [\#76](https://github.com/ZupIT/ritchie-cli/pull/76) ([kaduartur](https://github.com/kaduartur))
- \[FEATURE\] Enhancement/delete repo list [\#75](https://github.com/ZupIT/ritchie-cli/pull/75) ([victor-schumacher](https://github.com/victor-schumacher))

## [1.0.0-beta.6](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.6) (2020-05-06)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.5...1.0.0-beta.6)

**Merged pull requests:**

- Fix create formula Go [\#73](https://github.com/ZupIT/ritchie-cli/pull/73) ([ernelio](https://github.com/ernelio))

## [1.0.0-beta.5](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.5) (2020-05-06)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.4...1.0.0-beta.5)

**Merged pull requests:**

- \[FIX\] updating single repo commons url [\#71](https://github.com/ZupIT/ritchie-cli/pull/71) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))

## [1.0.0-beta.4](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.4) (2020-05-06)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.3...1.0.0-beta.4)

**Merged pull requests:**

- changing the name of the dist directory regarding macos [\#69](https://github.com/ZupIT/ritchie-cli/pull/69) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))

## [1.0.0-beta.3](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.3) (2020-05-06)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.2...1.0.0-beta.3)

**Closed issues:**

- \[FEATURE\] Create new languages in create formula. [\#51](https://github.com/ZupIT/ritchie-cli/issues/51)
- \[FEATURE\] Group commands by core/repo [\#43](https://github.com/ZupIT/ritchie-cli/issues/43)
- \[FEATURE\] Add golint in circleci [\#37](https://github.com/ZupIT/ritchie-cli/issues/37)
- \[FEATURE\] Command for set serverURL and remove from build [\#30](https://github.com/ZupIT/ritchie-cli/issues/30)
- \[FEATURE\] CircleCI pipeline [\#28](https://github.com/ZupIT/ritchie-cli/issues/28)
- Enhancement test for pkg/cmd [\#23](https://github.com/ZupIT/ritchie-cli/issues/23)

**Merged pull requests:**

- Feature/warning [\#67](https://github.com/ZupIT/ritchie-cli/pull/67) ([victor-schumacher](https://github.com/victor-schumacher))
- Fix/create formula python [\#66](https://github.com/ZupIT/ritchie-cli/pull/66) ([ernelio](https://github.com/ernelio))
- Feature/group commands [\#65](https://github.com/ZupIT/ritchie-cli/pull/65) ([sandokandias](https://github.com/sandokandias))
- Feature/default commons repo [\#64](https://github.com/ZupIT/ritchie-cli/pull/64) ([victor-schumacher](https://github.com/victor-schumacher))
- Feature/create formula any languages [\#63](https://github.com/ZupIT/ritchie-cli/pull/63) ([ernelio](https://github.com/ernelio))
- \[FEATURE\] Set server cmd update [\#62](https://github.com/ZupIT/ritchie-cli/pull/62) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))
- Feature/add bin other so tree [\#61](https://github.com/ZupIT/ritchie-cli/pull/61) ([ernelio](https://github.com/ernelio))
- Feature/requirements lint [\#58](https://github.com/ZupIT/ritchie-cli/pull/58) ([ernelio](https://github.com/ernelio))
- Feature/pr ci [\#54](https://github.com/ZupIT/ritchie-cli/pull/54) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- Feature/code of conduct [\#53](https://github.com/ZupIT/ritchie-cli/pull/53) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))
- adding stale bot to ritchie-cli [\#52](https://github.com/ZupIT/ritchie-cli/pull/52) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- Feature/any languages [\#48](https://github.com/ZupIT/ritchie-cli/pull/48) ([erneliojuniorzup](https://github.com/erneliojuniorzup))
- Update README.md [\#46](https://github.com/ZupIT/ritchie-cli/pull/46) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- applying standards [\#45](https://github.com/ZupIT/ritchie-cli/pull/45) ([sandokandias](https://github.com/sandokandias))
- Require lint [\#41](https://github.com/ZupIT/ritchie-cli/pull/41) ([ernelio](https://github.com/ernelio))
- changes [\#40](https://github.com/ZupIT/ritchie-cli/pull/40) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- changing badges [\#39](https://github.com/ZupIT/ritchie-cli/pull/39) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- updating ernelio github user [\#38](https://github.com/ZupIT/ritchie-cli/pull/38) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))
- Feature/lint [\#36](https://github.com/ZupIT/ritchie-cli/pull/36) ([erneliojuniorzup](https://github.com/erneliojuniorzup))
- \[FIX \] use $\(HOME\) instead of ~ in Makefile \(main\) [\#35](https://github.com/ZupIT/ritchie-cli/pull/35) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))
- \[FEATURE\] Creation set server cmd \(TEAM\) [\#34](https://github.com/ZupIT/ritchie-cli/pull/34) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))
- fix docs url [\#32](https://github.com/ZupIT/ritchie-cli/pull/32) ([rodrigomedeirosf](https://github.com/rodrigomedeirosf))
- Circleci project setup [\#29](https://github.com/ZupIT/ritchie-cli/pull/29) ([viniciusramosdefaria](https://github.com/viniciusramosdefaria))
- Feature/enhancement\_test\_pkg\_cmd [\#27](https://github.com/ZupIT/ritchie-cli/pull/27) ([sandokandias](https://github.com/sandokandias))
- change to apache licence [\#26](https://github.com/ZupIT/ritchie-cli/pull/26) ([rodrigomedeirosf](https://github.com/rodrigomedeirosf))
- \[FIX\] Removing templates [\#25](https://github.com/ZupIT/ritchie-cli/pull/25) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))
- Fix hello world in tpl\_main [\#24](https://github.com/ZupIT/ritchie-cli/pull/24) ([erneliojuniorzup](https://github.com/erneliojuniorzup))
- Feature/improve\_session\_validator [\#21](https://github.com/ZupIT/ritchie-cli/pull/21) ([kaduartur](https://github.com/kaduartur))
- \[FEATURE\] updating issues contribution templates [\#20](https://github.com/ZupIT/ritchie-cli/pull/20) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))
- \[Enhancement\] Contributing file [\#17](https://github.com/ZupIT/ritchie-cli/pull/17) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))
- Release 1.0.0-beta.2 merge [\#15](https://github.com/ZupIT/ritchie-cli/pull/15) ([zup-ci](https://github.com/zup-ci))

## [1.0.0-beta.2](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.2) (2020-04-09)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta.1...1.0.0-beta.2)

**Closed issues:**

- The autocomplete to create formula doesn't work [\#12](https://github.com/ZupIT/ritchie-cli/issues/12)

**Merged pull requests:**

- feature/prompt\_interface [\#14](https://github.com/ZupIT/ritchie-cli/pull/14) ([sandokandias](https://github.com/sandokandias))
- fix/create\_formula\_autocomplete [\#13](https://github.com/ZupIT/ritchie-cli/pull/13) ([kaduartur](https://github.com/kaduartur))
- \[Enhancement\] Contributing file [\#11](https://github.com/ZupIT/ritchie-cli/pull/11) ([GuillaumeFalourd](https://github.com/GuillaumeFalourd))
- Release 1.0.0-beta.1 merge [\#10](https://github.com/ZupIT/ritchie-cli/pull/10) ([zup-ci](https://github.com/zup-ci))

## [1.0.0-beta.1](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta.1) (2020-04-09)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/1.0.0-beta...1.0.0-beta.1)

**Closed issues:**

- The server url is showing while running any commands [\#7](https://github.com/ZupIT/ritchie-cli/issues/7)

**Merged pull requests:**

- fix issue \#7 [\#8](https://github.com/ZupIT/ritchie-cli/pull/8) ([sandokandias](https://github.com/sandokandias))
- Release 1.0.0-beta merge [\#6](https://github.com/ZupIT/ritchie-cli/pull/6) ([zup-ci](https://github.com/zup-ci))

## [1.0.0-beta](https://github.com/zupit/ritchie-cli/tree/1.0.0-beta) (2020-04-09)

[Full Changelog](https://github.com/zupit/ritchie-cli/compare/da1809eba73786a35a6c211d937e2d27b08d6361...1.0.0-beta)

**Merged pull requests:**

- fix build [\#5](https://github.com/ZupIT/ritchie-cli/pull/5) ([marcosgmgm](https://github.com/marcosgmgm))
- Create CODEOWNERS [\#4](https://github.com/ZupIT/ritchie-cli/pull/4) ([marcosgmgm](https://github.com/marcosgmgm))
- Fix create formula and binaries [\#3](https://github.com/ZupIT/ritchie-cli/pull/3) ([sandokandias](https://github.com/sandokandias))
- Fix load repo in root [\#2](https://github.com/ZupIT/ritchie-cli/pull/2) ([kaduartur](https://github.com/kaduartur))
- fix/load\_tree [\#1](https://github.com/ZupIT/ritchie-cli/pull/1) ([kaduartur](https://github.com/kaduartur))



\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
