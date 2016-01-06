// Package binpacking is a Golang 3D Bin Packing Implementation
//
package binpacking

import (
	"fmt"
	"sort"
)

var (
	BoxSamples = []Box{
		{Width: 220, Height: 160, Depth: 100, Weight: 110, Name: "Box1"},
		{Width: 260, Height: 145, Depth: 145, Weight: 120, Name: "Box2"},
		{Width: 270, Height: 185, Depth: 110, Weight: 140, Name: "Box3"},
		{Width: 310, Height: 220, Depth: 140, Weight: 210, Name: "Box4"},
		{Width: 300, Height: 210, Depth: 200, Weight: 250, Name: "Box5"},
		{Width: 300, Height: 300, Depth: 130, Weight: 290, Name: "Box6"},
		{Width: 370, Height: 270, Depth: 150, Weight: 300, Name: "Box7"},
		{Width: 300, Height: 300, Depth: 250, Weight: 360, Name: "Box8"},
		{Width: 470, Height: 280, Depth: 210, Weight: 400, Name: "Box9"},
		{Width: 430, Height: 315, Depth: 200, Weight: 430, Name: "Box10"},
		{Width: 330, Height: 330, Depth: 350, Weight: 500, Name: "Box11"},
		{Width: 465, Height: 350, Depth: 370, Weight: 650, Name: "Box12"},
	}
)

type RotationType int

const (
	RT1 RotationType = iota // w, h, d
	RT2                     // h, w, d
	RT3                     // h, d, w
	RT4                     // d, h, w
	RT5                     // d, w, h
	RT6                     // w, d, h
)

func (rt RotationType) String() string {
	switch rt {
	case RT1:
		return "RT1 (w, h, d)"
	case RT2:
		return "RT2 (h, w, d)"
	case RT3:
		return "RT3 (h, d, w)"
	case RT4:
		return "RT4 (d, h, w)"
	case RT5:
		return "RT5 (d, w, h)"
	case RT6:
		return "RT6 (w, d, h)"
	}

	return "wrong type"
}

type Box struct {
	Name   string
	Width  int // unit: mm
	Height int // unit: mm
	Depth  int // unit: mm
	Weight int // unit: g

	Items []BoxItem
}

func (b Box) String() string {
	r := fmt.Sprintf(
		"box (w: %d, h: %d, d: %d, weight: %d, name: %s) itemsCount: %d",
		b.Width, b.Height, b.Depth, b.Weight, b.Name, len(b.Items),
	)
	for i, item := range b.Items {
		r += fmt.Sprintf("\n  item %d: %s", i, item)
	}
	return r
}

func (b Box) IsValid() bool {
	return b.volume() != 0
}

func (b Box) volume() int {
	return b.Width * b.Height * b.Depth
}

func (b Box) TotalWeight() (w int) {
	w += b.Weight
	for _, item := range b.Items {
		w += item.GetWeight()
	}
	return
}

func (b Box) size() (s int) {
	return b.Width + b.Height + b.Depth
}

func (b Box) nonBoxItems() (r []Item) {
	for _, item := range b.Items {
		r = append(r, item.Item)
	}
	return
}

type BoxItem struct {
	Item
	Pos   [3]int // (w, h, d)
	RType RotationType
}

func (bi BoxItem) String() string {
	return fmt.Sprintf(
		"item(w: %d, h: %d, d: %d, weight: %d) pos(w: %d, h: %d, d: %d) rtype(%s)",
		bi.Item.GetWidth(), bi.Item.GetHeight(), bi.Item.GetDepth(), bi.Item.GetWeight(),
		bi.Pos[0], bi.Pos[1], bi.Pos[2],
		bi.RType,
	)
}

func (bi *BoxItem) volume() int {
	return bi.GetWidth() * bi.GetHeight() * bi.GetDepth()
}

//     +-----------------+
//    /|                /|
//   / |               / |
//  /  |              /  |
// +-----------------+   |
// |   |             |   |
// |   |             |   |
// |   H             |   |
// |   |             |   |
// |   |             |   |
// |   O----W--------|---+
// |  /              |  /
// | D               | /
// |/                |/
// +-----------------+
//
func (bi1 BoxItem) intersect(bi2 BoxItem) bool {
	d1 := bi1.Dimensions()
	d2 := bi2.Dimensions()
	return intersect([2]int{bi1.Pos[0], bi1.Pos[1]}, [2]int{bi2.Pos[0], bi2.Pos[1]}, d1[0], d1[1], d2[0], d2[1]) &&
		intersect([2]int{bi1.Pos[1], bi1.Pos[2]}, [2]int{bi2.Pos[1], bi2.Pos[2]}, d1[1], d1[2], d2[1], d2[2]) &&
		intersect([2]int{bi1.Pos[0], bi1.Pos[2]}, [2]int{bi2.Pos[0], bi2.Pos[2]}, d1[0], d1[2], d2[0], d2[2])
}

//
// O------X-------+
// |              |
// Y      *       |
// |              |
// +--------------+
// intersect checks if two rectangles overlap
func intersect(o1, o2 [2]int, x1, y1, x2, y2 int) bool {
	centerx1 := o1[0] + x1/2
	centery1 := o1[1] + y1/2
	centerx2 := o2[0] + x2/2
	centery2 := o2[1] + y2/2
	var x, y int
	if centerx1 > centerx2 {
		x = centerx1 - centerx2
	} else {
		x = centerx2 - centerx1
	}
	if centery1 > centery2 {
		y = centery1 - centery2
	} else {
		y = centery2 - centery1
	}
	return x < (x1+x2)/2 && y < (y1+y2)/2
}

