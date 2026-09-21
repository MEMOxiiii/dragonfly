package entity

import (
	"container/heap"
	"math"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

const (
	maxPathNodes = 384
	maxPathDrop  = 3
)

// pathNodeLimit scales the number of nodes searched with the distance to the
// target, so that a search for a target nearby gives up early.
func pathNodeLimit(start, target cube.Pos) int {
	return min(int(distance(start, target))*12+32, maxPathNodes)
}

type navigator struct {
	path  []cube.Pos
	index int
	speed float64

	stuckPos   mgl64.Vec3
	stuckTicks int

	open   pathHeap
	nodes  map[cube.Pos]*pathNode
	closed map[cube.Pos]struct{}
}

func (n *navigator) navigating() bool {
	return n.index < len(n.path)
}

func (n *navigator) stop() {
	n.path, n.index, n.stuckTicks = n.path[:0], 0, 0
}

// moveTo searches a path to the target and starts following it, reporting
// whether one was found.
func (n *navigator) moveTo(m *Mob, target cube.Pos, speed float64) bool {
	n.path = n.findPath(m, cube.PosFromVec3(m.Position()), target)
	if len(n.path) == 0 {
		return false
	}
	n.index, n.speed = 0, speed
	n.stuckPos, n.stuckTicks = m.Position(), 0
	return true
}

func (n *navigator) tick(m *Mob) {
	if !n.navigating() {
		return
	}
	pos := m.Position()
	node := n.path[n.index].Vec3Middle()
	if math.Abs(node[0]-pos[0]) < 0.35 && math.Abs(node[2]-pos[2]) < 0.35 && math.Abs(node[1]-pos[1]) < 1 {
		if n.index++; !n.navigating() {
			n.stop()
			return
		}
		node = n.path[n.index].Vec3Middle()
	}
	if pos.Sub(n.stuckPos).Len() < 0.05 {
		if n.stuckTicks++; n.stuckTicks > 40 {
			n.stop()
			return
		}
	} else {
		n.stuckPos, n.stuckTicks = pos, 0
	}

	m.Move(mgl64.Vec3{node[0] - pos[0], 0, node[2] - pos[2]}, n.speed)
	m.LookAt(mgl64.Vec3{node[0], pos[1] + m.eyeHeight(), node[2]})
	if node[1] > pos[1]+0.1 || m.mob().horizontallyCollided {
		m.Jump()
	}
}

type pathNode struct {
	pos    cube.Pos
	parent *pathNode
	g, f   float64
	index  int
}

type pathHeap []*pathNode

func (h pathHeap) Len() int           { return len(h) }
func (h pathHeap) Less(i, j int) bool { return h[i].f < h[j].f }
func (h pathHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i]; h[i].index, h[j].index = i, j }
func (h *pathHeap) Push(x any)        { n := x.(*pathNode); n.index = len(*h); *h = append(*h, n) }
func (h *pathHeap) Pop() (x any)      { old := *h; x, *h = old[len(old)-1], old[:len(old)-1]; return x }
func (h pathHeap) update(n *pathNode) { heap.Fix(&h, n.index) }

func distance(a, b cube.Pos) float64 { return a.Vec3().Sub(b.Vec3()).Len() }

var pathOffsets = [8][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}

