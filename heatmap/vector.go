package main

import (
	"math"
)

func (a *App) GetRayCast(tick int) []Vector {

	var vectorRays []Vector

	for _, player := range a.demo.gameTicks[tick].players {

		playerVector := a.vectorFromPlayer(player)

		newPlayerVector, _, hasIntersection := a.findShortestBoundary(playerVector)

		//invert the coordinates back for html canvas
		newPlayerVector.A.Y = -newPlayerVector.A.Y
		newPlayerVector.B.Y = -newPlayerVector.B.Y

		if hasIntersection {
			vectorRays = append(vectorRays, newPlayerVector)
		}
	}
	return vectorRays
}

func (a *App) CheckControlVectors(tick int) bool {

	checkNextVector:
	for _, vector := range a.demo.controlVectors {

		controlState := NoTeams

		//for each player
		for _, player := range a.demo.gameTicks[tick].players {
			//fmt.Println("Checking player ", player.Team)

			if player.Team == 0 || player.Team == 1 {
				continue
			}
			
			playerVector := a.vectorFromPlayer(player)

			intersection, intersects := intersectPoint(playerVector.A, playerVector.B, vector.vec.A, vector.vec.B)

			if intersects {
				controlDistX := intersection.X - playerVector.A.X
				controlDistY := intersection.Y - playerVector.A.Y

				controlDistance := math.Sqrt(float64(controlDistX*controlDistX + controlDistY*controlDistY))


				_, boundaryDist, hasIntersection := a.findShortestBoundary(playerVector)

				if hasIntersection && boundaryDist < controlDistance {
					continue
				}

				if controlState == TeamCounterTerrorists && player.Team == TeamTerrorists {
					if vector.state != AllTeams {
						return false
					}
					continue checkNextVector
				}

				if controlState == TeamTerrorists && player.Team == TeamCounterTerrorists {
					if vector.state != AllTeams {
						return false
					}
					continue checkNextVector
				}

				controlState = player.Team

			}
		}

		if vector.state != controlState {
			return false
		}

	}
	//all vectors have been checked and they match, so we can return true
	return true
}

func (a *App) GetViewDirections(tick int) []Vector {

	var directionRays []Vector

	for _, player := range a.demo.gameTicks[tick].players {

		playerVector := a.vectorFromPlayer(player)

		directionRays = append(directionRays, playerVector)

	}
	return directionRays
}

//DONT FORGET TO HAVE A STATE PROVIDED
func (a *App) DefineControlVector(vector Vector, state byte)  {
	a.demo.controlVectors = append(a.demo.controlVectors, controlVector{Vector{Point{vector.A.X, -vector.A.Y}, Point{vector.B.X, -vector.B.Y}}, state})
}

//returns the shortest intersection found (vector it intersects with, the distance), if it finds no intersection, returns false for safety.
func (a *App) findShortestBoundary(vec Vector) (Vector, float64, bool) {

	var shortestDistance float64
	var shortestVector Vector

	intersectionFound := false

	//check through all boundary vectors
	for _, boundaryVector := range a.demo.boundaries {

		intersection, intersects := intersectPoint(vec.A, vec.B, boundaryVector.A, boundaryVector.B)


		if intersects {
			distY := intersection.X - vec.A.X
			distX := intersection.Y - vec.A.Y

			dist := math.Sqrt(float64(distX*distX + distY*distY))

			if !intersectionFound {
				shortestDistance = dist
				shortestVector = Vector{vec.A, intersection}
				intersectionFound = true
				continue
			}

			if dist < shortestDistance {
				shortestDistance = dist
				shortestVector = Vector{vec.A, intersection}
			}
		}
	}
	return shortestVector, shortestDistance, intersectionFound
}

func (a *App) vectorFromPlayer(player Player) (Vector) {
	legb := 15 * math.Cos(float64(player.ViewDirection) * (math.Pi / 180))
	lega := 15 * math.Sin(float64(player.ViewDirection) * (math.Pi / 180))

	playerX, playerY := a.demo.mapMetadata.TranslateScale(player.X, player.Y)
	playerVectorPointA := Point{playerX, -playerY}
	playerVectorPointB := Point{playerX + legb, -playerY + lega}

	return Vector{playerVectorPointA, playerVectorPointB}
}

func intersectPoint(a Point, b Point, c Point, d Point) (Point, bool) {

	//formulas from wikipedia

	denom := (c.X-d.X)*(a.Y-b.Y) - (c.Y-d.Y)*(a.X-b.X)

	if denom == 0 {
		return Point{}, false
	}

	t := ((c.X-a.X)*(a.Y-b.Y) - (c.Y-a.Y)*(a.X-b.X)) / denom


	u := ((c.X-a.X)*(c.Y-d.Y) - (c.Y-a.Y)*(c.X-d.X)) / denom

	if t > 0 && t < 1 && u > 0 {

		return Point{c.X + t*(d.X-c.X), c.Y + t*(d.Y-c.Y)}, true

	} else {
		return Point{}, false
	}
}