func (bi BoxItem) Dimensions() (d []int) {
	switch bi.RType {
	case RT1:
		d = []int{bi.GetWidth(), bi.GetHeight(), bi.GetDepth()}
	case RT2:
		d = []int{bi.GetHeight(), bi.GetWidth(), bi.GetDepth()}
	case RT3:
		d = []int{bi.GetHeight(), bi.GetDepth(), bi.GetWidth()}
	case RT4:
		d = []int{bi.GetDepth(), bi.GetHeight(), bi.GetWidth()}
	case RT5:
		d = []int{bi.GetDepth(), bi.GetWidth(), bi.GetHeight()}
	case RT6:
		d = []int{bi.GetWidth(), bi.GetDepth(), bi.GetHeight()}
	}
	return
}

type Item interface {
	GetHeight() int // unit: mm
	GetWidth() int  // unit: mm
	GetDepth() int  // unit: mm
	GetWeight() int // unit: g
}

type Items []Item

func (is Items) Len() int {
	return len(is)
}

func (is Items) Less(i int, j int) bool {
	return is[i].GetWidth()*is[i].GetHeight()*is[i].GetDepth() > is[j].GetWidth()*is[j].GetHeight()*is[j].GetDepth()
}

func (is Items) Swap(i int, j int) {
	swap := is[i]
	is[i] = is[j]
	is[j] = swap
}

// Original Algorithm: https://github.com/bom-d-van/binpacking/blob/master/erick_dube_507-034.pdf
// The current implementation is based on it, but with some tweaks to fit our requirements.
//
// The original algorithm is designed for identical bins but our requirements is made for
// bins in various sizes
func Pack(notPacked []Item) (boxes []Box, err error) {
	sort.Sort(Items(notPacked))
	for len(notPacked) > 0 {
		toPack := notPacked
		// notPacked = []Item{} // clear notPacked

		currentBin := pickBox(toPack[0])
		if !currentBin.IsValid() {
			err = fmt.Errorf(
				"item too big: {width: %d, height: %d, depth: %d, weight: %d}",
				toPack[0].GetWidth(),
				toPack[0].GetHeight(),
				toPack[0].GetDepth(),
				toPack[0].GetWeight(),
			)
			return
		}

		notPacked = pack(&currentBin, toPack, true)

		if len(currentBin.Items) > 0 {
			boxes = append(boxes, currentBin)
		}
	}

	return
}

func pack(currentBin *Box, toPack []Item, replaceBin bool) (notPacked []Item) {
	if !currentBin.place(toPack[0], [3]int{}) {
		if nbin := getBiggerBox(*currentBin); nbin.IsValid() {
			*currentBin = nbin
			return pack(currentBin, toPack, replaceBin)
		}

		return toPack
	}

	for _, currenItem := range toPack[1:] {
		var fitted bool
	lookup:
		for p := 0; p < 3; p++ {
			for _, binItem := range currentBin.Items {
				var pos [3]int
				switch p {
				case 0:
					pos = [3]int{binItem.Pos[0] + binItem.GetWidth(), binItem.Pos[1], binItem.Pos[2]}
				case 1:
					pos = [3]int{binItem.Pos[0], binItem.Pos[1] + binItem.GetHeight(), binItem.Pos[2]}
				case 2:
					pos = [3]int{binItem.Pos[0], binItem.Pos[1], binItem.Pos[2] + binItem.GetDepth()}
				}

				if currentBin.place(currenItem, pos) {
					fitted = true
					break lookup
				}
			}
		}
		if !fitted {
			if replaceBin {
				for nbin := getBiggerBox(*currentBin); nbin.IsValid(); nbin = getBiggerBox(nbin) {
					left := pack(&nbin, append(currentBin.nonBoxItems(), currenItem), false)
					if len(left) == 0 {
						*currentBin = nbin
						fitted = true
						break
					}
				}
			}

			if !fitted {
				notPacked = append(notPacked, currenItem)
			}
		}
	}

	return
}

func (b *Box) place(item Item, pos [3]int) (fit bool) {
	bi := BoxItem{Item: item, Pos: pos}
	for i := 0; i < 6; i++ {
		bi.RType = RotationType(i)
		d := bi.Dimensions()
		if b.Width < pos[0]+d[0] || b.Height < pos[1]+d[1] || b.Depth < pos[2]+d[2] {
			continue
		}
		fit = true
		for _, item := range b.Items {
			if item.intersect(bi) {
				fit = false
				break
			}
		}
		if fit {
			b.Items = append(b.Items, bi)
			break
		}
		return
	}

	return
}

func pickBox(item Item) Box {
	for _, b := range BoxSamples {
		if !b.place(item, [3]int{}) {
			continue
		}
		b.Items = []BoxItem{}
		return b
	}
	return Box{}
}

func getBiggerBox(box Box) Box {
	v := box.volume()
	for _, b := range BoxSamples {
		if b.volume() > v {
			return b
		}
	}

	return Box{}
}
