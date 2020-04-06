# datadog_sample

The sample of Datadog APM.

<img width="650" alt="datadog_sample" src="https://user-images.githubusercontent.com/301822/78652870-451fea80-78fd-11ea-821b-57e594ce2b1a.png">

Trace Result in APM

![FlameGraph_grpc_exp](https://user-images.githubusercontent.com/301822/78885124-01f97f00-7a97-11ea-93a2-12b7e21edc56.png)

## Requirements

### Datadog Agent

[Install](https://app.datadoghq.com/account/login?next=%2Faccount%2Fsettings#agent)

### DB(mysql)

```
CREATE DATABASE `datadog_sample`;

USE `datadog_sample`;

CREATE TABLE `message` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `text` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT 'text',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='message'
```

## config file

Configuration information such as gRPC Server address and DB address is in the config file.
please edit [config.yaml](https://github.com/muroon/datadog_sample/blob/master/config/config.yaml)

## gRPC Server

```
go run grpcserver/main.go
```

## HTTP Server

```
go run httpserver/main.go
```

## Access

```
curl http://localhost:8080/grpc/
```