func (n *navigator) findPath(m *Mob, start, target cube.Pos) []cube.Pos {
	start = groundPos(m, start)
	if start == target || !walkable(m, target) {
		if t, ok := nearestWalkable(m, target); ok {
			target = t
		} else {
			return nil
		}
	}
	if start == target {
		return nil
	}
	if n.nodes == nil {
		n.nodes, n.closed = make(map[cube.Pos]*pathNode, maxPathNodes), make(map[cube.Pos]struct{}, maxPathNodes)
		n.open = make(pathHeap, 0, maxPathNodes)
	}
	clear(n.nodes)
	clear(n.closed)
	open, nodes, closed := &n.open, n.nodes, n.closed
	*open = (*open)[:0]
	nodes[start] = &pathNode{pos: start, f: distance(start, target)}
	heap.Push(open, nodes[start])

	var best *pathNode
	limit := pathNodeLimit(start, target)
	for open.Len() > 0 && len(closed) < limit {
		current := heap.Pop(open).(*pathNode)
		if best == nil || distance(current.pos, target) < distance(best.pos, target) {
			best = current
		}
		if current.pos == target {
			return tracePath(current, n.path[:0])
		}
		closed[current.pos] = struct{}{}

		for _, o := range pathOffsets {
			next, ok := step(m, current.pos, o[0], o[1])
			if !ok {
				continue
			}
			if _, done := closed[next]; done {
				continue
			}
			g := current.g + distance(current.pos, next)
			n, ok := nodes[next]
			if !ok {
				n = &pathNode{pos: next, parent: current, g: g, f: g + distance(next, target)}
				nodes[next] = n
				heap.Push(open, n)
				continue
			}
			if g < n.g {
				n.parent, n.g, n.f = current, g, g+distance(next, target)
				open.update(n)
			}
		}
	}
	if best == nil || best.pos == start {
		return nil
	}
	return tracePath(best, n.path[:0])
}

func tracePath(n *pathNode, path []cube.Pos) []cube.Pos {
	for ; n.parent != nil; n = n.parent {
		path = append(path, n.pos)
	}
	slicesReverse(path)
	return path
}

func slicesReverse(p []cube.Pos) {
	for i, j := 0, len(p)-1; i < j; i, j = i+1, j-1 {
		p[i], p[j] = p[j], p[i]
	}
}

// step returns the position a mob at pos ends up in when moving by the offset
// passed, walking up a block or falling down at most maxPathDrop blocks.
func step(m *Mob, pos cube.Pos, dx, dz int) (cube.Pos, bool) {
	next := pos.Add(cube.Pos{dx, 0, dz})
	if dx != 0 && dz != 0 {
		if !passable(m, pos.Add(cube.Pos{dx, 0, 0})) || !passable(m, pos.Add(cube.Pos{0, 0, dz})) {
			return cube.Pos{}, false
		}
	}
	if walkable(m, next) {
		return next, true
	}
	if up := next.Add(cube.Pos{0, 1, 0}); passable(m, pos.Add(cube.Pos{0, 1, 0})) && walkable(m, up) {
		return up, true
	}
	if !passable(m, next) {
		return cube.Pos{}, false
	}
	for i := 1; i <= maxPathDrop; i++ {
		down := next.Sub(cube.Pos{0, i, 0})
		if walkable(m, down) {
			return down, true
		}
		if !passable(m, down) {
			break
		}
	}
	return cube.Pos{}, false
}

func groundPos(m *Mob, pos cube.Pos) cube.Pos {
	for i := 0; i < 4 && !solidBelow(m.tx, pos); i++ {
		pos = pos.Sub(cube.Pos{0, 1, 0})
	}
	return pos
}

func nearestWalkable(m *Mob, pos cube.Pos) (cube.Pos, bool) {
	for dy := 0; dy <= 2; dy++ {
		if up := pos.Add(cube.Pos{0, dy, 0}); walkable(m, up) {
			return up, true
		}
		if down := pos.Sub(cube.Pos{0, dy, 0}); walkable(m, down) {
			return down, true
		}
	}
	return cube.Pos{}, false
}

func walkable(m *Mob, pos cube.Pos) bool {
	return solidBelow(m.tx, pos) && passable(m, pos)
}

func solidBelow(tx *world.Tx, pos cube.Pos) bool {
	below := pos.Side(cube.FaceDown)
	boxes := tx.Block(below).Model().BBox(below, tx)
	for _, box := range boxes {
		if box.Max()[1] >= 0.9 {
			return true
		}
	}
	return false
}

func passable(m *Mob, pos cube.Pos) bool {
	for y := 0; y < m.mob().pathHeight; y++ {
		at := pos.Add(cube.Pos{0, y, 0})
		b := m.tx.Block(at)
		if len(b.Model().BBox(at, m.tx)) != 0 {
			return false
		}
	}
	return true
}
