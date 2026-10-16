// Package geohash 提供 Geohash 地理编码的完整算法实现：
// 编码（经纬度 → Geohash 字符串）、解码（Geohash → 中心点/边界框/误差）、
// 邻域计算、区域覆盖、精度误差表与球面距离计算。
//
// 所有函数均为纯算法实现，不依赖任何第三方库，可离线确定性运行。
// Geohash 是一种将二维经纬度编码为一维短字符串的地理网格系统，
// 被 Elasticsearch、Redis GEO、PostGIS 等广泛采用。
package geohash

import (
	"errors"
	"math"
	"strings"
)

// Base32 为 Geohash 使用的 32 字符编码表（去除了 a,i,l,o 以避免混淆）。
const Base32 = "0123456789bcdefghjkmnpqrstuvwxyz"

// MaxPrecision 为 Geohash 支持的最大精度（字符数）。
const MaxPrecision = 12

// 地球平均半径（千米），用于球面距离计算。
const earthRadiusKm = 6371.0088

// 经纬度合法范围。
const (
	minLat = -90.0
	maxLat = 90.0
	minLon = -180.0
	maxLon = 180.0
)

// 预定义错误，便于调用方进行精确的错误判断。
var (
	ErrInvalidPrecision = errors.New("geohash: 精度必须在 1 到 12 之间")
	ErrInvalidLatitude  = errors.New("geohash: 纬度必须在 -90 到 90 之间")
	ErrInvalidLongitude = errors.New("geohash: 经度必须在 -180 到 180 之间")
	ErrInvalidHash      = errors.New("geohash: 包含非法的 Geohash 字符")
	ErrTooManyCells     = errors.New("geohash: 覆盖单元格数量过多，请降低精度")
)

// Box 表示一个 Geohash 网格的地理边界框。
type Box struct {
	MinLat float64 `json:"min_lat"`
	MinLon float64 `json:"min_lon"`
	MaxLat float64 `json:"max_lat"`
	MaxLon float64 `json:"max_lon"`
}

// Decoded 表示解码后的完整信息。
type Decoded struct {
	Latitude  float64 `json:"latitude"`  // 单元格中心纬度
	Longitude float64 `json:"longitude"` // 单元格中心经度
	LatError  float64 `json:"lat_error"` // 纬度半宽（误差）
	LonError  float64 `json:"lon_error"` // 经度半宽（误差）
	Box       Box     `json:"box"`       // 完整边界框
}

// PrecisionInfo 描述某一精度下的误差与近似尺寸。
type PrecisionInfo struct {
	Precision  int     `json:"precision"`
	LatError   float64 `json:"lat_error_deg"`
	LonError   float64 `json:"lon_error_deg"`
	LatErrorKm float64 `json:"lat_error_km"`
	LonErrorKm float64 `json:"lon_error_km"`
}

// Encode 将经纬度编码为指定精度的 Geohash 字符串。
// precision 取值范围为 1..12，精度越高字符串越长、地理范围越小。
func Encode(lat, lon float64, precision int) (string, error) {
	if precision < 1 || precision > MaxPrecision {
		return "", ErrInvalidPrecision
	}
	if lat < minLat || lat > maxLat {
		return "", ErrInvalidLatitude
	}
	if lon < minLon || lon > maxLon {
		return "", ErrInvalidLongitude
	}

	latRange := [2]float64{minLat, maxLat}
	lonRange := [2]float64{minLon, maxLon}

	var sb strings.Builder
	var idx, bits int
	even := true // even 为 true 时处理经度，false 时处理纬度

	for sb.Len() < precision {
		if even {
			mid := (lonRange[0] + lonRange[1]) / 2
			if lon >= mid {
				idx = idx*2 + 1
				lonRange[0] = mid
			} else {
				idx = idx * 2
				lonRange[1] = mid
			}
		} else {
			mid := (latRange[0] + latRange[1]) / 2
			if lat >= mid {
				idx = idx*2 + 1
				latRange[0] = mid
			} else {
				idx = idx * 2
				latRange[1] = mid
			}
		}
		even = !even
		bits++
		if bits == 5 {
			sb.WriteByte(Base32[idx])
			bits = 0
			idx = 0
		}
	}
	return sb.String(), nil
}

// Decode 将 Geohash 字符串解码为经纬度中心、误差与边界框。
func Decode(hash string) (Decoded, error) {
	hash = strings.ToLower(strings.TrimSpace(hash))
	if hash == "" {
		return Decoded{}, ErrInvalidHash
	}

	latRange := [2]float64{minLat, maxLat}
	lonRange := [2]float64{minLon, maxLon}
	even := true

	for _, c := range hash {
		idx := strings.IndexRune(Base32, c)
		if idx < 0 {
			return Decoded{}, ErrInvalidHash
		}
		for i := 4; i >= 0; i-- {
			bit := (idx >> uint(i)) & 1
			if even {
				mid := (lonRange[0] + lonRange[1]) / 2
				if bit == 1 {
					lonRange[0] = mid
				} else {
					lonRange[1] = mid
				}
			} else {
				mid := (latRange[0] + latRange[1]) / 2
				if bit == 1 {
					latRange[0] = mid
				} else {
					latRange[1] = mid
				}
			}
			even = !even
		}
	}

	box := Box{
		MinLat: latRange[0],
		MinLon: lonRange[0],
		MaxLat: latRange[1],
		MaxLon: lonRange[1],
	}
	return Decoded{
		Latitude:  (latRange[0] + latRange[1]) / 2,
		Longitude: (lonRange[0] + lonRange[1]) / 2,
		LatError:  (latRange[1] - latRange[0]) / 2,
		LonError:  (lonRange[1] - lonRange[0]) / 2,
		Box:       box,
	}, nil
}

