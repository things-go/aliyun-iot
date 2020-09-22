## aliyun 物联网设备接入的golang实现 
WIP

- [x] infra 公共包
- [x] sign: 实现MQTT签名,独立使用,不依赖第三方任何包
- [x] dynamic: 直连设备动态注册
- [x] ahttp: http 上云实现
- [ ] dm: 物联型mqtt上云实现,独立使用
- [ ] dataflow: 服务器订阅数据流定义


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
    - [ ] extend RRPC

- gateway
    - [ ] event property pack post
    - [ ] event property history post
