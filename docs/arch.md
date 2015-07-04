## 文档

分成三个组件

- 核心组件，用来处理各种业务逻辑（开会提醒，提供信息等）
- 接入组件，实现核心组件提供的接口，把各个接入端（QQ，Slack，Email等）提供的信息转换成核心组件识别的结构
- 管理员界面，用于提交各个业务逻辑所需要的基本信息（也可以认为是核心组件的一部分


### 如何接入

```
type Postman interface {
    Type() string
    SendMsg(to string, msg string) error
    RecvMsg() chan Msg
    SendGroupMsg(gid string, msg string) error
    RecvGroupMsg() chan GroupMsg
}
```

所有的接入组件（邮差）只要实现上述功能就可以接入QBot机器人
