package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/daishe/go-future"
	"golang.org/x/sync/errgroup"
)

var dagASCII = []string{
	" ┌───────────────┐    ┌────────────┐    ┌──────────────┐┌───────────────────┐┌───────────┐    ┌────────────┐┌────────────┐ ",
	" │ get spaghetti │    │ boil water │    │ get tomatoes ││ get slicing board ││ get onion │    │ get garlic ││ get grater │ ",
	" └─────────────┬─┘    └─┬──────────┘    └─────────┬────┘└───────┬───────────┘└────┬──────┘    └──────────┬─┘└─┬──────────┘ ",
	" raw spaghetti │        │ boiling water  tomatoes │             │ slicing board   │ onion         garlic │    │ grater     ",
	"               │        ├──────┐                  │             │                 │                      │    │            ",
	"          ┌────▼────────▼──┐   │                  └──────┐      │      ┌──────────┘                      │    │            ",
	"          │ cook spaghetti │   │                       ┌─▼──────▼──────▼─┐                          ┌────▼────▼────┐       ",
	"          └──────────────┬─┘   │                       │ chop vegetables │                          │ grate garlic │       ",
	"        cooked spaghetti │     │                       └───────────┬───┬─┘                          └─┬────────────┘       ",
	"                         │     │                  chopped tomatoes │   │ chopped onion                │ grated garlic      ",
	"                         │     │                                   │   │                              │                    ",
	"                         │     └────────────────────────────────┐  │   │      ┌───────────────────────┘                    ",
	"                         │                                    ┌─▼──▼───▼──────▼─┐                                          ",
	"                         │                                    │ cook vegetables │                                          ",
	"                         │                                    └────────┬────────┘                                          ",
	"                         │                                             │ cooked vegetables                                 ",
	"                         │                                             │                                                   ",
	"                         └─────────────────────────────┐          ┌────┘                                                   ",
	"                                                     ┌─▼──────────▼─┐                                                      ",
	"                                                     │ put on plate │                                                      ",
	"                                                     └──────┬───────┘                                                      ",
	"                                                            │ dish                                                         ",
	"                                                            │                                                              ",
	"                                                            │                                                              ",
	"                                                      ┌─────▼─────┐                                                        ",
	"                                                      │ completed │                                                        ",
	"                                                      └───────────┘                                                        ",
}

func main() {
	fmt.Printf("Preparing spaghetti with tomatoes:\n %s", strings.Join(dagASCII, "\n "))
	fmt.Println("")

	eg, ctx := errgroup.WithContext(context.Background())

	boilingWater := BoilWater(ctx, eg)
	rawSpaghetti := GetSpaghetti(ctx, eg)
	tomatoes := GetTomatoes(ctx, eg)
	slicingBoard := GetSlicingBoard(ctx, eg)
	onion := GetOnion(ctx, eg)
	garlic := GetGarlic(ctx, eg)
	grater := GetGrater(ctx, eg)
	cookedSpaghetti := CookSpaghetti(ctx, eg, boilingWater, rawSpaghetti)
	choppedTomatoes, choppedOnion := ChopVegetables(ctx, eg, slicingBoard, tomatoes, onion)
	grateGarlic := GrateGarlic(ctx, eg, grater, garlic)
	cookedVegetables := CookVegetables(ctx, eg, boilingWater, choppedTomatoes, choppedOnion, grateGarlic)
	dish := PutOnPlate(ctx, eg, cookedSpaghetti, cookedVegetables)

	if err := eg.Wait(); err != nil {
		fmt.Printf("error: %s\n", err.Error())
	} else {
		fmt.Printf("result:\n%s\n", Format(dish.Get()))
	}
}

func Do(name string) {
	fmt.Println(name)
	<-time.After(time.Second + time.Duration(rand.Int63n(int64(time.Millisecond*200))))
}

func Format(deps ...any) string {
	d := []string{}
	for _, x := range deps {
		d = append(d, fmt.Sprint(x))
	}
	return "  " + strings.ReplaceAll(strings.Join(d, "\n"), "\n", "\n  ")
}

type BoilingWater string

func BoilWater(ctx context.Context, eg *errgroup.Group) *future.Future[BoilingWater] {
	bw := &future.Future[BoilingWater]{}
	eg.Go(func() error {
		Do("preparing boiling water")
		bw.Resolve(BoilingWater("BoilingWater"))
		return nil
	})
	return bw
}

type RawSpaghetti string

func GetSpaghetti(ctx context.Context, eg *errgroup.Group) *future.Future[RawSpaghetti] {
	bw := &future.Future[RawSpaghetti]{}
	eg.Go(func() error {
		Do("getting raw spaghetti")
		bw.Resolve(RawSpaghetti("RawSpaghetti"))
		return nil
	})
	return bw
}

type Tomatoes string

