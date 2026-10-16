<p align="center">
  <img src="https://img.shields.io/badge/GeoHashForge-v1.0.0-blue?style=flat-square" alt="version"/>
  <img src="https://img.shields.io/badge/license-MIT-green?style=flat-square" alt="license"/>
  <img src="https://img.shields.io/badge/dependencies-zero-brightgreen?style=flat-square" alt="zero deps"/>
  <img src="https://img.shields.io/badge/go-1.22%2B-00ADD8?style=flat-square" alt="go"/>
</p>

<h1 align="center">🌍 GeoHashForge</h1>
<p align="center"><b>零依赖 · 跨平台 · 单机可运行的 Geohash 地理编码工具箱</b></p>

<p align="center">
  <a href="#简体中文">简体中文</a> &nbsp;|&nbsp;
  <a href="#繁體中文">繁體中文</a> &nbsp;|&nbsp;
  <a href="#english">English</a>
</p>

---

<a id="简体中文"></a>
## 🎉 项目介绍（简体中文）

**GeoHashForge** 是一个用 Go 编写、零第三方依赖的 Geohash 地理编码工具箱，同时提供 **命令行工具（CLI）** 与 **可复用的 Go 库**。它能将经纬度编码为简短的 Geohash 字符串，并支持解码、邻域计算、边界框、球面距离、区域覆盖与精度误差分析。

### 核心价值
- 🧭 **解决真实痛点**：在地理搜索、邻近推荐、热力图、地图瓦片、位置分享等场景中，Geohash 是把二维坐标压缩为一维字符串、实现高效前缀检索与聚合的标准方案（被 Elasticsearch、Redis GEO、PostGIS 等广泛采用）。
- ⚡ **零依赖、单二进制**：仅使用 Go 标准库，编译后是一个跨平台可执行文件，无运行时依赖，可离线运行。
- 🧩 **CLI + 库双形态**：既能在终端直接调用，也能作为 `import` 库集成进你的 Go 项目。

### 自研差异化亮点
- 相比社区零散的 Geohash 片段，GeoHashForge 提供了**完整闭环功能**（编码 / 解码 / 邻域 / 边界框 / 距离 / 区域覆盖 / 精度表），并通过**单元测试**保证算法正确性。
- 内置 **`cover` 区域覆盖**子命令，可一键生成覆盖任意矩形范围的 Geohash 网格，直接服务于地图瓦片与地理聚合。
- 全量 **中 / 英双语 README** 与清晰的错误处理，开箱即用。

### 灵感来源
源于地理信息系统中对"把坐标变成可索引字符串"这一经典需求的再实现，目标是提供一个轻量、正确、可审计的开源工具。

---

## ✨ 核心特性（简体中文）

- 🔐 **编码 `encode`**：经纬度 → Geohash，精度 1–12 可选。
- 🔓 **解码 `decode`**：Geohash → 中心点坐标、误差与完整边界框。
- 🧭 **邻域 `neighbors`**：计算 8 个相邻单元（含跨 ±180° 经线环绕处理）。
- 📐 **边界框 `bbox`**：输出单元格的经纬度范围。
- 📏 **距离 `distance`**：基于 Haversine 公式计算两单元中心的球面距离（km）。
- 🗺️ **区域覆盖 `cover`**：生成覆盖指定矩形范围的最优 Geohash 网格集合。
- 📊 **精度表 `table`**：列出各精度下的经 / 纬度误差与近似公里尺寸。
- 🤖 **JSON 输出**：所有子命令支持 `-j` 以 JSON 输出，便于脚本解析。

---

## 🚀 快速开始（简体中文）

### 环境要求
- Go 1.22 或更高版本（仅用于从源码安装；使用发布版二进制无需 Go）。

### 安装方式一：go install（推荐）
```bash
go install github.com/gitstq/GeoHashForge@latest
# 安装后 geohashforge 会位于 $GOPATH/bin，请确保其在 PATH 中
geohashforge --help
```

### 安装方式二：下载发布版二进制
前往 [Releases](../../releases) 页面，下载对应平台（Windows / macOS / Linux，amd64 / arm64）的预编译二进制，解压后即可直接使用。

