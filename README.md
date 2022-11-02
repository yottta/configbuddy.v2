[![CircleCI](https://circleci.com/gh/yottta/configbuddy.v2.svg?style=shield)](https://circleci.com/gh/yottta/configbuddy.v2)
[![Coverage Status](https://coveralls.io/repos/github/yottta/configbuddy.v2/badge.svg?branch=master)](https://coveralls.io/github/yottta/configbuddy.v2?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/yottta/configbuddy.v2)](https://goreportcard.com/report/github.com/yottta/configbuddy.v2)

# configbuddy v2

A tool initially designed to help with installing dotfiles but can link others too.

## Installing
### Artifact
You can find the available versions in the [releases page](https://github.com/yottta/configbuddy.v2/releases).
### Go
This can be also installed using go by just running:
```shell
go install github.com/yottta/configbuddy.v2
```
> Tip: In order to be sure that you can run the app right after running install, be sure that you have `go env GOBIN` added to your system PATH.
> ```shell
> export PATH="`go env GOBIN`;$PATH"
> ```

### Integration
There is pretty much no moment when you are going to have Go installed in order to use the Go installation method, therefore for simplicity, 
you can create a script in your dotfiles repository with the following content:
```shell
#!/bin/bash

current_system=$(uname -s)
configbuddy_url=$(curl -s  -i https://api.github.com/repos/yottta/configbuddy.v2/releases/latest | grep download_url | grep tar.gz | grep -i $current_system | sed 's/": /#/' | cut -d'#' -f2 | sed 's/"//g')

cd $HOME
mkdir __delete_me_configbuddy
cd __delete_me_configbuddy

wget $configbuddy_url -O configbuddy.tar.gz
tar -xf configbuddy.tar.gz

mv configbuddy.v2 $HOME
cd $HOME
```

And after the script is executed you will have configbuddy installed on your home directory.
Now you can just run:
```shell
~/configbuddy.v2 -c $HOME/.dotfiles/configs/config_file.yml -l debug -b
```

## Features
* For installing configuration files it gives a lot of flexibility on how the installation should be done
* Setup configurations in multiple files and to add imports in each other config file of other config file  
    ```yaml
    includes:
      - other_config_file.yml
    ```
* Support for the most used env vars when we are talking about dotfiles
  * To indicate that an env var is wanted to be used, in the config files the env vars needs to be surrounded by `$#` and `#$`
  * The supported env vars are the following (TODO: refactor to make use of `os.Getenv`):
    * HOME - The home directory of the user that it's running the app
    * USER - The name of the user that it's running the app
    * DISTRO - To be implemented
    * PCK_MANAGER - To be implemented
  
## Configuration

```yaml
Globals: # global configuration
  exitOnError: true # configure the app to exit on the first encountered error
  confirmEveryPackage: true # when it's installing packages, ask for confirmation. not used yet

includes: # define a list of other configuration files
  - tmux.yml 

FileAction: # file actions
  aliases: # if the `name` field is not configured, the target file name will be this key
    name: not_default_name_aliases # the target file name
    hidden: true # as this is true, the file name will have `.` in front
    source: ../home # the place where the file can be found
    command: ln -s # the command that should be used for configuring this aliases file
    destination: $#HOME#$ # the destination directory where the file will be configured
    # in this case, the file <current directory>/../home/aliases will be symlinked to $HOME/.not_default_name_aliases
    
  bashrc:
    hidden: true
    source: ../home
    command: ln -s
    destination: $#HOME#$
    # in this case, the file from the <current directory>/../home/bashrc will be symlinked to $HOME/.bashrc
    
  exports:
    hidden: false
    source: ../home
    command: ln -s
    destination: $#HOME#$
    # in this case, the file from the <current directory>/../home/exports will be symlinked to $HOME/exports
```


## Future features
* Install apps from package managers
  * Option for alternatives in case the main package name does not exit 
* Install packages from URLs (zip, tar.gz, deb, dpk, rpm)
