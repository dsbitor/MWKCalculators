package main

import (
	"math"
	"testing"
)

// A single self-consistent feed rate scenario: speed 1000 rpm, 4
// cutting edges, chip load 0.002 in/edge gives feed 8 in/min and
// 0.008 in/rev. Used to check that solving from any three of the
// five quantities recovers the other two.
const (
	testSpeed      = 1000.0
	testEdges      = 4.0
	testLoad       = 0.002
	testFeedPerMin = 8.0
	testFeedPerRev = 0.008
)

func TestSolveFeedRate_AllFourSolvableCombinations(t *testing.T) {
	cases := []struct {
		name                                       string
		speed, edges, load, feedPerMin, feedPerRev float64
	}{
		{"speed+edges+load", testSpeed, testEdges, testLoad, 0, 0},
		{"speed+edges+feedPerMin", testSpeed, testEdges, 0, testFeedPerMin, 0},
		{"speed+edges+feedPerRev", testSpeed, testEdges, 0, 0, testFeedPerRev},
		{"speed+load+feedPerMin", testSpeed, 0, testLoad, testFeedPerMin, 0},
		{"speed+load+feedPerRev", testSpeed, 0, testLoad, 0, testFeedPerRev},
		{"edges+load+feedPerMin", 0, testEdges, testLoad, testFeedPerMin, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := solveFeedRate(c.speed, c.edges, c.load, c.feedPerMin, c.feedPerRev)
			if err != nil {
				t.Fatalf("solveFeedRate() error = %v", err)
			}
			checkClose(t, "Speed", got.Speed, testSpeed)
			checkClose(t, "Edges", got.Edges, testEdges)
			checkClose(t, "Load", got.Load, testLoad)
			checkClose(t, "FeedPerMin", got.FeedPerMin, testFeedPerMin)
			checkClose(t, "FeedPerRev", got.FeedPerRev, testFeedPerRev)
		})
	}
}

func TestSolveFeedRate_EdgesAndLoadAloneIsInsufficient(t *testing.T) {
	// Knowing only edges and load (with no speed and no feed figure)
	// doesn't determine anything: an error, not a nonsensical
	// result, matching the original program's own "whoops" case.
	_, err := solveFeedRate(0, testEdges, testLoad, 0, 0)
	if err == nil {
		t.Fatal("solveFeedRate() error = nil, want an error for edges+load alone")
	}
}

func TestSolveFeedRate_EdgesLoadAndFeedPerRevIsNotSolvable(t *testing.T) {
	// edges+load+feedPerRev is deliberately not one of the four
	// solvable combinations: feedPerRev is already fully determined
	// by edges and load alone (feedPerRev = load*edges), so it adds
	// no information toward finding speed. The original program
	// never even prompts for feedPerRev in this situation.
	_, err := solveFeedRate(0, testEdges, testLoad, 0, testFeedPerRev)
	if err == nil {
		t.Fatal("solveFeedRate() error = nil, want an error for edges+load+feedPerRev")
	}
}

func TestSolveFeedRate_TooFewKnownValuesIsInsufficient(t *testing.T) {
	_, err := solveFeedRate(testSpeed, 0, 0, 0, 0)
	if err == nil {
		t.Fatal("solveFeedRate() error = nil, want an error for a single known value")
	}
}

func checkClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
