// Command geohashforge 是 GeoHashForge 的命令行入口，
// 提供编码、解码、邻域、边界框、距离、区域覆盖与精度表等子命令。
//
// 示例：
//
//	geohashforge encode 39.9042 116.4074 -p 8
//	geohashforge decode wx4g0c
//	geohashforge neighbors wx4g0c
//	geohashforge bbox wx4g0c
//	geohashforge distance wx4g0c xn774c
//	geohashforge cover -minlat 39.9 -minlon 116.4 -maxlat 39.91 -maxlon 116.41 -p 6
//	geohashforge table
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gitstq/GeoHashForge/geohash"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}
	cmd := os.Args[1]
	args := os.Args[2:]

	var err error
	switch cmd {
	case "encode":
		err = cmdEncode(args)
	case "decode":
		err = cmdDecode(args)
	case "neighbors":
		err = cmdNeighbors(args)
	case "bbox":
		err = cmdBbox(args)
	case "distance":
		err = cmdDistance(args)
	case "cover":
		err = cmdCover(args)
	case "table":
		err = cmdTable(args)
	case "-h", "--help", "help":
		printUsage()
		return
	default:
		fmt.Fprintf(os.Stderr, "未知命令: %s\n\n", cmd)
		printUsage()
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`GeoHashForge — 零依赖 Geohash 地理编码工具箱

用法:
  geohashforge <命令> [参数]

命令:
  encode     将经纬度编码为 Geohash         geohashforge encode <纬度> <经度> [-p 精度] [-j]
  decode     解码 Geohash 为坐标与边界框     geohashforge decode <geohash> [-j]
  neighbors  计算 8 个相邻单元              geohashforge neighbors <geohash> [-j]
  bbox       输出 Geohash 的边界框          geohashforge bbox <geohash> [-j]
  distance   计算两个 Geohash 中心的球面距离 geohashforge distance <gh1> <gh2> [-j]
  cover      覆盖经纬度矩形区域的所有单元      geohashforge cover -minlat -minlon -maxlat -maxlon -p 精度 [-j]
  table      输出精度误差对照表              geohashforge table [-j]

选项:
  -j, --json   以 JSON 格式输出（便于脚本解析）
  -h, --help   显示帮助信息
`)
}

// boolFlag 解析 -j / --json 布尔标志。
func newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	return fs
}

func parseJSON(fs *flag.FlagSet) bool {
	j := fs.Bool("json", false, "JSON 输出")
	fs.Parse(fs.Args())
	return *j
}

func cmdEncode(args []string) error {
	fs := newFlagSet("encode")
	precision := fs.Int("p", 9, "编码精度 (1-12)")
	j := fs.Bool("json", false, "JSON 输出")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 2 {
		return fmt.Errorf("用法: encode <纬度> <经度>")
	}
	lat, err := strconv.ParseFloat(fs.Arg(0), 64)
	if err != nil {
		return fmt.Errorf("纬度解析失败: %w", err)
	}
	lon, err := strconv.ParseFloat(fs.Arg(1), 64)
	if err != nil {
		return fmt.Errorf("经度解析失败: %w", err)
	}
	h, err := geohash.Encode(lat, lon, *precision)
	if err != nil {
		return err
	}
	if *j {
		return emitJSON(map[string]any{"geohash": h, "latitude": lat, "longitude": lon, "precision": *precision})
	}
	fmt.Println(h)
	return nil
}

func cmdDecode(args []string) error {
	fs := newFlagSet("decode")
	j := parseJSON(fs)
	if fs.NArg() < 1 {
		return fmt.Errorf("用法: decode <geohash>")
	}
	d, err := geohash.Decode(fs.Arg(0))
	if err != nil {
		return err
	}
	if j {
		return emitJSON(d)
	}
	fmt.Printf("geohash : %s\n", fs.Arg(0))
	fmt.Printf("纬度    : %.6f (±%.6f)\n", d.Latitude, d.LatError)
	fmt.Printf("经度    : %.6f (±%.6f)\n", d.Longitude, d.LonError)
	fmt.Printf("边界框  : 纬度 [%.6f, %.6f]  经度 [%.6f, %.6f]\n",
		d.Box.MinLat, d.Box.MaxLat, d.Box.MinLon, d.Box.MaxLon)
	return nil
}

