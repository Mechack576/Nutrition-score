package main

import (
	"fmt"
)

func main() {
	ns := GetNutritionalScore(NutritionalData{
		Energy:              EnergyFromKcal(100),
		Sugar:               SugarGram(10),
		SaturatedFattyAcids: SaturatedFattyAcids(2),
		Sodium:              SodiumMilligram(500),
		Fruits:              FruitsPercent(60),
		Fiber:               FiberGram(4),
		Protien:             ProtienGram(2),
	}, Food)

	// Output:
	// Nutritional score: 3
	// NutriScore: C

	fmt.Printf("Nutrional Score:%d\n", ns.Value)
	fmt.Printf("NutriScore: %s\n", ns.GetNutriScore())
}
