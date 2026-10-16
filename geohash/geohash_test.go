package geohash

import (
	"math"
	"testing"
)

// 已知对照值（来源：公共 Geohash 参考资料）。
func TestEncodeKnown(t *testing.T) {
	cases := []struct {
		lat, lon float64
		prec     int
		want     string
	}{
		{57.64911, 10.40744, 12, "u4pruydqqvj8"},
		{57.64911, 10.40744, 6, "u4pruy"},
		{39.9042, 116.4074, 6, "wx4g0c"},  // 北京
		{35.6895, 139.6917, 6, "xn774c"},  // 东京
		{-33.8688, 151.2093, 6, "r3gx2f"}, // 悉尼
	}
	for _, c := range cases {
		got, err := Encode(c.lat, c.lon, c.prec)
		if err != nil {
			t.Fatalf("Encode(%v,%v,%v) err=%v", c.lat, c.lon, c.prec, err)
		}
		if got != c.want {
			t.Errorf("Encode(%v,%v,%v)=%q, want %q", c.lat, c.lon, c.prec, got, c.want)
		}
	}
}

func TestEncodeErrors(t *testing.T) {
	if _, err := Encode(0, 0, 0); err != ErrInvalidPrecision {
		t.Errorf("precision 0: got %v, want ErrInvalidPrecision", err)
	}
	if _, err := Encode(0, 0, 13); err != ErrInvalidPrecision {
		t.Errorf("precision 13: got %v, want ErrInvalidPrecision", err)
	}
	if _, err := Encode(91, 0, 5); err != ErrInvalidLatitude {
		t.Errorf("lat 91: got %v, want ErrInvalidLatitude", err)
	}
	if _, err := Encode(0, 181, 5); err != ErrInvalidLongitude {
		t.Errorf("lon 181: got %v, want ErrInvalidLongitude", err)
	}
}

func TestDecodeRoundTrip(t *testing.T) {
	for _, h := range []string{"u4pruydqqvj8", "wx4g0c", "xn774c", "r3gx2f", "9q8yyk8"} {
		d, err := Decode(h)
		if err != nil {
			t.Fatalf("Decode(%q) err=%v", h, err)
		}
		// 中心点重新编码应得到相同哈希
		got, err := Encode(d.Latitude, d.Longitude, len(h))
		if err != nil {
			t.Fatalf("re-encode err=%v", err)
		}
		if got != h {
			t.Errorf("round-trip %q -> (%v,%v) -> %q", h, d.Latitude, d.Longitude, got)
		}
	}
}

func TestDecodeErrors(t *testing.T) {
	if _, err := Decode(""); err != ErrInvalidHash {
		t.Errorf("empty: got %v, want ErrInvalidHash", err)
	}
	if _, err := Decode("u4pruy!"); err != ErrInvalidHash {
		t.Errorf("bad char: got %v, want ErrInvalidHash", err)
	}
}

func TestBoundingBox(t *testing.T) {
	box, err := BoundingBox("u4pruy")
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	// 中心应落在边界框内
	d, _ := Decode("u4pruy")
	if d.Latitude < box.MinLat || d.Latitude > box.MaxLat {
		t.Errorf("lat %v outside box %v..%v", d.Latitude, box.MinLat, box.MaxLat)
	}
	if d.Longitude < box.MinLon || d.Longitude > box.MaxLon {
		t.Errorf("lon %v outside box %v..%v", d.Longitude, box.MinLon, box.MaxLon)
	}
}

func TestNeighbors(t *testing.T) {
	// 以 "u4pruy" 为例验证北邻与南邻的方向关系
	n, err := Neighbors("u4pruy")
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	for _, dir := range []string{"n", "s", "e", "w", "ne", "nw", "se", "sw"} {
		if n[dir] == "" {
			t.Errorf("neighbor %q empty", dir)
		}
	}
	// 北邻中心纬度应高于原中心纬度
	d0, _ := Decode("u4pruy")
	dn, _ := Decode(n["n"])
	if dn.Latitude <= d0.Latitude {
		t.Errorf("north neighbor lat %v should be > %v", dn.Latitude, d0.Latitude)
	}
	ds, _ := Decode(n["s"])
	if ds.Latitude >= d0.Latitude {
		t.Errorf("south neighbor lat %v should be < %v", ds.Latitude, d0.Latitude)
	}
}

func TestDistance(t *testing.T) {
	// 北京到东京的近似距离约 2100 km
	d, err := DistanceKm("wx4g0c", "xn774c")
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if math.Abs(d-2100) > 100 {
		t.Errorf("Beijing-Tokyo distance=%v km, want ~2100", d)
	}
	// 同一点的距离为 0
	z, err := DistanceKm("wx4g0c", "wx4g0c")
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if z != 0 {
		t.Errorf("self distance=%v, want 0", z)
	}
}

func TestPrecisionTable(t *testing.T) {
	table := PrecisionTable()
	if len(table) != MaxPrecision {
		t.Fatalf("table len=%d, want %d", len(table), MaxPrecision)
	}
	// 精度越高，误差越小
	if table[0].LatError <= table[1].LatError {
		t.Errorf("precision 1 error %v should be > precision 2 error %v",
			table[0].LatError, table[1].LatError)
	}
}

func TestCover(t *testing.T) {
	// 覆盖一个很小的区域，应返回有限且合理的单元集合
	cells, err := Cover(39.90, 116.40, 39.91, 116.41, 6, 10000)
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if len(cells) == 0 {
		t.Errorf("cover returned no cells")
	}
	// 每个单元都应与查询框相交（粗略校验：中心在框内附近）
	for _, c := range cells {
		d, err := Decode(c)
		if err != nil {
			t.Fatalf("decode %q err=%v", c, err)
		}
		if d.Latitude < 39.89 || d.Latitude > 39.92 ||
			d.Longitude < 116.39 || d.Longitude > 116.42 {
			t.Errorf("cell %q center (%v,%v) outside expected region",
				c, d.Latitude, d.Longitude)
		}
	}
}

func TestCoverTooMany(t *testing.T) {
	_, err := Cover(minLat, minLon, maxLat, maxLon, 4, 10)
	if err != ErrTooManyCells {
		t.Errorf("got %v, want ErrTooManyCells", err)
	}
}