func cmdNeighbors(args []string) error {
	fs := newFlagSet("neighbors")
	j := parseJSON(fs)
	if fs.NArg() < 1 {
		return fmt.Errorf("用法: neighbors <geohash>")
	}
	n, err := geohash.Neighbors(fs.Arg(0))
	if err != nil {
		return err
	}
	if j {
		return emitJSON(n)
	}
	for _, dir := range []string{"n", "ne", "e", "se", "s", "sw", "w", "nw"} {
		v := n[dir]
		if v == "" {
			v = "(无)"
		}
		fmt.Printf("%-3s : %s\n", dir, v)
	}
	return nil
}

func cmdBbox(args []string) error {
	fs := newFlagSet("bbox")
	j := parseJSON(fs)
	if fs.NArg() < 1 {
		return fmt.Errorf("用法: bbox <geohash>")
	}
	box, err := geohash.BoundingBox(fs.Arg(0))
	if err != nil {
		return err
	}
	if j {
		return emitJSON(box)
	}
	fmt.Printf("纬度 [%.6f, %.6f]\n", box.MinLat, box.MaxLat)
	fmt.Printf("经度 [%.6f, %.6f]\n", box.MinLon, box.MaxLon)
	return nil
}

func cmdDistance(args []string) error {
	fs := newFlagSet("distance")
	j := parseJSON(fs)
	if fs.NArg() < 2 {
		return fmt.Errorf("用法: distance <gh1> <gh2>")
	}
	km, err := geohash.DistanceKm(fs.Arg(0), fs.Arg(1))
	if err != nil {
		return err
	}
	if j {
		return emitJSON(map[string]any{"geohash1": fs.Arg(0), "geohash2": fs.Arg(1), "distance_km": km})
	}
	fmt.Printf("%.3f km\n", km)
	return nil
}

func cmdCover(args []string) error {
	fs := newFlagSet("cover")
	minLat := fs.Float64("minlat", 0, "最小纬度")
	minLon := fs.Float64("minlon", 0, "最小经度")
	maxLat := fs.Float64("maxlat", 0, "最大纬度")
	maxLon := fs.Float64("maxlon", 0, "最大经度")
	precision := fs.Int("p", 6, "编码精度 (1-12)")
	j := fs.Bool("json", false, "JSON 输出")
	if err := fs.Parse(args); err != nil {
		return err
	}
	cells, err := geohash.Cover(*minLat, *minLon, *maxLat, *maxLon, *precision, 0)
	if err != nil {
		return err
	}
	if *j {
		return emitJSON(map[string]any{
			"precision": *precision,
			"count":     len(cells),
			"cells":     cells,
		})
	}
	fmt.Printf("# 覆盖单元数: %d (精度 %d)\n", len(cells), *precision)
	fmt.Println(strings.Join(cells, "\n"))
	return nil
}

func cmdTable(args []string) error {
	fs := newFlagSet("table")
	j := parseJSON(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	table := geohash.PrecisionTable()
	if j {
		return emitJSON(table)
	}
	fmt.Println("精度 | 纬度误差(°) | 经度误差(°) | 纬度误差(km) | 经度误差(km)")
	fmt.Println("------+------------+------------+--------------+--------------")
	for _, r := range table {
		fmt.Printf("%4d  | %9.6f  | %9.6f  | %11.3f   | %11.3f\n",
			r.Precision, r.LatError, r.LonError, r.LatErrorKm, r.LonErrorKm)
	}
	return nil
}

func emitJSON(v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