### 安装方式三：从源码构建
```bash
git clone https://github.com/gitstq/GeoHashForge.git
cd GeoHashForge
go build -o geohashforge .
./geohashforge --help
```

### 一分钟上手
```bash
# 把北京坐标编码为 8 位 Geohash
geohashforge encode 39.9042 116.4074 -p 8
# => wx4g0c9w

# 解码
geohashforge decode wx4g0c

# 计算北京到东京的距离
geohashforge distance wx4g0c xn774c
```

---

## 📖 详细使用指南（简体中文）

### `encode` 编码
```bash
geohashforge encode <纬度> <经度> [-p 精度] [-j]
geohashforge encode 35.6895 139.6917 -p 6
```
| 参数 | 说明 | 默认值 |
|---|---|---|
| `-p` | 编码精度（1–12） | 9 |
| `-j` | JSON 格式输出 | 关闭 |

### `decode` 解码
```bash
geohashforge decode <geohash> [-j]
geohashforge decode wx4g0c -j
```

### `neighbors` 邻域
```bash
geohashforge neighbors <geohash> [-j]
geohashforge neighbors wx4g0c
# 输出 n / ne / e / se / s / sw / w / nw 八个方向
```

### `bbox` 边界框
```bash
geohashforge bbox <geohash> [-j]
```

### `distance` 距离
```bash
geohashforge distance <gh1> <gh2> [-j]
geohashforge distance wx4g0c xn774c
# => 2132.xxx km
```

### `cover` 区域覆盖
```bash
geohashforge cover -minlat 39.9 -minlon 116.4 -maxlat 39.91 -maxlon 116.41 -p 6 [-j]
```
返回覆盖该矩形范围的所有 Geohash 单元，适合地图瓦片与地理聚合。

### `table` 精度表
```bash
geohashforge table
```
展示精度 1–12 对应的经 / 纬度误差（度与公里），帮助你按精度权衡精度与粒度。

> 💡 典型场景：用 `cover` 获取某城市范围内的网格，配合数据库前缀查询即可实现"附近的人 / 店"等邻近检索。

---

## 💡 设计思路与迭代规划（简体中文）

### 设计理念
- **纯算法、零副作用**：所有核心函数为确定性纯函数，不触网、不依赖外部数据，便于测试与审计。
- **标准优先**：严格遵循 Geohash 规范（Base32 去除了 `a,i,l,o`）。
- **双形态交付**：CLI 满足即开即用，库形态满足工程集成。

### 后续迭代计划
- 🧱 增加 `point-in-geohash` 判定与多边形覆盖。
- 🖥️ 可选 TUI 交互模式。
- 🔌 提供 Python / Node 的调用示例与绑定说明。
- 📦 增加 Homebrew / Scoop 等包管理器分发。

### 社区贡献方向
欢迎补充更多单元测试用例、性能基准（benchmarks）与多语言文档。

---

## 📦 打包与部署指南（简体中文）

- **作为 CLI 工具**：从 [Releases](../../releases) 下载对应平台二进制，或 `go install` 即可，无需部署。
- **作为 Go 库**：在你的模块中执行
  ```bash
  go get github.com/gitstq/GeoHashForge
  ```
  然后 `import "github.com/gitstq/GeoHashForge/geohash"` 调用 `geohash.Encode` / `geohash.Decode` 等函数。
- **跨平台构建**：本项目通过 GitHub Actions 在 `windows / macos / ubuntu` + `amd64 / arm64` 矩阵云端构建，产物自动发布至 Releases，无需本地打包。

---

## 🤝 贡献指南（简体中文）

请阅读 [CONTRIBUTING.md](./CONTRIBUTING.md)。核心要求：遵循 Angular 提交规范、保持零依赖、新增功能须补充单元测试、提交前通过 `gofmt -e -l .` 检查。

---

## 📄 开源协议说明（简体中文）

本项目基于 **MIT License** 开源，可自由用于商业与非商业场景，请保留版权声明。详见 [LICENSE](./LICENSE)。

---

<a id="繁體中文"></a>
## 🎉 專案介紹（繁體中文）