// BoundingBox 返回 Geohash 对应的边界框。
func BoundingBox(hash string) (Box, error) {
	d, err := Decode(hash)
	if err != nil {
		return Box{}, err
	}
	return d.Box, nil
}

// Neighbors 计算 Geohash 的 8 个相邻单元（北、南、东、西及四个对角方向）。
// 返回以方向名为键的映射；位于极地或跨越 ±180° 经线时对应方向可能为空字符串。
func Neighbors(hash string) (map[string]string, error) {
	d, err := Decode(hash)
	if err != nil {
		return nil, err
	}
	precision := len(hash)
	latCenter := d.Latitude
	lonCenter := d.Longitude
	latErr := d.LatError
	lonErr := d.LonError

	dirs := map[string][2]float64{
		"n":  {2, 0},
		"ne": {2, 2},
		"e":  {0, 2},
		"se": {-2, 2},
		"s":  {-2, 0},
		"sw": {-2, -2},
		"w":  {0, -2},
		"nw": {2, -2},
	}

	result := make(map[string]string, len(dirs))
	for name, delta := range dirs {
		nlat := latCenter + delta[0]*latErr
		nlon := lonCenter + delta[1]*lonErr
		// 经度环绕到 [-180, 180)
		for nlon >= maxLon {
			nlon -= 360
		}
		for nlon < minLon {
			nlon += 360
		}
		if nlat > maxLat || nlat < minLat {
			result[name] = ""
			continue
		}
		nh, err := Encode(nlat, nlon, precision)
		if err != nil {
			result[name] = ""
			continue
		}
		result[name] = nh
	}
	return result, nil
}

// DistanceKm 计算两个 Geohash 中心点之间的球面距离（单位：千米）。
func DistanceKm(hash1, hash2 string) (float64, error) {
	d1, err := Decode(hash1)
	if err != nil {
		return 0, err
	}
	d2, err := Decode(hash2)
	if err != nil {
		return 0, err
	}
	return haversineKm(d1.Latitude, d1.Longitude, d2.Latitude, d2.Longitude), nil
}

// PrecisionTable 返回精度 1..12 的误差与近似尺寸信息。
func PrecisionTable() []PrecisionInfo {
	out := make([]PrecisionInfo, 0, MaxPrecision)
	for p := 1; p <= MaxPrecision; p++ {
		box, err := BoundingBox(strings.Repeat("s", p)) // "s" 位于网格中部，用作参考单元
		if err != nil {
			continue
		}
		latErr := (box.MaxLat - box.MinLat) / 2
		lonErr := (box.MaxLon - box.MinLon) / 2
		out = append(out, PrecisionInfo{
			Precision:  p,
			LatError:   latErr,
			LonError:   lonErr,
			LatErrorKm: latErr * 110.574,
			LonErrorKm: lonErr * 111.320,
		})
	}
	return out
}

// Cover 返回覆盖给定经纬度边界框（minLat,minLon 到 maxLat,maxLon）的所有 Geohash 单元。
// precision 决定单元格粒度；当生成的单元数量超过 maxCells 时返回 ErrTooManyCells。
func Cover(minLatV, minLonV, maxLatV, maxLonV float64, precision, maxCells int) ([]string, error) {
	if precision < 1 || precision > MaxPrecision {
		return nil, ErrInvalidPrecision
	}
	if minLatV < minLat || maxLatV > maxLat || minLatV > maxLatV {
		return nil, ErrInvalidLatitude
	}
	if minLonV < minLon || maxLonV > maxLon || minLonV > maxLonV {
		return nil, ErrInvalidLongitude
	}

	ref, err := BoundingBox(strings.Repeat("s", precision))
	if err != nil {
		return nil, err
	}
	latStep := ref.MaxLat - ref.MinLat
	lonStep := ref.MaxLon - ref.MinLon

	iLo := int(math.Floor((minLatV - minLat) / latStep))
	iHi := int(math.Floor((maxLatV - minLat) / latStep))
	jLo := int(math.Floor((minLonV - minLon) / lonStep))
	jHi := int(math.Floor((maxLonV - minLon) / lonStep))

	total := (iHi - iLo + 1) * (jHi - jLo + 1)
	if maxCells > 0 && total > maxCells {
		return nil, ErrTooManyCells
	}

	var out []string
	for i := iLo; i <= iHi; i++ {
		for j := jLo; j <= jHi; j++ {
			clat := minLat + (float64(i)+0.5)*latStep
			clon := minLon + (float64(j)+0.5)*lonStep
			g, err := Encode(clat, clon, precision)
			if err != nil {
				continue
			}
			bb, err := BoundingBox(g)
			if err != nil {
				continue
			}
			if bb.MaxLat < minLatV || bb.MinLat > maxLatV ||
				bb.MaxLon < minLonV || bb.MinLon > maxLonV {
				continue
			}
			out = append(out, g)
		}
	}
	return out, nil
}

// haversineKm 使用 Haversine 公式计算两点间的球面距离（千米）。
func haversineKm(lat1, lon1, lat2, lon2 float64) float64 {
	φ1 := lat1 * math.Pi / 180
	φ2 := lat2 * math.Pi / 180
	dφ := (lat2 - lat1) * math.Pi / 180
	dλ := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dφ/2)*math.Sin(dφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*math.Sin(dλ/2)*math.Sin(dλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}
