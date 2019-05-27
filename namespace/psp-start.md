# summary

PSP, Pod Security Policy, 是用于检查 Pod 安全的对象。他可以限制 Pod 是否可以使用特权模式，挂载主机目录等等。

## 限制范围

||||
|:-:|:-:|:-:|
|是否特权模式|Running of privileged containers|privileged|
|是否root namespace|Usage of the root namespaces|hostPID, hostIPC|
|是否主机网络模式|Usage of host networking and ports|hostNetwork, hostPorts|
|可以选择存储类型|Usage of volume types|volumes|
|可以挂载主机哪些目录|Usage of the host filesystem|allowedHostPaths|
|lvm?|White list of FlexVolume drivers|allowedFlexVolumes|
||Allocating an FSGroup that owns the pod’s volumes|fsGroup|
|read only root file|Requiring the use of a read only root file system|readOnlyRootFilesystem|
|user in ctr|The user and group IDs of the container|runAsUser, supplementalGroups|
||Restricting escalation to root privileges|allowPrivilegeEscalation, defaultAllowPrivilegeEscalation|
||Linux capabilities|defaultAddCapabilities, requiredDropCapabilities, allowedCapabilities|
|SELinux|The SELinux context of the container|seLinux|
||The AppArmor profile used by containers|annotations|
||The seccomp profile used by containers|annotations|
||The sysctl profile used by containers|annotations|

## 策略实例

<span id="privileged" ></span>
最宽松策略

```yaml
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: privileged
  annotations:
    seccomp.security.alpha.kubernetes.io/allowedProfileNames: '*'
spec:
  privileged: true
  allowPrivilegeEscalation: true
  allowedCapabilities:
  - '*'
  volumes:
  - '*'
  hostNetwork: true
  hostPorts:
  - min: 0
    max: 65535
  hostIPC: true
  hostPID: true
  runAsUser:
    rule: 'RunAsAny'
  seLinux:
    rule: 'RunAsAny'
  supplementalGroups:
    rule: 'RunAsAny'
  fsGroup:
    rule: 'RunAsAny'
```

最严格的限制
```yaml
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: restricted
  annotations:
    seccomp.security.alpha.kubernetes.io/allowedProfileNames: 'docker/default'
    apparmor.security.beta.kubernetes.io/allowedProfileNames: 'runtime/default'
    seccomp.security.alpha.kubernetes.io/defaultProfileName:  'docker/default'
    apparmor.security.beta.kubernetes.io/defaultProfileName:  'runtime/default'
spec:
  privileged: false
  # Required to prevent escalations to root.
  allowPrivilegeEscalation: false
  # This is redundant with non-root + disallow privilege escalation,
  # but we can provide it for defense in depth.
  requiredDropCapabilities:
    - ALL
  # Allow core volume types.
  volumes:
    - 'configMap'
    - 'emptyDir'
    - 'projected'
    - 'secret'
    - 'downwardAPI'
    # Assume that persistentVolumes set up by the cluster admin are safe to use.
    - 'persistentVolumeClaim'
  hostNetwork: false
  hostIPC: false
  hostPID: false
  runAsUser:
    # Require the container to run without root privileges.
    rule: 'MustRunAsNonRoot'
  seLinux:
    # This policy assumes the nodes are using AppArmor rather than SELinux.
    rule: 'RunAsAny'
  supplementalGroups:
    rule: 'MustRunAs'
    ranges:
      # Forbid adding the root group.
      - min: 1
        max: 65535
  fsGroup:
    rule: 'MustRunAs'
    ranges:
      # Forbid adding the root group.
      - min: 1
        max: 65535
  readOnlyRootFilesystem: false
```

## 开启 PSP 功能

管理员需要关心这一章节。普通用户可以跳过。

Pod 安全检查的功能已经编译在 `kube-apiserver` 中。当 api server 接收到创建 Pod 的请求，他会检查 Pod 的各个参数，匹配创建者所能使用的策略。如果没有策略可以匹配，则 Pod 不会被创建。
只需要在 `kube-apiserver` 启动参数 `--admission-control=` 的值列表中加入 `PodSecurityPolicy` 就会开启检查。

从 1.10 开始，当 `kube-apiserver` 使用 Static Pod 的方式启动时会有一些问题。`kube-apiserver` 启动参数加入 `PodSecurityPolicy` 之后，他本身的 Pod 会因为没有匹配的策略而无法启动。
启动 `kube-apiserver` 的创建者是 `group` `system:node`。所以需要提前给他绑定绑定一个拥有 PSP 策略的角色。

