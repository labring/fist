# Namespace管理模块
namespace在多租户管理系统中扮演重要角色

我们希望在namespace创建时就给namespace特定配额与绑定PSP，这在多租户时特别重要

首先需要禁止掉所有普通用户的创建namespace权限

两种模式下namespace管理的区别：

1. 使用管理员创建namespace 贴用户信息label，组信息label (组信息动态变更如何处理？)
2. 使用管理员创建quota
3. 使用管理员给用户绑定角色，角色里包含PSP
4. 管理员给新建的namespace的default service account绑定包含PSP的角色

# 类似公有云(有计量计费)
用户可以自行申请namespace，因为需要在申请namespace时进行付钱，所以不用太担心用户胡乱申请资源, 不过系统会设置所有租户的资源上限。

```
User                fist                    k8s
 |   create namespace |                      |
 |------------------->|using admin create ns |
 |                    |--------------------->|
 |                    |create role           |
 |                    |quota                 |
 |                    |PSP                   |
 |<-------------------|                      |
 |                    |                      |
```

# 私有云 （仅管理员有权限创建namespace）
这类比较适合管理员创建namespace，给定特定配额，创建一个角色具有访问该namespace权限，再把角色分配给特定的用户或者组。
