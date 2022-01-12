package main

import "fmt"

//BarnesHut is our highest level function.
//Input: initial Universe object, a number of generations, and a time interval.
//Output: collection of Universe objects corresponding to updating the system
//over indicated number of generations every given time interval.
func BarnesHut(initialUniverse *Universe, numGens int, time, theta float64) []*Universe {
	timePoints := make([]*Universe, numGens+1)
	timePoints[0] = initialUniverse
	// Your code goes here. Use subroutines! :)
	for gen := 0; gen < numGens; gen++ {
		timePoints[gen+1] = UpdateUniverse(*timePoints[gen], theta, time)
		fmt.Println("Finished Universe Update", gen+1)
	}

	return timePoints
}

// UpdateUniverse returns a new universe after time t
func UpdateUniverse(univ Universe, theta, t float64) *Universe {
	newUniverse := CopyUniverse(&univ)
	uniQuadTree := QuadTree{
		root: univ.uniQuad.MakeQuadTree(),
	}
	uniQuadTree.SetRootPathsAndLeaves()

	// Compute force on each star using QuadTree
	newStars := make([]*Star, 0, len(univ.stars))
	for _, node := range uniQuadTree.leaves {
		newStars = append(newStars, UpdateStar(node, uniQuadTree, theta, t))
	}

	// Set list of updates stars to list of stars for newUniverse and reset the quadrant for new universe
	newUniverse.stars = newStars
	newUniverse.InitializeUniQuad()

	return newUniverse
}
