# ContainDB
24 files | 2272 logic lines

## Tree
```
ContainDB/
├─ npm
│  └─ InstallController.js
└─ src
   ├─ base
   │  ├─ Banner.go
   │  ├─ BaseCaseHandler.go
   │  ├─ DatabaseSelector.go
   │  ├─ DockerStarterPack.go
   │  ├─ FilePathSelector.go
   │  ├─ flagHandler.go
   │  └─ StartContainer.go
   ├─ Core
   │  └─ main.go
   ├─ Docker
   │  ├─ docker_container.go
   │  ├─ docker_installation.go
   │  ├─ Docker_Network.go
   │  ├─ docker.go
   │  ├─ DockerComposeMaker.go
   │  ├─ ImportDockerServices.go
   │  ├─ platform.go
   │  └─ SysRequirement.go
   └─ tools
      ├─ AfterContainerToolInstaller.go
      ├─ askForInput.go
      ├─ MongoDB_Tools.go
      ├─ PgAdmin.go
      ├─ PhpMyAdmin.go
      ├─ Redis_Insight.go
      └─ rollback.go
```

## Modules

### npm
**InstallController.js** (29L)
deps:child_process,path,os

### src/Core
**main.go** (76L)
pkg:main | uses:Docker,base,tools,fmt,os,signal,runtime
  main() @13

### src/Docker
**DockerComposeMaker.go** (275L)
pkg:Docker | uses:fmt,os,exec,filepath,strings,template
structs: ContainerInfo
  MakeDockerComposeWithAllServices() -> string [E] @26
  getContainerInfo(containerName) -> (ContainerInfo, error) @82
  generateComposeYAML(containers) -> string @248
**Docker_Network.go** (16L)
pkg:Docker | uses:os,exec
  CreateDockerNetworkIfNotExists() -> error [E] @8
**ImportDockerServices.go** (103L)
pkg:Docker | uses:errors,fmt,ioutil,net,os,exec,strings,yaml.v2
structs: DockerComposeConfig, DockerComposeService
  ImportDockerServices(composeFilePath) -> error [E] @31
  volumeExists(name) -> bool @108
  getRunningContainers() -> ([]string, error) @115
  isPortAvailable(port) -> bool @125
**SysRequirement.go** (147L)
pkg:Docker | uses:fmt,os,exec,runtime,strconv,strings
  CheckSystemRequirements() [E] @12
  checkDockerInstallation() -> error @35
  checkRAM(minGB) -> error @46
  checkDiskSpace(minGB) -> error @115
**docker.go** (139L)
pkg:Docker | uses:fmt,os,exec,strings
  ListRunningDatabases() -> ([]string, error) [E] @11
  RemoveDatabase(name) -> error [E] @28
  ListDatabaseImages() -> ([]string, error) [E] @68
  IsImageInUse(image) -> (bool, string, error) [E] @103
  RemoveImage(image) -> error [E] @123
  ListContainDBVolumes() -> ([]string, error) [E] @134
  IsVolumeInUse(volume) -> (bool, string, error) [E] @168
**docker_container.go** (92L)
pkg:Docker | uses:fmt,os,exec,strings,promptui
  AskYesNo(label) -> bool [E] @12
  IsContainerRunning(nameOrImage, checkByName) -> bool [E] @31
  ListOfContainers(images) -> []string [E] @42
  VolumeExists(name) -> bool [E] @83
  CreateVolume(name) -> error [E] @90
  RemoveVolume(name) -> error [E] @98
**docker_installation.go** (161L)
pkg:Docker | uses:fmt,os,exec,runtime
  IsDockerInstalled() -> bool [E] @10
  InstallDocker() -> error [E] @16
  installDockerLinux() -> error @31
  installDockerWindows() -> error @59
  installDockerMacOS() -> error @76
  UninstallDocker() -> error [E] @120
  uninstallDockerLinux() -> error @135
  uninstallDockerWindows() -> error @158
  uninstallDockerMacOS() -> error @172
**platform.go** (91L)
pkg:Docker | uses:fmt,os,exec,runtime
  IsAdmin() -> bool [E] @12
  GetTempDir() -> string [E] @25
  IsWindows() -> bool [E] @33
  IsMacOS() -> bool [E] @38
  IsLinux() -> bool [E] @43
  GetOSName() -> string [E] @48
  CheckOSSupport() -> error [E] @62
  GetOSRelease() -> string [E] @75
  CheckDockerCommand(cmd) -> error [E] @88
  GetShell() -> string [E] @98
  ExecuteCommand(name, args) -> *exec.Cmd [E] @114
  BuildDockerRunCommand(args) -> []string [E] @119

### src/base
**Banner.go** (51L)
pkg:base | uses:fmt,os,runtime,strings,color
  ShowBanner() [E] @14
**BaseCaseHandler.go** (230L)
pkg:base | uses:Docker,tools,fmt,os,exec,runtime,promptui
  BaseCaseHandler() [E] @14
**DatabaseSelector.go** (24L)
pkg:base | uses:tools,fmt,os,promptui
  SelectDatabase() -> string [E] @11
**DockerStarterPack.go** (28L)
pkg:base | uses:Docker,fmt,os,promptui
  DockerStarter() [E] @11
**FilePathSelector.go** (96L)
pkg:base | uses:fmt,os,filepath,strings,promptui
  SelectFilePath(label, defaultPath, extension) -> (string, error) [E] @13
**StartContainer.go** (146L)
pkg:base | uses:Docker,tools,fmt,os,exec,strings,promptui
  StartContainer(database) [E] @14
**flagHandler.go** (50L)
pkg:base | uses:Docker,fmt,os
  FlagHandler() [E] @10

### src/tools
**AfterContainerToolInstaller.go** (47L)
pkg:tools | uses:Docker,fmt
  AfterContainerToolInstaller(database) [E] @22
**MongoDB_Tools.go** (70L)
pkg:tools | uses:Docker,fmt,io,http,os,exec,filepath,runtime
  DownloadMongoDBCompass() [E] @14
**PgAdmin.go** (93L)
pkg:tools | uses:Docker,fmt,os,exec,strings,promptui
  StartPgAdmin() [E] @13
**PhpMyAdmin.go** (185L)
pkg:tools | uses:Docker,fmt,os,exec,strings,promptui
structs: CloudDBConfig
  StartPHPMyAdmin() [E] @23
  selectConnectionType(hasLocalContainers) -> string @62
  startPHPMyAdminLocal(sqlContainers) @91
  startPHPMyAdminCloud() @139
  getCloudConnectionConfig() -> CloudDBConfig @209
**Redis_Insight.go** (72L)
pkg:tools | uses:Docker,fmt,os,exec,strings,promptui
  StartRedisInsight() [E] @13
**askForInput.go** (17L)
pkg:tools | uses:bufio,fmt,os,strings
  AskForInput(label, defaultValue) -> string [E] @10
**rollback.go** (34L)
pkg:tools | uses:Docker,fmt,os,exec,filepath,strings
  Cleanup() [E] @13

