# TRPC接口同一管理

## 说明
trpc插件使用见 [tRPC-Go Database 插件](https://github.com/trpc-ecosystem/go-database/blob/main/README.zh_CN.md)

## 功能
- 所有插件采用统一风格管理，维护一份配置文件
- 省略掉插件的创建过程，统一在初始化时创建，代码层随用随取
- 代码层面，保持引用后编译的原则，不引入使用仅加载文件不会引入额外插件
- 个别插件做了二次开发，增强场景使用