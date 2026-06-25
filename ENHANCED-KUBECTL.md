# Enhanced kubectl

基于 Kubernetes master 分支，包含以下 kubectl 增强改动。

## 增强特性

### 1. kubectl logs -c 支持 container name 补全

`kubectl logs -c` 的 `--container` flag 注册了 shell completion 函数，输入 `-c <Tab>` 时会自动列出 Pod 中的容器名。

- 文件: `staging/src/k8s.io/kubectl/pkg/cmd/logs/logs.go`
- 改动: 调用 `cmd.RegisterFlagCompletionFunc("container", completion.ContainerCompletionFunc(f))`

### 2. container name 补全从 containerStatuses 获取

container name 补全数据源从 `spec.containers` 改为 `status.containerStatuses`。

原因: EKS DaemonSet 注入场景下，注入的 container 可能只存在于 `status.containerStatuses` 而不在 `spec.containers` 中，从 spec 取会漏掉注入的 container。

- 文件: `staging/src/k8s.io/kubectl/pkg/util/completion/completion.go`
- 改动: template 从 `.spec.containers` 改为 `.status.containerStatuses`

### 3. 删除 kubectl logs 的 container 名称校验

`kubectl logs` 不再校验 container 名称是否存在于 pod spec 中。

原因: 超级节点 DaemonSet 注入场景下，注入的 container 不在 `spec.containers` 中，原校验逻辑会报 "container xxx is not valid for pod" 导致无法查看注入 container 的日志。

- 文件: `staging/src/k8s.io/kubectl/pkg/polymorphichelpers/logsforobject.go`
- 改动: 移除 `FindContainerByName` 返回 nil 时的报错

### 4. 删除 kubectl exec 的 container 名称校验

`kubectl exec -c` 指定 container name 时不再校验 container 是否存在于 pod spec 中。

原因: 与 logs 相同，超级节点注入的隐藏 container 只在 `status.containerStatuses` 中，不在 `spec.containers` 中，校验会导致无法 exec 到注入的 container。

- 文件: `staging/src/k8s.io/kubectl/pkg/cmd/exec/exec.go`
- 改动: 移除指定 container name 时的 `FindContainerByName` 校验

## Build

### 前置要求

- Go (参考 `hack/lib/golang.sh` 中的版本要求)
- 本地构建需要 Go 环境；Docker 构建需要 Docker

### 编译当前平台的 kubectl

```bash
make kubectl
```

输出路径: `_output/local/bin/<host_os>/<host_arch>/kubectl`

host platform 还会在 `_output/local/bin/kubectl` 建符号链接。

### 编译指定 os + arch 的 kubectl

```bash
KUBE_BUILD_PLATFORMS=linux/amd64 make kubectl
```

支持的平台:

| OS      | Arch    |
| ------- | ------- |
| linux   | amd64   |
| linux   | 386     |
| linux   | arm     |
| linux   | arm64   |
| linux   | s390x   |
| linux   | ppc64le |
| darwin  | amd64   |
| darwin  | arm64   |
| windows | amd64   |
| windows | 386     |
| windows | arm64   |

输出路径: `_output/local/bin/<os>/<arch>/kubectl` (windows 下为 `kubectl.exe`)

### 编译多种架构的 kubectl

```bash
# 编译所有支持平台的 client 二进制（包括 kubectl）
make kubectl KUBE_BUILD_PLATFORMS="linux/amd64 linux/arm64 darwin/arm64"
```

每个平台的二进制输出到 `_output/local/bin/<os>/<arch>/kubectl`。

### Docker 容器化构建

如果本地 Go 环境不便配置，可以用 Docker 容器构建:

```bash
# 编译当前平台
build/run.sh make kubectl

# 编译指定平台
build/run.sh make kubectl KUBE_BUILD_PLATFORMS=darwin/arm64
```

Docker 构建输出路径: `_output/bin/<os>/<arch>/kubectl`

### 快速构建（仅 host 架构）

```bash
KUBE_FASTBUILD=true make kubectl
```
