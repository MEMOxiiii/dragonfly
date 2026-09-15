package block_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

// TestTorchBreaksWithoutSupport verifies that a torch is broken by a neighbour
// update on the tick after its supporting block is removed, using a
// synchronous World to make the tick deterministic.
func TestTorchBreaksWithoutSupport(t *testing.T) {
	w := world.Config{Synchronous: true, Entities: entity.DefaultRegistry}.New()
	defer w.Close()

	support, torch := cube.Pos{0, 0, 0}, cube.Pos{0, 1, 0}
	w.Do(func(tx *world.Tx) {
		tx.SetBlock(support, block.Stone{}, nil)
		tx.SetBlock(torch, block.Torch{Facing: cube.FaceDown}, nil)
		tx.SetBlock(support, block.Air{}, nil)
	})
	w.AdvanceTick()

	b, err := world.Call(context.Background(), w, func(tx *world.Tx) (world.Block, error) {
		return tx.Block(torch), nil
	})
	if err != nil {
		t.Fatalf("read torch block: %v", err)
	}
	if b != (block.Air{}) {
		t.Errorf("expected torch to break after removing its support, got %v", b)
	}
}

// TestNyliumDecaysUnderOpaqueBlock verifies that nylium reverts to netherrack on a random tick once a block
// that blocks all light is placed above it, and that it survives under a light-passing block.
func TestNyliumDecaysUnderOpaqueBlock(t *testing.T) {
	w := world.Config{Synchronous: true, Entities: entity.DefaultRegistry}.New()
	defer w.Close()

	covered, open := cube.Pos{0, 0, 0}, cube.Pos{4, 0, 0}
	res, err := world.Call(context.Background(), w, func(tx *world.Tx) ([2]world.Block, error) {
		r := rand.New(rand.NewPCG(0, 0))
		for _, pos := range []cube.Pos{covered, open} {
			tx.SetBlock(pos, block.Nylium{}, nil)
		}
		tx.SetBlock(covered.Side(cube.FaceUp), block.Stone{}, nil)
		tx.SetBlock(open.Side(cube.FaceUp), block.Glass{}, nil)

		block.Nylium{}.RandomTick(covered, tx, r)
		block.Nylium{}.RandomTick(open, tx, r)
		return [2]world.Block{tx.Block(covered), tx.Block(open)}, nil
	})
	if err != nil {
		t.Fatalf("tick nylium: %v", err)
	}
	if res[0] != (block.Netherrack{}) {
		t.Errorf("expected nylium under stone to decay into netherrack, got %v", res[0])
	}
	if res[1] != (block.Nylium{}) {
		t.Errorf("expected nylium under glass to survive, got %v", res[1])
	}
}

// TestBoneMealOnNetherrackSpreadsNylium verifies that bone meal turns netherrack neighbouring nylium into
// nylium of that same type, and leaves netherrack without a nylium neighbour untouched.
func TestBoneMealOnNetherrackSpreadsNylium(t *testing.T) {
	w := world.Config{Synchronous: true, Entities: entity.DefaultRegistry}.New()
	defer w.Close()

	neighboured, isolated := cube.Pos{0, 0, 0}, cube.Pos{8, 0, 0}
	res, err := world.Call(context.Background(), w, func(tx *world.Tx) ([2]world.Block, error) {
		tx.SetBlock(neighboured, block.Netherrack{}, nil)
		tx.SetBlock(neighboured.Side(cube.FaceNorth), block.Nylium{Warped: true}, nil)
		tx.SetBlock(isolated, block.Netherrack{}, nil)

		if r := (block.Netherrack{}).BoneMeal(neighboured, tx); r != item.BoneMealResultSmall {
			return [2]world.Block{}, fmt.Errorf("expected a small bone meal result next to nylium, got %v", r)
		}
		if r := (block.Netherrack{}).BoneMeal(isolated, tx); r != item.BoneMealResultNone {
			return [2]world.Block{}, fmt.Errorf("expected no bone meal result without nylium, got %v", r)
		}
		return [2]world.Block{tx.Block(neighboured), tx.Block(isolated)}, nil
	})
	if err != nil {
		t.Fatalf("bone meal netherrack: %v", err)
	}
	if res[0] != (block.Nylium{Warped: true}) {
		t.Errorf("expected netherrack to convert into warped nylium, got %v", res[0])
	}
	if res[1] != (block.Netherrack{}) {
		t.Errorf("expected isolated netherrack to stay netherrack, got %v", res[1])
	}
}

// TestMyceliumSpreadsOntoDirt verifies that mycelium spreads onto neighbouring dirt on a random tick, but
// never onto coarse dirt.
func TestMyceliumSpreadsOntoDirt(t *testing.T) {
	w := world.Config{Synchronous: true, Entities: entity.DefaultRegistry}.New()
	defer w.Close()

	source, dirt, coarse := cube.Pos{0, 0, 0}, cube.Pos{1, 0, 0}, cube.Pos{-1, 0, 0}
	res, err := world.Call(context.Background(), w, func(tx *world.Tx) ([2]world.Block, error) {
		tx.SetBlock(source, block.Mycelium{}, nil)
		tx.SetBlock(dirt, block.Dirt{}, nil)
		tx.SetBlock(coarse, block.Dirt{Coarse: true}, nil)

		for i := range 64 {
			block.Mycelium{}.RandomTick(source, tx, rand.New(rand.NewPCG(uint64(i), 0)))
		}
		return [2]world.Block{tx.Block(dirt), tx.Block(coarse)}, nil
	})
	if err != nil {
		t.Fatalf("tick mycelium: %v", err)
	}
	if res[0] != (block.Mycelium{}) {
		t.Errorf("expected mycelium to spread onto adjacent dirt, got %v", res[0])
	}
	if res[1] != (block.Dirt{Coarse: true}) {
		t.Errorf("expected mycelium not to spread onto coarse dirt, got %v", res[1])
	}
}

// TestRootedDirtToolInteractions verifies that a hoe strips rooted dirt into dirt rather than tilling it into
// farmland, and that a shovel turns it into a dirt path.
func TestRootedDirtToolInteractions(t *testing.T) {
	tilled, ok := (block.RootedDirt{}).Till()
	if !ok || tilled != (block.Dirt{}) {
		t.Errorf("expected a hoe to turn rooted dirt into dirt, got %v (%v)", tilled, ok)
	}
	shovelled, ok := (block.RootedDirt{}).Shovel()
	if !ok || shovelled != (block.DirtPath{}) {
		t.Errorf("expected a shovel to turn rooted dirt into a dirt path, got %v (%v)", shovelled, ok)
	}
}
