# migrate-mc
用于迁移我的MineCraft。

## 用例
我的第一目的是：将Minecraft的mods迁移到高版本，之前这是手动完成的。

这迁移自动的可能在于：不依赖于魔改mod和配置（之后会考虑迁移少量配置文件比如按键绑定）。也就是我的有几百个mod的“原版”(😆)Minecraft。

粗糙点处理就是直接调api完成，专门为此则是为处理一些琐碎细节（TODO)。

## 安装
```shell
go install github.com/wind-mask/migrate-mc
```