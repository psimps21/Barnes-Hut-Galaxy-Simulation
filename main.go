package main

import (
	"fmt"
	"gifhelper"
	"math"
	"os"
	"time"
)

func main() {

	filename := os.Args[1]

	var (
		initialUniverse                 *Universe
		numGens, canvasWidth, frequency int
		step, theta, scalingFactor      float64
	)

	switch filename {
	case "jupiter":
		const G = 2.0 * 6.67408e-11
		width := 4e9

		var g Galaxy

		// Io
		var jupiter, io, europa, ganymede, callisto Star

		jupiter.mass = 1.898e27
		jupiter.red, jupiter.green, jupiter.blue = 223, 227, 202
		jupiter.position.x, jupiter.position.y = 2000000000, 2000000000
		jupiter.blue = 255
		jupiter.radius = 71000000
		jupiter.velocity.x, jupiter.velocity.y = 0, 0

		io.red, io.green, io.blue = 249, 249, 165
		io.mass = 8.9319e22
		io.radius = 1821000 * 6
		io.position.x, io.position.y = 2e9-4.216e8, 2e9
		io.radius = 1.821e6
		io.velocity.x, io.velocity.y = 0, -17320

		europa.red, europa.green, europa.blue = 132, 83, 52
		europa.mass = 4.7998e22
		europa.radius = 1569000 * 6
		europa.position.x, europa.position.y = 2e9, 2e9+6.709e8
		europa.velocity.x, europa.velocity.y = -13740, 0

		ganymede.red, ganymede.green, ganymede.blue = 76, 0, 153
		ganymede.mass = 1.4819 * math.Pow(10, 23)
		ganymede.radius = 2631000 * 6
		ganymede.position.x, ganymede.position.y = 2000000000+1070400000, 2000000000
		ganymede.velocity.x, ganymede.velocity.y = 0, 10870

		callisto.red, callisto.green, callisto.blue = 0, 153, 76
		callisto.mass = 1.0759 * math.Pow(10, 23)
		callisto.radius = 2410000 * 6
		callisto.position.x, callisto.position.y = 2000000000, 2000000000-1882700000
		callisto.velocity.x, callisto.velocity.y = 8200, 0

		g = append(g, &io)
		g = append(g, &europa)
		g = append(g, &ganymede)
		g = append(g, &callisto)
		g = append(g, &jupiter)

		initialUniverse = InitializeUniverse([]Galaxy{g}, width)
		numGens = 10000
		step = 75
		canvasWidth = 100
		scalingFactor = 1 // a scaling factor is needed to inflate size of stars when drawn because galaxies are very sparse
		frequency = 50

	case "galaxy":
		g0 := InitializeGalaxy0(500, 5.5e21, 6e22, 4e22)
		width := 1.0e23
		galaxies := []Galaxy{g0} //, g1}

		initialUniverse = InitializeUniverse(galaxies, width)

		// now evolve the universe: feel free to adjust the following parameters.
		numGens = 1700
		step = 1.85e15
		theta = 0.5
		canvasWidth = 1000
		scalingFactor = 1e11 // a scaling factor is needed to inflate size of stars when drawn because galaxies are very sparse
		frequency = 50

	case "collision":
		g0 := InitializeGalaxy0(500, 9e21, 6e22, 4e22)
		g1 := InitializeGalaxy1(500, 9e21, 4.5e22, 4e22)

		width := 1.0e23
		galaxies := []Galaxy{g0, g1}

		initialUniverse = InitializeUniverse(galaxies, width)

		// now evolve the universe: feel free to adjust the following parameters.
		numGens = 1000
		step = 1.8e15
		theta = 0.5
		canvasWidth = 1000
		scalingFactor = 1e11 // a scaling factor is needed to inflate size of stars when drawn because galaxies are very sparse
		frequency = 10
	}

	startTime := time.Now()
	timePoints := BarnesHut(initialUniverse, numGens, step, theta)
	elapTime := time.Since(startTime)
	fmt.Println("\nTime to run barnes hut for", numGens, "generations:", elapTime, "\n")

	fmt.Println("Simulation run. Now drawing images.")

	imageList := AnimateSystem(timePoints, canvasWidth, frequency, scalingFactor)

	fmt.Println("Images drawn. Now generating GIF.")
	gifhelper.ImagesToGIF(imageList, filename)
	fmt.Println("GIF drawn.")
}
