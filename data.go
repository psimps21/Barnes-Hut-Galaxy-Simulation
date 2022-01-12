package main

import (
	"math"
)

// Universe contains a slice of pointers to stars and a width parameter.
// We conceptualize the universe as a square -- stars may go outside the universe
// but the width dictates relative distances when drawing the universe.
type Universe struct {
	stars   []*Star
	width   float64
	uniQuad *Quadrant
}

// InitializeUniQuad
func (u *Universe) InitializeUniQuad() {
	u.uniQuad = &Quadrant{u.width, u.width, u.width, u.stars, &Node{parent: nil}}
	u.uniQuad.node.sector = u.uniQuad
}

// Galaxy is a potentially useful object holding a list of star positions
type Galaxy []*Star

// Star is analogous to the "Body" object from the jupiter simulations.
type Star struct {
	position, velocity, acceleration OrderedPair
	mass                             float64
	radius                           float64
	red, blue, green                 uint8
}

//OrderedPair represents a point or vector.
type OrderedPair struct {
	x float64
	y float64
}

//QuadTree simply contains a pointer to the root.
//Another way of doing this would be type QuadTree *Node
type QuadTree struct {
	root   *Node
	leaves []*Node
}

//Node object contains a slice of children (this could just as easily be an array of length 4).
//A node refers to a star. Sometimes, the star will be a "dummy" star, sometimes it is a star in the
//universe, and sometimes it is nil. Every internal node points to a dummy star.
type Node struct {
	children     []*Node
	star         *Star
	sector       *Quadrant
	parent       *Node
	pathFromRoot []*Node
}

//Quadrant is an object representing a sub-square within a larger universe.
type Quadrant struct {
	x     float64 //bottom left corner x coordinate
	y     float64 //bottom right corner y coordinate
	width float64
	stars []*Star
	node  *Node
}

// Compute the Euclidian Distance between two bodies
func Dist(s1, s2 Star) float64 {
	dx := s1.position.x - s2.position.x
	dy := s1.position.y - s2.position.y
	return math.Sqrt(dx*dx + dy*dy)
}

// ComputeGravityForce computes the gravity force between star 1 and star 2.
func ComputeGravityForce(s1, s2 *Star) OrderedPair {
	d := Dist(*s1, *s2)
	deltaX := s2.position.x - s1.position.x
	deltaY := s2.position.y - s1.position.y
	F := G * s1.mass * s2.mass / (d * d)

	return OrderedPair{
		x: F * deltaX / d,
		y: F * deltaY / d,
	}
}

//ForceFromTree recursively computes the force acting on a star s for all stars in a quad tree
func ForceFromTree(n, root *Node, theta float64, netForce OrderedPair) OrderedPair {
	// Base Cases
	if root.star == n.star { // If computing force on self do nothing
		return netForce
	} else if root.children == nil { // if leaf node add force from that star to net force
		f := ComputeGravityForce(n.star, root.star)
		netForce.Add(f)
		return netForce
	}

	// if internal node
	s := MostRecentAncestor(n.pathFromRoot, root.pathFromRoot).sector.width
	d := Dist(*n.star, *root.star)

	if s/d < theta { // compute force with dummy star at node r
		f := ComputeGravityForce(n.star, root.star)
		netForce.Add(f)
		return netForce
	}
	// If s/d >= theta compute force of the children of r
	var childrenForce OrderedPair
	for _, child := range root.children {
		if child != nil {
			childrenForce.Add(ForceFromTree(n, child, theta, netForce))
		}
	}

	netForce.Add(childrenForce)
	return netForce
}

//ComputeNetForce computes the net force acting on a star s given a quadtree
func (n *Node) ComputeNetForce(qTree QuadTree, theta float64) OrderedPair {
	return ForceFromTree(n, qTree.root, theta, OrderedPair{0.0, 0.0})
}

// NewVelocity makes the velocity of this object consistent with the acceleration.
func (s *Star) NewVelocity(t float64) OrderedPair {
	return OrderedPair{
		x: s.velocity.x + s.acceleration.x*t,
		y: s.velocity.y + s.acceleration.y*t,
	}
}

// NewPosition computes the new poosition given the updated acc and velocity.
//
// Assumputions: constant acceleration over a time step.
// => DeltaX = v_avg * t
//    DeltaX = (v_start + v_final)*t/ 2
// because v_final = v_start + acc*t:
//	  DeltaX = (v_start + v_start + acc*t)t/2
// Simplify:
//	DeltaX = v_start*t + 0.5acc*t*t
// =>
//  NewX = v_start*t + 0.5acc*t*t + OldX
//
func (s *Star) NewPosition(t float64) OrderedPair {
	return OrderedPair{
		x: s.position.x + s.velocity.x*t + 0.5*s.acceleration.x*t*t,
		y: s.position.y + s.velocity.y*t + 0.5*s.acceleration.y*t*t,
	}
}

// UpdateAccel computes the new accerlation vector for b
func (n *Node) NewAccel(qt QuadTree, theta float64) OrderedPair {
	F := n.ComputeNetForce(qt, theta)
	return OrderedPair{
		x: F.x / n.star.mass,
		y: F.y / n.star.mass,
	}
}

//UpdateStar returns a star with updated position, velocity, and acceleration
func UpdateStar(node *Node, qt QuadTree, theta, t float64) *Star {
	newStar := CopyStar(node.star)
	newStar.acceleration = node.NewAccel(qt, theta)
	newStar.velocity = newStar.NewVelocity(t)
	newStar.position = newStar.NewPosition(t)
	return newStar
}