0. 创建规则

    创建上一章节中的 [`最宽松策略`](#privileged)。或者其他你觉得合适的规则。

<span id="privileged-user"></span>
1. 创建角色。

```yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: privileged-user
rules:
- apiGroups: ["extensions"]
  resourceNames:
  - privileged # <-- 上文中 `最宽松策略`
  resources: ["podsecuritypolicies"]
  verbs:
  - use
```

2. 绑定账户和角色

```yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: privileged-users
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: privileged-user # <-- 上一步中创建的角色
subjects:
- kind: Group
  apiGroup: rbac.authorization.k8s.io
  name: system:nodes # <-- 创建 kube-apiserver 的账户
```

3. 修改 `kube-apiserver` 参数

- `cp /etc/kubernetes/manifests/kube-apiserver.yaml ./`
- `vim kube-apiserver.yaml`
- 在 `--admission-control=` 末尾加 `,PodSecurityPolicy`
- `cp kube-apiserver.yaml /etc/kubernetes/manifests/kube-apiserver.yaml`

4. 确保 `kube-apiserver` 启动

- `kubectl -n kube-system get pods | grep api` 能看到 pod 。
- 如果看不到，请从头看一遍文章。
- 如果还不行请参考 k8s 最新文档。

## 使用

- 用户创建 Pod
- Service Account 创建 Pod

### 用户创建 Pod

通常我们并不直接创建 Pod。这里只是为了理解各对象之间的关系。
可以快速浏览至下一小节[Service Account 创建 Pod]

确定使用哪个用户

`kubectl config view`

可以看到当前命令行的配置。参考 `current-context`, `context` 可以知道默认用的那个用户。
这里如下用户作为例子。

```yaml
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: REDACTED
    server: https://*.*.*.*:*
  name: kubernetes
contexts:
- context: 
    cluster: kubernetes
    user: kubernetes-admin # <--- 3. 具体使用的用户
  name: kubernetes-admin@kubernetes # <--- 2. 
current-context: kubernetes-admin@kubernetes # <--- 1. current-context
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: REDACTED
    client-key-data: REDACTED
```

这里看到 kubectl 使用的就是 user kubernetes-admin@kubernetes。我们假设 `kubernetes-admin@kubernetes` 是个普通用户，目前不能创建 Pod 。

1. 创建规则

这里说的规则就是 `PodSecutityPolicy` 对象本身，也可以说是 psp 资源。psp 资源是 cluster scoped，不分 namespace。

这里使用前文中的 [`最宽松策略`](#privileged)。也可以使用其他你觉得合适的规则。如果已经创建，则不需要重复创建。

`kubectl get psp` 查看已有 psp 资源。

最宽松策略

```yaml
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: privileged
  annotations:
    seccomp.security.alpha.kubernetes.io/allowedProfileNames: '*'
spec:
  privileged: true
  allowPrivilegeEscalation: true
  allowedCapabilities:
  - '*'
  volumes:
  - '*'
  hostNetwork: true
  hostPorts:
  - min: 0
    max: 65535
  hostIPC: true
  hostPID: true
  runAsUser:
    rule: 'RunAsAny'
  seLinux:
    rule: 'RunAsAny'
  supplementalGroups:
    rule: 'RunAsAny'
  fsGroup:
    rule: 'RunAsAny'
```

2. 创建角色

在 default namespace 下创建一个角色。

```yaml
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: default
  name: privileged-user
rules:
- apiGroups: ["extensions"]
  resourceNames:
  - privileged # <--- psp 资源
  resources: ["podsecuritypolicies"]
  verbs:
  - use
```

3. 绑定角色和用户

```yaml
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: a-privileged-user
  namespace: default
subjects:
- kind: User
  name: kubernetes-admin@kubernetes # <--- 我们假设的用户
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: privileged-user # <--- 要绑的角色
  apiGroup: rbac.authorization.k8s.io
```

4. 创建 Pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: nginx
    image: hub.xfyun.cn/devops/alpine-curl:3.6
    command: ["/bin/sh", "-c", "sleep 1000000"]
    ports:
    - containerPort: 80
```

### Service Account 创建 Pod

当用命令行直接创建 Pod 时(`kubectl apply -f pod.yaml`)。通常都是用某个 `User` 作为创建者。
当通过 Deployment，StatefulSet 这类对象创建 Pod 时。Pod 的创建者是 Service Account 类型的。
当没有明确指定 Service Account 时，将会使用所在 namespace 下名为 default 的 Service Account。所以只要给 default Service Account 绑定 privileged-user 就可以了。
当在 Pod 的 spec 中明确定义了所使用的 `serviceAccountName` ，那么将会使用指定的。
所以只要给对应的 Service Account 绑定 privileged-user 就可以了。可以参考如下绑定。

绑定 default Service Account

```yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: privileged-user
subjects:
- kind: ServiceAccount
  name: default # <--- 或者其他 sa
  namespace: default
```

以上是 psp 的基本使用。下边是扩展技巧。

### 绑定角色使用 `ClusterRole`

`Role` 资源是 namespace scoped。仅限所在的 namespace 使用。当在不同 namespace 中使用同一个 psp 策略时，需要创建不同的 `Role`。这个时候可以使用 `ClusterRole`。

`ClusterRole` 是 cluster scoped 资源。当用 `ClusterRoleBinding` 绑定 `ClusterRole` 和账户时。此账户就在集群范围内拥有了 `ClusterRole` 定义的权限。

当在 `RoleBinding` 中绑定 `ClusterRole` 和账户时。此账户就只在 `RoleBinding` 所在 namespace 中拥有 `ClusterRole` 定义的权限。
这样就可以给不同账户绑定同一个角色。

```yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding # <--- RoleBinding
metadata:
  name: default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole # <--- ClusterRole
  name: privileged-user
subjects:
- kind: ServiceAccount
  name: default
  namespace: default
```

### 绑定账户使用 `Group`

当使用 `RoleBinding` 时可以在 namespace 范围内给所有 Service Account 或者所有认证用户绑定角色。

```yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding # <--- RoleBinding
metadata:
  name: default
  namespace: default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole 
  name: privileged-user
subjects:
# Authorize all service accounts in a namespace:
- kind: Group
  apiGroup: rbac.authorization.k8s.io
  name: system:serviceaccounts
# Or equivalently, all authenticated users in a namespace:
- kind: Group
  apiGroup: rbac.authorization.k8s.io
  name: system:authenticated
```

对于 kube-system ns 和集群管理方的 ns 可以直接授权所有 service account 高权限角色。

```yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding # <--- RoleBinding
metadata:
  name: default
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole 
  name: privileged-user
subjects:
# Authorize all service accounts in a namespace:
- kind: Group
  apiGroup: rbac.authorization.k8s.io
  name: system:serviceaccounts
```