**GeoHashForge** 是一套以 Go 撰寫、零第三方依賴的 Geohash 地理編碼工具箱，同時提供 **命令列工具（CLI）** 與 **可重複使用的 Go 函式庫**。它能將經緯度編碼為簡短的 Geohash 字串，並支援解碼、鄰域計算、邊界框、球面距離、區域覆蓋與精度誤差分析。

### 核心價值
- 🧭 **解決真實痛點**：在地圖搜尋、鄰近推薦、熱力圖、地圖瓦片、位置分享等場景中，Geohash 是把二維座標壓縮為一維字串、實現高效前綴檢索與聚合的標準方案（被 Elasticsearch、Redis GEO、PostGIS 等廣泛採用）。
- ⚡ **零依賴、單一二進位**：僅使用 Go 標準函式庫，編譯後是跨平台執行檔，無執行期依賴，可離線執行。
- 🧩 **CLI + 函式庫雙形態**：既可在終端機直接呼叫，也能以 `import` 方式整合進你的 Go 專案。

### 自研差異化亮點
- 相較於社群零散的 Geohash 片段，GeoHashForge 提供**完整閉環功能**（編碼 / 解碼 / 鄰域 / 邊界框 / 距離 / 區域覆蓋 / 精度表），並以**單元測試**確保演算法正確性。
- 內建 **`cover` 區域覆蓋**子命令，可一鍵產生覆蓋任意矩形範圍的 Geohash 網格，直接服務於地圖瓦片與地理聚合。
- 全量 **中 / 英雙語 README** 與清晰的錯誤處理，開箱即用。

### 靈感來源
源於地理資訊系統中對「把座標變成可索引字串」這一經典需求的再實作，目標是提供一個輕量、正確、可審計的開源工具。

---

## ✨ 核心特性（繁體中文）

- 🔐 **編碼 `encode`**：經緯度 → Geohash，精度 1–12 可選。
- 🔓 **解碼 `decode`**：Geohash → 中心點座標、誤差與完整邊界框。
- 🧭 **鄰域 `neighbors`**：計算 8 個相鄰單元（含跨 ±180° 經線環繞處理）。
- 📐 **邊界框 `bbox`**：輸出單元的經緯度範圍。
- 📏 **距離 `distance`**：基於 Haversine 公式計算兩單元中心的球面距離（km）。
- 🗺️ **區域覆蓋 `cover`**：產生覆蓋指定矩形範圍的最佳 Geohash 網格集合。
- 📊 **精度表 `table`**：列出各精度下的經 / 緯度誤差與近似公里尺寸。
- 🤖 **JSON 輸出**：所有子命令支援 `-j` 以 JSON 輸出，便於腳本解析。

---

## 🚀 快速開始（繁體中文）

### 環境需求
- Go 1.22 或更高版本（僅用於從原始碼安裝；使用發行版二進位無需 Go）。

### 安裝方式一：go install（推薦）
```bash
go install github.com/gitstq/GeoHashForge@latest
# 安裝後 geohashforge 會位於 $GOPATH/bin，請確保其在 PATH 中
geohashforge --help
```

### 安裝方式二：下載發行版二進位
前往 [Releases](../../releases) 頁面，下載對應平台（Windows / macOS / Linux，amd64 / arm64）的預編譯二進位，解壓後即可直接使用。

### 安裝方式三：從原始碼建置
```bash
git clone https://github.com/gitstq/GeoHashForge.git
cd GeoHashForge
go build -o geohashforge .
./geohashforge --help
```

### 一分鐘上手
```bash
# 把北京座標編碼為 8 位 Geohash
geohashforge encode 39.9042 116.4074 -p 8

# 解碼
geohashforge decode wx4g0c

# 計算北京到東京的距離
geohashforge distance wx4g0c xn774c
```

---

## 📖 詳細使用指南（繁體中文）

### `encode` 編碼
```bash
geohashforge encode <緯度> <經度> [-p 精度] [-j]
geohashforge encode 35.6895 139.6917 -p 6
```
| 參數 | 說明 | 預設值 |
|---|---|---|
| `-p` | 編碼精度（1–12） | 9 |
| `-j` | JSON 格式輸出 | 關閉 |

