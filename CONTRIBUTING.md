# 贡献指南 / Contributing

感谢你关注 **GeoHashForge**！本项目欢迎所有形式的贡献。

## 🧩 如何参与

1. **提交 Issue**：发现 Bug、提出新功能或改进建议，请先搜索是否已有类似 Issue。
2. **提交 PR**：
   - 从 `main` 分支切出特性分支（`feat/xxx`、`fix/xxx`）。
   - 提交信息遵循 Angular 规范：`feat:` / `fix:` / `docs:` / `refactor:` / `test:`。
   - 确保通过本地语法检查：`gofmt -e -l .` 无输出。
   - 新增功能请同步补充单元测试（`geohash/geohash_test.go`）。

## ✅ 开发约定

- 语言：Go 1.22+，仅使用标准库，**不引入第三方依赖**。
- 代码风格：以 `gofmt` 格式化（Tab 缩进）。
- 算法实现保持纯函数、确定性、可离线运行。
- 文档：README 的多语言版本需同步更新。

## 🚀 本地验证（无需安装依赖）

```bash
gofmt -e -l .        # 语法与格式检查（无输出即通过）
go vet ./...         # 静态检查
go test ./...        # 运行单元测试
```

> 注：本仓库的构建、测试与多平台发布均由 GitHub Actions 云端完成，
> 贡献者无需在本地打包。

## 📮 联系方式

欢迎通过 GitHub Issue 与我们交流。
