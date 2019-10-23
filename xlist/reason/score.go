// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package reason

import (
	"fmt"
	"strconv"
	"strings"
)

// scoreFromString loads from string the value
func scoreFromString(s string) (int, error) {
	slower := strings.ToLower(s)
	if !strings.HasPrefix(slower, "[score]") || !strings.HasSuffix(slower, "[/score]") {
		return 0, fmt.Errorf("invalid score '%s'", s)
	}
	s = s[7 : len(slower)-8] //len('[score]')==7, len('[/score]')==8
	if s == "" {
		return 0, nil
	}
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid score '%s'", s)
	}
	return value, nil
}

// scoreToString returns the score string
func scoreToString(score int) string {
	return fmt.Sprintf("[score]%v[/score]", score)
}

// WithScore inserts a score inside a reason string. If there is a score
// inside, WithScore will replace it.
func WithScore(score int, s string) string {
	s = cleanScore(s)
	if score == 0 {
		return s
	}
	return fmt.Sprintf("%s%s", scoreToString(score), s)
}

// ExtractScore extracts a score from a reason string. It returns the score,
// an string reason without the score and error
func ExtractScore(s string) (int, string, error) {
	value := 0
	scores, reason := extractScoreStr(s)
	if len(scores) > 0 {
		for _, scoreStr := range scores {
			v, err := scoreFromString(scoreStr)
			if err != nil {
				return value, reason, fmt.Errorf("invalid score '%s': %v", scoreStr, err)
			}
			value = value + v
		}
	}
	return value, reason, nil
}

func cleanScore(s string) (rest string) {
	rest = ""
	for len(s) > 0 {
		slow := strings.ToLower(s)
		first := strings.Index(slow, "[score]")
		last := strings.Index(slow, "[/score]")
		if first < 0 || last < 0 {
			rest = rest + s
			return
		}
		if last < first {
			//len("[/score]") == 8
			rest = rest + s[:last+8]
			s = s[last+8:]
			continue
		}
		rest = rest + s[:first]
		s = s[last+8:]
	}
	return
}

func extractScoreStr(s string) (scores []string, rest string) {
	rest = ""
	scores = make([]string, 0)
	for len(s) > 0 {
		slow := strings.ToLower(s)
		first := strings.Index(slow, "[score]")
		last := strings.Index(slow, "[/score]")
		if first < 0 || last < 0 {
			rest = rest + s
			return
		}
		if last < first {
			//len("[/score]") == 8
			rest = rest + s[:last+8]
			s = s[last+8:]
			continue
		}
		scores = append(scores, s[first:last+8])
		rest = rest + s[:first]
		s = s[last+8:]
	}
	return
}