### `decode` 解碼
```bash
geohashforge decode <geohash> [-j]
```

### `neighbors` 鄰域
```bash
geohashforge neighbors <geohash> [-j]
```

### `bbox` 邊界框
```bash
geohashforge bbox <geohash> [-j]
```

### `distance` 距離
```bash
geohashforge distance <gh1> <gh2> [-j]
```

### `cover` 區域覆蓋
```bash
geohashforge cover -minlat 39.9 -minlon 116.4 -maxlat 39.91 -maxlon 116.41 -p 6 [-j]
```

### `table` 精度表
```bash
geohashforge table
```

> 💡 典型場景：用 `cover` 取得某城市範圍內的網格，配合資料庫前綴查詢即可實作「附近的人 / 店」等鄰近檢索。

---

## 💡 設計思路與迭代規劃（繁體中文）

### 設計理念
- **純演算法、零副作用**：所有核心函式為確定性純函式，不連網、不依賴外部資料，便於測試與審計。
- **標準優先**：嚴格遵循 Geohash 規範（Base32 去除 `a,i,l,o`）。
- **雙形態交付**：CLI 滿足即開即用，函式庫形態滿足工程整合。

### 後續迭代計畫
- 🧱 增加 `point-in-geohash` 判定與多邊形覆蓋。
- 🖥️ 可選 TUI 互動模式。
- 🔌 提供 Python / Node 的呼叫範例與綁定說明。
- 📦 增加 Homebrew / Scoop 等套件管理器分發。

---

## 📦 打包與部署指南（繁體中文）

- **作為 CLI 工具**：從 [Releases](../../releases) 下載對應平台二進位，或 `go install` 即可，無需部署。
- **作為 Go 函式庫**：在模組中執行 `go get github.com/gitstq/GeoHashForge`，再以 `import "github.com/gitstq/GeoHashForge/geohash"` 呼叫 `geohash.Encode` / `geohash.Decode` 等函式。
- **跨平台建置**：本專案透過 GitHub Actions 在 `windows / macos / ubuntu` + `amd64 / arm64` 矩陣雲端建置，產物自動發佈至 Releases，無需本地打包。

---

## 🤝 貢獻指南（繁體中文）

請閱讀 [CONTRIBUTING.md](./CONTRIBUTING.md)。核心要求：遵循 Angular 提交規範、保持零依賴、新增功能須補充單元測試、提交前通過 `gofmt -e -l .` 檢查。

---

## 📄 開源協議說明（繁體中文）

本專案基於 **MIT License** 開源，可自由用於商業與非商業場景，請保留版權聲明。詳見 [LICENSE](./LICENSE)。

---

<a id="english"></a>
## 🎉 Introduction (English)

**GeoHashForge** is a zero-dependency, cross-platform Geohash geocoding toolkit written in Go. It ships both a **command-line interface (CLI)** and a **reusable Go library**. It encodes latitude/longitude into compact Geohash strings and supports decoding, neighbor lookup, bounding boxes, great-circle distance, region covering, and a precision/error table.

### Core Value
- 🧭 **Solves a real problem**: In geo-search, nearby recommendations, heatmaps, map tiles, and location sharing, Geohash is the standard way to compress 2D coordinates into a 1D string for efficient prefix search and aggregation (adopted by Elasticsearch, Redis GEO, PostGIS, and more).
- ⚡ **Zero dependencies, single binary**: Built only on the Go standard library. The compiled binary is cross-platform, has no runtime dependencies, and runs fully offline.
- 🧩 **CLI + library**: Use it directly from the terminal, or import it as a library in your Go project.

### Differentiation Highlights
- Unlike scattered Geohash snippets, GeoHashForge delivers a **complete feature loop** (encode / decode / neighbors / bbox / distance / cover / table) backed by **unit tests** for algorithmic correctness.
- The built-in **`cover`** subcommand generates the optimal Geohash grid covering any rectangular area — ideal for map tiles and geo aggregation.
- Fully **bilingual (zh/en) README** with clear error handling, ready to use out of the box.

