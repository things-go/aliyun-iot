## aliyun 物联网设备接入的golang实现 (WIP)

[![GoDoc](https://godoc.org/github.com/thinkgos/aliyun-iot?status.svg)](https://godoc.org/github.com/thinkgos/aliyun-iot)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/thinkgos/aliyun-iot?tab=doc)
[![Build Status](https://www.travis-ci.org/thinkgos/aliyun-iot.svg?branch=master)](https://www.travis-ci.org/thinkgos/aliyun-iot)
[![codecov](https://codecov.io/gh/thinkgos/aliyun-iot/branch/master/graph/badge.svg)](https://codecov.io/gh/thinkgos/aliyun-iot)
![Action Status](https://github.com/thinkgos/aliyun-iot/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/thinkgos/aliyun-iot)](https://goreportcard.com/report/github.com/thinkgos/aliyun-iot)
[![Licence](https://img.shields.io/github/license/thinkgos/aliyun-iot)](https://raw.githubusercontent.com/thinkgos/aliyun-iot/master/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/thinkgos/aliyun-iot)](https://github.com/thinkgos/aliyun-iot/tags)

- [x] infra 公共包
- [x] sign: 实现MQTT签名,独立使用,不依赖第三方任何包
- [x] dynamic: 直连设备动态注册
- [x] ahttp: http 上云实现
- [x] dataflow: 服务器订阅数据流定义


## Feature 

- device
    - [x] raw up and raw up reply
    - [x] raw down
    - [x] event property post and reply
    - [x] event post and reply
    - [x] ntp
    - [x] config get and push
    - [x] label update and delete
    - [x] RRPC
    - [x] extend RRPC

- gateway
    - [x] event property pack post
    - [x] event property history post

## Sponsor
**alipay**

![alipay](https://raw.githubusercontent.com/thinkgos/thinkgos/master/asserts/alipay.jpg)

**wxpay**

![wxpay](https://raw.githubusercontent.com/thinkgos/thinkgos/master/asserts/wxpay.jpg)
