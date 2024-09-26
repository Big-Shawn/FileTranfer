# FileTransfer


### 文件传输工具

- 基于`TCP`协议，传输可靠稳定。
- 使用`protobuf`作为传输帧协议，封装方便，使用简单，数据负载占帧负载约为99%，传输高效。
- 支持设置文件夹批量传输，支持大部分文件格式。
- C/S 架构，服务端可配置，客户端接入后自动开始文件传输。
- 传输速度快，实际测试可跑满路由器上传/下载峰值速度。

### 使用方法

#### 服务端
`./server config.yaml`

基于配置文件 `config.yaml`, 示例如下
```yaml
# zap 日志设置
logger:
  level: -1
  development: true
  encoding : "console"
  path:
    - "stdout"
    - "srv.log"

port: ":8879"  #监听端口
package_size: 2048 #传输帧大小
path:  #共享路径
  - "file.txt"
  - "/path/to/file"
  - "/path/will/be/work/too"
```

#### 客户端

`./client -p serverhost -c file2store`

- serverhost: 服务端地址，如：127.0.0.1:8879
- file2store: 文件存放地址，如：/file/will/be/stored