### Inspiration
A clean-room reimplementation of the classic "turn coordinates into an indexable string" need in GIS — lightweight, correct, and auditable.

---

## ✨ Core Features (English)

- 🔐 **encode**: latitude/longitude → Geohash (precision 1–12).
- 🔓 **decode**: Geohash → center coordinates, error, and full bounding box.
- 🧭 **neighbors**: 8 surrounding cells (with ±180° longitude wrap handling).
- 📐 **bbox**: output the cell's lat/lon range.
- 📏 **distance**: great-circle distance (km) between two cell centers via Haversine.
- 🗺️ **cover**: generate the Geohash grid covering a given bounding box.
- 📊 **table**: per-precision lat/lon error and approximate size in km.
- 🤖 **JSON output**: every subcommand supports `-j` for script-friendly JSON.

---

## 🚀 Quick Start (English)

### Requirements
- Go 1.22+ (only needed to install from source; the release binaries need no Go).

### Install via go install (recommended)
```bash
go install github.com/gitstq/GeoHashForge@latest
geohashforge --help
```

### Install via release binary
Download the prebuilt binary for your platform (Windows / macOS / Linux, amd64 / arm64) from the [Releases](../../releases) page and run it directly.

### Build from source
```bash
git clone https://github.com/gitstq/GeoHashForge.git
cd GeoHashForge
go build -o geohashforge .
./geohashforge --help
```

### One-minute tour
```bash
geohashforge encode 39.9042 116.4074 -p 8   # Beijing
geohashforge decode wx4g0c
geohashforge distance wx4g0c xn774c          # Beijing -> Tokyo
```

---

## 📖 Usage Guide (English)

### `encode`
```bash
geohashforge encode <lat> <lon> [-p precision] [-j]
```
| Flag | Description | Default |
|---|---|---|
| `-p` | precision (1–12) | 9 |
| `-j` | JSON output | off |

### `decode`
```bash
geohashforge decode <geohash> [-j]
```

### `neighbors`
```bash
geohashforge neighbors <geohash> [-j]   # n ne e se s sw w nw
```

### `bbox`
```bash
geohashforge bbox <geohash> [-j]
```

### `distance`
```bash
geohashforge distance <gh1> <gh2> [-j]
```

### `cover`
```bash
geohashforge cover -minlat 39.9 -minlon 116.4 -maxlat 39.91 -maxlon 116.41 -p 6 [-j]
```

### `table`
```bash
geohashforge table
```

> 💡 Typical use: use `cover` to fetch the grid over a city, then prefix-match in your database to implement "nearby people / stores" lookups.

---

## 💡 Design & Roadmap (English)

### Design Principles
- **Pure algorithms, no side effects**: all core functions are deterministic and offline — easy to test and audit.
- **Standards first**: strictly follows the Geohash spec (Base32 omits `a,i,l,o`).
- **Dual deliverable**: CLI for instant use, library for integration.

### Roadmap
- 🧱 Add `point-in-geohash` and polygon covering.
- 🖥️ Optional TUI mode.
- 🔌 Python / Node usage examples and binding notes.
- 📦 Homebrew / Scoop distribution.

---

## 📦 Packaging & Deployment (English)

- **As a CLI**: download a release binary or `go install` — no deployment needed.
- **As a Go library**: run `go get github.com/gitstq/GeoHashForge`, then `import "github.com/gitstq/GeoHashForge/geohash"` and call `geohash.Encode` / `geohash.Decode`.
- **Cross-platform builds**: GitHub Actions builds `windows / macos / ubuntu` × `amd64 / arm64` in the cloud; artifacts are published to Releases automatically — no local packaging required.

---

## 🤝 Contributing (English)

See [CONTRIBUTING.md](./CONTRIBUTING.md). Key rules: follow the Angular commit convention, keep zero dependencies, add unit tests for new features, and pass `gofmt -e -l .` before submitting.

---

## 📄 License (English)

Released under the **MIT License** — free for commercial and non-commercial use. Please retain the copyright notice. See [LICENSE](./LICENSE).
