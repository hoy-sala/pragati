package hpc

import (
	"sort"

	"github.com/pragati/backend/internal/models"
)

func ComputeGrade(marksScored, maxMarks float64, scheme []models.GradeDefinition) string {
	if maxMarks <= 0 {
		return "N/A"
	}
	pct := (marksScored / maxMarks) * 100

	sorted := make([]models.GradeDefinition, len(scheme))
	copy(sorted, scheme)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].MinPct > sorted[j].MinPct
	})

	for _, g := range sorted {
		if pct >= g.MinPct && pct <= g.MaxPct {
			return g.Grade
		}
	}
	return "E"
}

func ComputeGradePoint(grade string, scheme []models.GradeDefinition) float64 {
	points := map[string]float64{
		"A+": 10, "A": 9, "B+": 8, "B": 7,
		"C": 6, "D": 5, "E": 4,
	}
	if v, ok := points[grade]; ok {
		return v
	}
	return 0
}

func AggregateLOAssessments(assessments []models.LearningOutcomeAssessment) []models.LOAttainmentEntry {
	seen := make(map[string]*models.LOAttainmentEntry)
	for _, a := range assessments {
		if existing, ok := seen[a.LearningOutcomeID]; ok {
			if a.ProficiencyLevel > existing.ProficiencyLevel {
				existing.ProficiencyLevel = a.ProficiencyLevel
			}
		} else {
			seen[a.LearningOutcomeID] = &models.LOAttainmentEntry{
				LOID:             a.LearningOutcomeID,
				Code:             a.LearningOutcomeID,
				ProficiencyLevel: a.ProficiencyLevel,
			}
		}
	}
	result := make([]models.LOAttainmentEntry, 0, len(seen))
	for _, v := range seen {
		result = append(result, *v)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Code < result[j].Code
	})
	return result
}

func ComputeHPCScholastic(marksBySubject map[string]struct{ Scored, Max float64 }, scheme []models.GradeDefinition) []models.ScholasticEntry {
	entries := []models.ScholasticEntry{}
	for subjectID, m := range marksBySubject {
		grade := ComputeGrade(m.Scored, m.Max, scheme)
		entries = append(entries, models.ScholasticEntry{
			SubjectID:   subjectID,
			MaxMarks:    m.Max,
			MarksScored: m.Scored,
			Grade:       grade,
			GradePoint:  ComputeGradePoint(grade, scheme),
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].SubjectID < entries[j].SubjectID
	})
	return entries
}

func GradeFromPercentile(pct float64) (string, float64) {
	switch {
	case pct >= 90:
		return "A+", 10
	case pct >= 75:
		return "A", 9
	case pct >= 60:
		return "B+", 8
	case pct >= 45:
		return "B", 7
	case pct >= 33:
		return "C", 6
	case pct >= 20:
		return "D", 5
	default:
		return "E", 4
	}
}
