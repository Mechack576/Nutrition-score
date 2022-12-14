package main

type ScoreType int

// ScoreType is the type of the scored product
const (
	// Food is used when calculating nutritional score for general food items
	Food ScoreType = iota
	// Beverage is used when calculating nutritional score for beverages
	Beverage
	// Water is used when calculating nutritional score for water
	Water
	// Cheese is used for calculating the nutritional score for cheeses
	Cheese
)

// EnergyFromKcal converts energy density from kcal to EnergyKJ
func EnergyFromKcal(kcal float64) EnergyKJ {
	return EnergyKJ(kcal * 4.184)
}

// SodiumFromSalt converts salt mg/100g content to sodium content
func SodiumFromSalt(saltMg float64) SodiumMilligram {
	return SodiumFromSalt(saltMg / 2.5)
}

// NutritionalScore contains the numeric nutritional score value and type of product
type NutrionalScore struct {
	Value     int
	Positive  int
	Negative  int
	ScoreType ScoreType
}

var scoreToLetter = []string{"A", "B", "C", "D", "E"}

// EnergyKJ represents the energy density in kJ/100g
type EnergyKJ float64

// SugarGram represents amount of sugars in grams/100g
type SugarGram float64

// SaturatedFattyAcidsGram represents amount of saturated fatty acids in grams/100g
type SaturatedFattyAcids float64

// SodiumMilligram represents amount of sodium in mg/100g
type SodiumMilligram float64

// FruitsPercent represents fruits, vegetables, pulses, nuts, and rapeseed, walnut and olive oils
// as percentage of the total
type FruitsPercent float64

// FibreGram represents amount of fibre in grams/100g
type FiberGram float64

// ProteinGram represents amount of protein in grams/100g
type ProtienGram float64

// NutritionalData represents the source nutritional data used for the calculation
type NutritionalData struct {
	Energy              EnergyKJ
	Sugar               SugarGram
	SaturatedFattyAcids SaturatedFattyAcids
	Sodium              SodiumMilligram
	Fruits              FruitsPercent
	Fiber               FiberGram
	Protien             ProtienGram
	IsWater             bool
}

//Density levels per nutrient
var energyLevels = []float64{3350, 3015, 2680, 2345, 2010, 1675, 1340, 1005, 670, 335}
var sugarsLevels = []float64{45, 40, 36, 31, 27, 22.5, 18, 13.5, 9, 4.5}
var saturatedFattyAcidsLevels = []float64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
var sodiumLevels = []float64{900, 810, 720, 630, 540, 450, 360, 270, 180, 90}
var fiberLevels = []float64{4.7, 3.7, 2.8, 1.9, 0.9}
var proteinLevels = []float64{8, 6.4, 4.8, 3.2, 1.6}

var energyLevelsBeverage = []float64{270, 240, 210, 180, 150, 120, 90, 60, 30, 0}
var sugarsLevelsBeverage = []float64{13.5, 12, 10.5, 9, 7.5, 6, 4.5, 3, 1.5, 0}

// GetPoints returns the nutritional score
func (e EnergyKJ) GetPoints(st ScoreType) int {
	if st == Beverage {
		return getPointsFromRange(float64(e), energyLevelsBeverage)
	}
	return getPointsFromRange(float64(e), energyLevels)
}

func (s SugarGram) GetPoints(st ScoreType) int {
	if st == Beverage {
		return getPointsFromRange(float64(s), sugarsLevelsBeverage)
	}
	return getPointsFromRange(float64(s), sugarsLevels)
}

func (sfa SaturatedFattyAcids) GetPoints(st ScoreType) int {
	return getPointsFromRange(float64(sfa), saturatedFattyAcidsLevels)
}

func (s SodiumMilligram) GetPoints(st ScoreType) int {
	return getPointsFromRange(float64(s), sodiumLevels)
}

func (f FruitsPercent) GetPoints(st ScoreType) int {
	if st == Beverage {
		if f > 80 {
			return 10
		} else if f > 60 {
			return 4
		} else if f > 40 {
			return 2
		}
		return 0
	}
	if f > 80 {
		return 5
	} else if f > 60 {
		return 2
	} else if f > 40 {
		return 1
	}
	return 0
}

func (p ProtienGram) GetPoints(st ScoreType) int {
	return getPointsFromRange(float64(p), proteinLevels)
}

func (f FiberGram) GetPoints(st ScoreType) int {
	return getPointsFromRange(float64(f), fiberLevels)
}

// GetNutritionalScore calculates the nutritional score for nutritional data n of type st
func GetNutritionalScore(n NutritionalData, st ScoreType) NutrionalScore {

	value := 0
	positive := 0
	negative := 0
	// Water is always graded A
	if st != Water {
		fruitPoints := n.Fruits.GetPoints(st)
		fiberPoints := n.Fiber.GetPoints(st)
		//negative points are the negative things like calories (it says energy but these are what people are avoiding as these are calories)
		//sugars, saturated fats and sodium
		//positives are fruit points, fiber points and proteins
		negative = n.Energy.GetPoints(st) + n.Sugar.GetPoints(st) + n.SaturatedFattyAcids.GetPoints(st) + n.Sodium.GetPoints(st)
		positive = fruitPoints + fiberPoints + n.Protien.GetPoints(st)

		// Cheeses always use (negative - positive)
		if st == Cheese {
			value = negative - positive
		} else {
			if negative >= 11 && fruitPoints < 5 {
				value = negative - positive - fruitPoints
			} else {
				value = negative - positive
			}
		}
	}

	return NutrionalScore{
		Value:     value,
		Positive:  positive,
		Negative:  negative,
		ScoreType: st,
	}
}

// GetNutriScore returns the Nutri-Score rating
func (ns NutrionalScore) GetNutriScore() string {

	if ns.ScoreType == Food {
		return scoreToLetter[getPointsFromRange(float64(ns.Value), []float64{18, 10, 2, -1})]
	}
	if ns.ScoreType == Water {
		return scoreToLetter[0]
	}
	return scoreToLetter[getPointsFromRange(float64(ns.Value), []float64{9, 5, 1, -2})]

}

func getPointsFromRange(v float64, steps []float64) int {
	lenSteps := len(steps)
	for i, l := range steps {
		if v > l {
			return lenSteps - i
		}
	}
	return 0
}
