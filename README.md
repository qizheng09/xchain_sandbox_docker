# xchain_sandbox_docker
A multi xchain node sandbox running in docker

## Usage

### Requirements

* OS Support: Linux, Mac OS and Windows 
  (You need to build XuperUnion and xchain_sandbox_docker by docker contianer if you use MAC OS or Windows)
* Go 1.12.x or later
* Git

### Build

### Clone Repository

``` 
git clone https://github.com/qizheng09/xchain_sandbox_docker.git
```

#### Set Env

``` 
export XCHAIN_ROOT={{path to XuperUnion directory}}
export XCHAIN_SAND_ROOT={{path to xchain_sandbox_docker directory}}
```

#### Build Tool

``` 
cd sandbox && go build
```

#### Init Sandbox

Init the sandbox enviroment, including binaries, nodes, leadger and so on.

```
./sandbox init -N {{nodes number}} -M {{miner number}}
```
#### Update Sandbox

This used in update sandbox enviroment, only including binaries.

```
./sandbox update
```

#### Start Sandbox

Start sandbox enviroment by docker containers.

```
./sandbox start 
```

#### Stop Sandbox

Stop sandbox enviroment.

```
./sandbox stop 
```

#### Clear Sandbox

Clear the sandbox enviroment

```
./sandbox clear 
```

#### For more useage

```
./sandbox --help
```