package hpc

import (
	"testing"

	"github.com/pragati/backend/internal/models"
)

var testScheme = []models.GradeDefinition{
	{Grade: "A+", MinPct: 90, MaxPct: 100, Descriptor: "Outstanding"},
	{Grade: "A", MinPct: 75, MaxPct: 89, Descriptor: "Excellent"},
	{Grade: "B+", MinPct: 60, MaxPct: 74, Descriptor: "Very Good"},
	{Grade: "B", MinPct: 45, MaxPct: 59, Descriptor: "Good"},
	{Grade: "C", MinPct: 33, MaxPct: 44, Descriptor: "Satisfactory"},
	{Grade: "D", MinPct: 20, MaxPct: 32, Descriptor: "Needs Improvement"},
	{Grade: "E", MinPct: 0, MaxPct: 19, Descriptor: "Requires Remedial"},
}

func TestComputeGrade(t *testing.T) {
	tests := []struct {
		scored, max float64
		expected    string
	}{
		{95, 100, "A+"},
		{80, 100, "A"},
		{65, 100, "B+"},
		{50, 100, "B"},
		{38, 100, "C"},
		{25, 100, "D"},
		{10, 100, "E"},
		{0, 100, "E"},
		{100, 0, "N/A"},
	}
	for _, tc := range tests {
		got := ComputeGrade(tc.scored, tc.max, testScheme)
		if got != tc.expected {
			t.Errorf("ComputeGrade(%.0f, %.0f) = %s; want %s", tc.scored, tc.max, got, tc.expected)
		}
	}
}

func TestComputeGradePoint(t *testing.T) {
	tests := []struct {
		grade    string
		expected float64
	}{
		{"A+", 10},
		{"A", 9},
		{"B+", 8},
		{"B", 7},
		{"C", 6},
		{"D", 5},
		{"E", 4},
		{"F", 0},
	}
	for _, tc := range tests {
		got := ComputeGradePoint(tc.grade, testScheme)
		if got != tc.expected {
			t.Errorf("ComputeGradePoint(%s) = %.0f; want %.0f", tc.grade, got, tc.expected)
		}
	}
}

func TestComputeHPCScholastic(t *testing.T) {
	marks := map[string]struct{ Scored, Max float64 }{
		"math":  {85, 100},
		"eng":   {30, 100},
		"sci":   {92, 100},
	}
	entries := ComputeHPCScholastic(marks, testScheme)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	gradeMap := map[string]string{}
	for _, e := range entries {
		gradeMap[e.SubjectID] = e.Grade
	}
	if gradeMap["math"] != "A" {
		t.Errorf("math grade = %s; want A", gradeMap["math"])
	}
	if gradeMap["eng"] != "D" {
		t.Errorf("eng grade = %s; want D", gradeMap["eng"])
	}
	if gradeMap["sci"] != "A+" {
		t.Errorf("sci grade = %s; want A+", gradeMap["sci"])
	}
}

func TestGradeFromPercentile(t *testing.T) {
	tests := []struct {
		pct         float64
		wantGrade   string
		wantGp      float64
	}{
		{95, "A+", 10},
		{80, "A", 9},
		{65, "B+", 8},
		{50, "B", 7},
		{38, "C", 6},
		{25, "D", 5},
		{10, "E", 4},
	}
	for _, tc := range tests {
		g, gp := GradeFromPercentile(tc.pct)
		if g != tc.wantGrade || gp != tc.wantGp {
			t.Errorf("GradeFromPercentile(%.0f) = (%s, %.0f); want (%s, %.0f)", tc.pct, g, gp, tc.wantGrade, tc.wantGp)
		}
	}
}

func TestAggregateLOAssessments(t *testing.T) {
	assessments := []models.LearningOutcomeAssessment{
		{LearningOutcomeID: "lo1", ProficiencyLevel: 2},
		{LearningOutcomeID: "lo1", ProficiencyLevel: 4},
		{LearningOutcomeID: "lo2", ProficiencyLevel: 3},
	}
	result := AggregateLOAssessments(assessments)
	if len(result) != 2 {
		t.Fatalf("expected 2 aggregated entries, got %d", len(result))
	}
	for _, r := range result {
		if r.LOID == "lo1" && r.ProficiencyLevel != 4 {
			t.Errorf("lo1 proficiency level = %d; want 4 (highest)", r.ProficiencyLevel)
		}
		if r.LOID == "lo2" && r.ProficiencyLevel != 3 {
			t.Errorf("lo2 proficiency level = %d; want 3", r.ProficiencyLevel)
		}
	}
}
