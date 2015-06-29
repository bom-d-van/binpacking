package binpacking

import (
	"fmt"
	"reflect"

	"testing"
)

func BenchmarkPack(b *testing.B) {
	items := []Item{
		goods{20, 100, 30},
		goods{100, 20, 30},
		goods{20, 100, 30},
		goods{100, 20, 30},
		goods{100, 20, 30},
		goods{100, 100, 30},
		goods{100, 100, 30},
		goods{100, 100, 30},
		goods{100, 100, 30},
		goods{100, 100, 30},
		goods{100, 100, 30},
		goods{100, 100, 30},
		goods{100, 100, 30},
		goods{100, 100, 30},
	}
	for i := 0; i < b.N; i++ {
		_, err := Pack(items)
		if err != nil {
			b.Error(err)
		}
	}
}

type goods [4]int

func (g goods) GetWidth() int {
	return g[1]
}

func (g goods) GetHeight() int {
	return g[2]
}

func (g goods) GetDepth() int {
	return g[3]
}

func (g goods) GetWeight() int {
	return 10
}

func TestPack(t *testing.T) {
	items := []Item{
		goods{1, 20, 100, 30},
		goods{2, 100, 20, 30},
		goods{3, 20, 100, 30},
		goods{4, 100, 20, 30},
		goods{5, 100, 20, 30},
		goods{6, 100, 100, 30},
		goods{7, 100, 100, 30},
	}
	want := []Box{BoxSamples[0]}
	want[0].Items = []BoxItem{
		{Item: items[5], RType: 0, Pos: [3]int{0, 0, 0}},
		{Item: items[6], RType: 0, Pos: [3]int{100, 0, 0}},
		{Item: items[0], RType: 0, Pos: [3]int{200, 0, 0}},
		{Item: items[1], RType: 0, Pos: [3]int{0, 100, 0}},
		{Item: items[2], RType: 1, Pos: [3]int{100, 100, 0}},
		{Item: items[3], RType: 2, Pos: [3]int{200, 100, 0}},
		{Item: items[4], RType: 0, Pos: [3]int{0, 120, 0}},
	}

	got, err := Pack(items)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got:\n%swant:\n%s", printBoxes(got), printBoxes(want))
	}
}

func printBoxes(boxes []Box) (r string) {
	for i, box := range boxes {
		r += fmt.Sprintln("box", i, box.Width, box.Height, box.Depth, len(box.Items))
		for i, item := range box.Items {
			r += fmt.Sprintln("  ", i, item)
		}
	}

	return
}
