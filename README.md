[![Go Report Card](https://goreportcard.com/badge/github.com/cloudconsole/cloudspider)](https://goreportcard.com/report/github.com/cloudconsole/cloudspider) [![CircleCI](https://circleci.com/gh/cloudconsole/cloudspider/tree/master.svg?style=svg)](https://circleci.com/gh/cloudconsole/cloudspider/tree/master)

# Cloud Spider
Yet another dashboard and a crawler to visualize/integrate all your Cloud/IT infrastructure in a unified view.

It's in a very very early stage.

### What cloud services are supported so far
We are working hard to add more cloud services and soon you will see your favorite cloud in the supported list.

1. AWS
2. UltraDNS
3. Akamai
4. DnsMadeEasy

### Prerequisites
1. MongoDB installed
2. All your cloud provider credentials

### Install

```
go get github.com/cloudconsole/cloudspider
```

#### Configuring credentials
```
mkdir .cloudspider
cp $GOPATH/src/github.com/cloudconsole/cloudspider/conf/sample_config.yaml ~/.cloudspider/config.yaml
```
Modify the file `~/.cloudspider/config.yaml` to fill in your credentials or other application related settings. File is self documented and eazy to understand.

#### Usage
```
shell$ cloudspider --help
A Cloud Crawler process to crawl all your IT infraStructure and
stors then in to a MongoDB util.

Usage:
  cloudcrawler [command]

Available Commands:
  crawler     Unleash the spiders
  help        Help about any command
  ui          A high performance webserver
  version     Print the program version

Flags:
      --config string   config file (default is $HOME/.cloudspider/config.yaml)

Use "cloudcrawler [command] --help" for more information about a command.
```

To start a UI server

```
cloudspider ui start
```

To start the crawler to crawl all your cloud providers and index

```
cloudspider crawler start
```
