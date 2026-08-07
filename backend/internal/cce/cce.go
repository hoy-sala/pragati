package cce

import "math"

type ClassRange string

const (
	Class6to8 ClassRange = "6-8"
	Class9to10 ClassRange = "9-10"
)

type SubjectType string

const (
	FirstLanguage SubjectType = "first_language"
	OtherSubject  SubjectType = "other"
)

// ConversionResult holds the calculated marks for a subject.
type ConversionResult struct {
	FAConverted   float64 `json:"fa_converted"`   // Internal FA marks (after conversion)
	SAConverted   float64 `json:"sa_converted"`   // SA marks (after conversion)
	FinalMarks    float64 `json:"final_marks"`    // FA + SA
	FinalGrade    string  `json:"final_grade"`
	MaxMarks      float64 `json:"max_marks"`
	FARawTotal    float64 `json:"fa_raw_total"`
	SARawTotal    float64 `json:"sa_raw_total"`
	FAOutOf       float64 `json:"fa_out_of"`
	SAOutOf       float64 `json:"sa_out_of"`
}

// Grade boundaries for final marks (out of 100)
var gradeBoundaries = []struct {
	grade string
	min   float64
}{
	{"A+", 90},
	{"A", 70},
	{"B+", 50},
	{"B", 40},
	{"C+", 30},
	{"C", 0},
}

func GradeFromMarks(pct float64) string {
	for _, g := range gradeBoundaries {
		if pct >= g.min {
			return g.grade
		}
	}
	return "C"
}

// ConvertClasses68 applies CCE conversion for classes 6-8.
// Each FA is out of 30 marks (20 written + 10 activities), converted to 10.
// Each SA is out of 50 marks, converted to 30.
// Final = sum of all FA (out of 40) + sum of all SA (out of 60) = 100.
func ConvertClasses68(faMarks, saMarks []float64) ConversionResult {
	var faRawTotal float64
	for _, m := range faMarks {
		faRawTotal += m
	}
	// Each FA max 30, total FA max = 4 * 30 = 120, converted to 40
	faConverted := (faRawTotal / 120.0) * 40.0

	var saRawTotal float64
	for _, m := range saMarks {
		saRawTotal += m
	}
	// Each SA max 50, total SA max = 2 * 50 = 100, converted to 60
	saConverted := (saRawTotal / 100.0) * 60.0

	final := faConverted + saConverted
	final = math.Round(final*100) / 100

	return ConversionResult{
		FAConverted: round2(faConverted),
		SAConverted: round2(saConverted),
		FinalMarks:  final,
		FinalGrade:  GradeFromMarks(final),
		MaxMarks:    100,
		FARawTotal:  round2(faRawTotal),
		SARawTotal:  round2(saRawTotal),
		FAOutOf:     120,
		SAOutOf:     100,
	}
}

// ConvertClasses910 applies CCE conversion for classes 9-10.
// Each FA is out of 50 marks. Total FA (all 4) = 200 marks raw.
// First Language: FA converted to 25, SA (out of 125) converted to 75.
// Other Subjects: FA converted to 20, SA (out of 100) converted to 80.
// Final = FA_converted + SA_converted = 100.
func ConvertClasses910(faMarks, saMarks []float64, subjType SubjectType) ConversionResult {
	var faRawTotal float64
	for _, m := range faMarks {
		faRawTotal += m
	}

	var saRawTotal float64
	for _, m := range saMarks {
		saRawTotal += m
	}

	var faConverted, saConverted float64

	if subjType == FirstLanguage {
		// FA: 200 raw → 25, SA: 125 raw → 75
		faConverted = (faRawTotal / 200.0) * 25.0
		saConverted = (saRawTotal / 125.0) * 75.0
	} else {
		// FA: 200 raw → 20, SA: 100 raw → 80
		faConverted = (faRawTotal / 200.0) * 20.0
		saConverted = (saRawTotal / 100.0) * 80.0
	}

	final := faConverted + saConverted
	final = math.Round(final*100) / 100

	saOutOf := 100.0
	if subjType == FirstLanguage {
		saOutOf = 125.0
	}

	return ConversionResult{
		FAConverted: round2(faConverted),
		SAConverted: round2(saConverted),
		FinalMarks:  final,
		FinalGrade:  GradeFromMarks(final),
		MaxMarks:    100,
		FARawTotal:  round2(faRawTotal),
		SARawTotal:  round2(saRawTotal),
		FAOutOf:     200,
		SAOutOf:     saOutOf,
	}
}

// IsFirstLanguage checks if a subject code is a first language.
func IsFirstLanguage(subjectCode string) bool {
	return subjectCode == "KAN" // Kannada is first language in this school
}

// DetermineClassRange returns the class range from a class name.
func DetermineClassRange(className string) ClassRange {
	switch className {
	case "Class 6", "Class 7", "Class 8", "6", "7", "8":
		return Class6to8
	case "Class 9", "Class 10", "9", "10":
		return Class9to10
	}
	return Class6to8
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
