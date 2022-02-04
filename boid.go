package main

import (
	"math"
	"math/rand"
	"time"
)

type Boid struct {
	position Vector2d
	velocity Vector2d
	id int
} 
func (b *Boid) calcAcceleration() Vector2d {
	upper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)
	sumVelocity := Vector2d{0,0}
	sumPosition := Vector2d{0,0}
	separation := Vector2d{0,0}
	count := 0

	rwLock.RLock()
	for i := math.Max(lower.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j <= math.Min(upper.y, screenHeight); j++ {
			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.id {
				if distance := boids[otherBoidId].position.Distance(b.position); distance < viewRadius {
					count++
					sumVelocity = sumVelocity.Add(boids[otherBoidId].velocity)
					sumPosition = sumPosition.Add(boids[otherBoidId].position)
					separation = separation.Add(b.position.Subtract(boids[otherBoidId].position).DivisionV(distance))
				
				}

			}
		}
	}
	rwLock.RUnlock()

	accel := Vector2d{b.borderBounce(b.position.x, screenWidth), b.borderBounce(b.position.y, screenHeight)}
	if count == 0 || rand.Intn(100) > 20 {
		return accel	
	}

	avgVelocity := sumVelocity.DivisionV(float64(count))
	avgPosition := sumPosition.DivisionV(float64(count))

	accelAlignment :=  avgVelocity.Subtract(b.velocity).MultiplyV(adjRate)
	accelCohesion  :=  avgPosition.Subtract(b.position).MultiplyV(adjRate)
	accelSeparation := separation.MultiplyV(adjRate)

	accel = accel.Add(accelAlignment).Add(accelCohesion).Add(accelSeparation)

	if rand.Intn(100) > 90 {
		accel = accel.MultiplyV(adjRate)
	}


	return accel
}

func (b *Boid) borderBounce(pos, maxBorderPos float64) float64 {
	if pos < viewRadius {
		return 1/ pos
	} else if pos > maxBorderPos - viewRadius {
		return 1/ (pos - maxBorderPos)
	}

	return 0
}

func (b * Boid) moveOne() {
	acc := b.calcAcceleration()
	rwLock.Lock()
	b.velocity = b.velocity.Add(acc).limit(-1, 1)
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.Add(b.velocity)
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	rwLock.Unlock()
}
func (b *Boid) start() {
	for {
		b.moveOne()
		time.Sleep(1 * time.Millisecond)
	}
}
func createBoid(bid int) {
	b := Boid {
		position: Vector2d{x: rand.Float64()* screenWidth, y: rand.Float64()* screenHeight},
		velocity: Vector2d{x: (rand.Float64() * 2) - 1.0, y: (rand.Float64() * 2) - 1.0},
		id: bid,
	}

	boids[bid] = &b

	boidMap[int(b.position.x)][int(b.position.y)] = b.id

	go b.start()

}