//StarsFromNodes returns a list of star pointers from a given list of nodes
func StarsFromNodes(nodes []*Node) []*Star {
	var stars []*Star
	for _, node := range nodes {
		if node != nil {
			stars = append(stars, node.star)
		}
	}
	return stars
}

//ComputeCenterOfGravity computes the center of gravity for a list of stars
func ComputeCenterOfGravity(stars []*Star) OrderedPair {
	var totalMass, weightedSumX, weightedSumY float64
	for _, star := range stars {
		totalMass += star.mass
		weightedSumX += star.mass * star.position.x
		weightedSumY += star.mass * star.position.y
	}

	return OrderedPair{
		x: weightedSumX / totalMass,
		y: weightedSumY / totalMass,
	}
}

//SumStarMasses sums the mass of stars
func SumStarMasses(stars ...*Star) float64 {
	var totalSum float64
	for _, star := range stars {
		totalSum += star.mass
	}
	return totalSum
}

//MostRecentAncestor returns the most recent ancestor given two paths from a node in the tree to root
// should never return nil because all nodes are in the same tree i.e. root is default common ancestor
func MostRecentAncestor(p1, p2 []*Node) *Node {
	for i := len(p1) - 1; i >= 0; i-- {
		for j := len(p2) - 1; j >= 0; j-- {
			if p1[i] == p2[j] {
				return p1[i]
			}
		}
	}
	return nil
}

//
/* Quadrant Methods */
//

//GenerateSubQuadrants returns a list of SubQuadrants for a given quadrant. In order [NW=0,NE=1,SW=2,SE=3]
// If quadrant has no stars that quadrant's node is nil
func (q *Quadrant) GenerateSubQuadrants() []*Quadrant {
	newQuads := make([]*Quadrant, 4)
	subXY := [][]float64{ // x,y coords for each subquadrant
		[]float64{q.x - (q.width / 2), q.y - (q.width / 2)}, //top left corner
		[]float64{q.x - (q.width / 2), q.y},                 // top right
		[]float64{q.x, q.y - (q.width / 2)},                 // bottom left
		[]float64{q.x, q.y},
	}

	// Populate the list of star for each sub quadrant
	for _, star := range q.stars {
		starSubQ := star.FindStarSubQuad(q)
		if starSubQ > -1 { // If star is in valid quadrant
			if newQuads[starSubQ] == nil { // If subquadrant not initialized
				var subQuad Quadrant
				subQuad.x, subQuad.y, subQuad.width = subXY[starSubQ][0], subXY[starSubQ][1], q.width/2
				subQuad.node = &Node{parent: q.node, sector: &subQuad}
				newQuads[starSubQ] = &subQuad
			}
			newQuads[starSubQ].stars = append(newQuads[starSubQ].stars, star)
		}
	}

	return newQuads
}

//MakeQuadTree recursively returns the root node of a QuadTree
func (q *Quadrant) MakeQuadTree() *Node {
	// Base Cases
	if len(q.stars) == 1 { // Quadrant contains only one star
		q.node.star = q.stars[0]
		return q.node
	}

	var dummyStar Star

	// Recursively call MakeQuadTree for each subQuadrant
	subQs := q.GenerateSubQuadrants()
	children := make([]*Node, 4)
	for i, subQ := range subQs {
		if subQ != nil {
			children[i] = subQ.MakeQuadTree()
		}
	}

	// Update values in dummyStar
	stars := StarsFromNodes(children)
	dummyStar.position, dummyStar.mass = ComputeCenterOfGravity(stars), SumStarMasses(stars...)

	// Update root node with children and star
	q.node.children = children
	q.node.star = &dummyStar
	return q.node
}

//
/* Star Methods */
//

//FindStarSubQuad Given a quadrant will return the number of the subquadrant the star resides in
func (s *Star) FindStarSubQuad(q *Quadrant) int {
	midX, midY := q.x-(q.width/2), q.y-(q.width/2)
	if q.x-q.width <= s.position.x && s.position.x < midX { // s is in NW=0 or NE=2
		if q.y-q.width <= s.position.y && s.position.y < midY { // s is in NW=0
			return 0
		} else if midY <= s.position.y && s.position.y <= q.y { // s is in NE= 1
			return 1
		}
	} else if midX <= s.position.x && s.position.x <= q.x { // s is in SW=1 or SE=3
		if q.y-q.width <= s.position.y && s.position.y < midY { // s is in SW=2
			return 2
		} else if midY <= s.position.y && s.position.y <= q.y { // s is in SE=3
			return 3
		}
	}
	return -1
}

//
/* Node Methods */
//

// SetPathsToRoot sets the path to the given node for all nodes under it and returns list of leaves in tree
func (n *Node) SetPathsToRoot(rootPath, leaves []*Node) []*Node {
	if len(n.children) == 0 { // If node is a leaf
		n.pathFromRoot = append(rootPath, n)
		return []*Node{n}
	}

	rootPath = append(rootPath, n)
	n.pathFromRoot = rootPath
	var childLeaves []*Node
	for _, child := range n.children {
		if child != nil {
			childLeaves = append(childLeaves, child.SetPathsToRoot(rootPath, leaves)...)
		}
	}
	return append(leaves, childLeaves...)
}

//
/* Ordered Pair Methods */
//

// Add adds the value of a given ordered pair to a current ordered pair
func (v *OrderedPair) Add(v2 OrderedPair) {
	v.x += v2.x
	v.y += v2.y
}

//
/* QuadTree Methods */
//

func (qt *QuadTree) SetRootPathsAndLeaves() {
	qt.leaves = qt.root.SetPathsToRoot([]*Node{}, []*Node{})
}