func GetTomatoes(ctx context.Context, eg *errgroup.Group) *future.Future[Tomatoes] {
	bw := &future.Future[Tomatoes]{}
	eg.Go(func() error {
		Do("getting tomatoes")
		bw.Resolve(Tomatoes("Tomatoes"))
		return nil
	})
	return bw
}

type SlicingBoard string

func GetSlicingBoard(ctx context.Context, eg *errgroup.Group) *future.Future[SlicingBoard] {
	bw := &future.Future[SlicingBoard]{}
	eg.Go(func() error {
		Do("getting slicing board")
		bw.Resolve(SlicingBoard("SlicingBoard"))
		return nil
	})
	return bw
}

type Onion string

func GetOnion(ctx context.Context, eg *errgroup.Group) *future.Future[Onion] {
	bw := &future.Future[Onion]{}
	eg.Go(func() error {
		Do("getting onion")
		bw.Resolve(Onion("Onion"))
		return nil
	})
	return bw
}

type Garlic string

func GetGarlic(ctx context.Context, eg *errgroup.Group) *future.Future[Garlic] {
	bw := &future.Future[Garlic]{}
	eg.Go(func() error {
		Do("getting garlic")
		bw.Resolve(Garlic("Garlic"))
		return nil
	})
	return bw
}

type Grater string

func GetGrater(ctx context.Context, eg *errgroup.Group) *future.Future[Grater] {
	bw := &future.Future[Grater]{}
	eg.Go(func() error {
		Do("getting grater")
		bw.Resolve(Grater("Grater"))
		return nil
	})
	return bw
}

type CookedSpaghetti string

func CookSpaghetti(ctx context.Context, eg *errgroup.Group, bw *future.Future[BoilingWater], rs *future.Future[RawSpaghetti]) *future.Future[CookedSpaghetti] {
	cs := &future.Future[CookedSpaghetti]{}
	eg.Go(func() error {
		if !future.Await(ctx, bw.Done(), rs.Done()) {
			return nil
		}
		Do("cooking spaghetti")
		cs.Resolve(CookedSpaghetti("CookedSpaghetti:\n" + Format(bw.Get(), rs.Get())))
		return nil
	})
	return cs
}

type ChoppedTomatoes string
type ChoppedOnion string

func ChopVegetables(ctx context.Context, eg *errgroup.Group, sb *future.Future[SlicingBoard], t *future.Future[Tomatoes], o *future.Future[Onion]) (*future.Future[ChoppedTomatoes], *future.Future[ChoppedOnion]) {
	ct, co := &future.Future[ChoppedTomatoes]{}, &future.Future[ChoppedOnion]{}
	eg.Go(func() error {
		if !future.Await(ctx, sb.Done(), t.Done(), o.Done()) {
			return nil
		}
		Do("chopping tomatoes and onion")
		ct.Resolve(ChoppedTomatoes("ChoppedTomatoes:\n" + Format(sb.Get(), t.Get())))
		co.Resolve(ChoppedOnion("ChoppedOnion:\n" + Format(sb.Get(), o.Get())))
		return nil
	})
	return ct, co
}

type GratedGarlic string

func GrateGarlic(ctx context.Context, eg *errgroup.Group, gr *future.Future[Grater], ga *future.Future[Garlic]) *future.Future[GratedGarlic] {
	gg := &future.Future[GratedGarlic]{}
	eg.Go(func() error {
		if !future.Await(ctx, gr.Done(), ga.Done()) {
			return nil
		}
		Do("grating garlic")
		gg.Resolve(GratedGarlic("GratedGarlic:\n" + Format(gr.Get(), ga.Get())))
		return nil
	})
	return gg
}

type CookedVegetables string

func CookVegetables(ctx context.Context, eg *errgroup.Group, bw *future.Future[BoilingWater], ct *future.Future[ChoppedTomatoes], co *future.Future[ChoppedOnion], gg *future.Future[GratedGarlic]) *future.Future[CookedVegetables] {
	cv := &future.Future[CookedVegetables]{}
	eg.Go(func() error {
		if !future.Await(ctx, bw.Done(), ct.Done(), co.Done(), gg.Done()) {
			return nil
		}
		Do("cooking vegetables")
		cv.Resolve(CookedVegetables("CookedVegetables:\n" + Format(bw.Get(), ct.Get(), co.Get(), gg.Get())))
		return nil
	})
	return cv
}

type Dish string

func PutOnPlate(ctx context.Context, eg *errgroup.Group, cs *future.Future[CookedSpaghetti], cv *future.Future[CookedVegetables]) *future.Future[Dish] {
	d := &future.Future[Dish]{}
	eg.Go(func() error {
		if !future.Await(ctx, cs.Done(), cv.Done()) {
			return nil
		}
		Do("putting everything on plate")
		d.Resolve(Dish("Dish:\n" + Format(cs.Get(), cv.Get())))
		return nil
	})
	return d
}
