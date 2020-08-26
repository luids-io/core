// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package reason

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Policy stores fields and values (both are strings). Warning, it's unsafe.
type Policy struct {
	m map[string]string
}

// NewPolicy returns a new policy.
func NewPolicy() Policy {
	return Policy{m: make(map[string]string, 0)}
}

// FromString loads from a encoded string the values of a policy.
// String must be in the format:
//  [policy]field1=value1,field2=value2[/policy]
func (p *Policy) FromString(s string) error {
	slow := strings.ToLower(s)
	if !strings.HasPrefix(slow, "[policy]") || !strings.HasSuffix(slow, "[/policy]") {
		return errors.New("invalid string")
	}
	s = s[8 : len(slow)-9] //len('[policy]')==8, len('[/policy]')==9

	//load fields an values
	m := make(map[string]string, 0)
	args := strings.Split(s, ",")
	for _, arg := range args {
		tmps := strings.SplitN(arg, "=", 2)
		field := strings.TrimSpace(tmps[0])
		value := ""
		if len(tmps) > 1 {
			value = strings.TrimSpace(tmps[1])
		}
		if !fieldRegExp.MatchString(field) {
			return fmt.Errorf("invalid field '%s'", field)
		}
		if !valueRegExp.MatchString(value) {
			return fmt.Errorf("invalid value '%s'", value)
		}
		m[field] = value
	}
	p.m = m
	return nil
}

// Empty returns true if the policy is empty.
func (p Policy) Empty() bool {
	return len(p.m) == 0
}

// Get returns the value of the field.
func (p Policy) Get(field string) (string, bool) {
	value, ok := p.m[field]
	return value, ok
}

// Set a new field policy or modify existing value.
func (p Policy) Set(field, value string) error {
	if !fieldRegExp.MatchString(field) {
		return errors.New("invalid field")
	}
	if !valueRegExp.MatchString(value) {
		return errors.New("invalid value")
	}
	p.m[field] = value
	return nil
}

// Merge policies.
func (p Policy) Merge(policies ...Policy) {
	for _, policy := range policies {
		for k, v := range policy.m {
			p.m[k] = v
		}
	}
}

// Fields returns the fields in the policy.
func (p Policy) Fields() []string {
	fields := make([]string, 0, len(p.m))
	for k := range p.m {
		fields = append(fields, k)
	}
	return fields
}

// String returns the policy encoded as string.
func (p Policy) String() string {
	if len(p.m) == 0 {
		return ""
	}
	s := "[policy]"
	sep := false
	for k, v := range p.m {
		if sep {
			s = s + ","
		} else {
			sep = true
		}
		s = s + fmt.Sprintf("%s=%s", k, v)
	}
	s = s + "[/policy]"
	return s
}

// WithPolicy inserts a policy inside a reason string. If there is a policy
// inside, WithPolicy will replace it.
func WithPolicy(policy Policy, s string) string {
	if policy.Empty() {
		return cleanPolicy(s)
	}
	return fmt.Sprintf("%s%s", policy.String(), cleanPolicy(s))
}

// ExtractPolicy extracts a policy from a reason string. It returns the policy,
// an string reason without the policy and error.
func ExtractPolicy(s string) (Policy, string, error) {
	policies, reason := extractPolicyStr(s)
	p := NewPolicy()
	if len(policies) > 0 {
		for _, policyStr := range policies {
			np := NewPolicy()
			err := np.FromString(policyStr)
			if err != nil {
				return p, reason, fmt.Errorf("invalid policy '%s': %v", policyStr, err)
			}
			p.Merge(np)
		}
	}
	return p, reason, nil
}

func cleanPolicy(s string) (rest string) {
	rest = ""
	for len(s) > 0 {
		slow := strings.ToLower(s)
		first := strings.Index(slow, "[policy]")
		last := strings.Index(slow, "[/policy]")
		if first < 0 || last < 0 {
			rest = rest + s
			return
		}
		if last < first {
			//len("[/policy]") == 9
			rest = rest + s[:last+9]
			s = s[last+9:]
			continue
		}
		rest = rest + s[:first]
		s = s[last+9:]
	}
	return
}

func extractPolicyStr(s string) (policies []string, rest string) {
	rest = ""
	policies = make([]string, 0)
	for len(s) > 0 {
		slow := strings.ToLower(s)
		first := strings.Index(slow, "[policy]")
		last := strings.Index(slow, "[/policy]")
		if first < 0 || last < 0 {
			rest = rest + s
			return
		}
		if last < first {
			//len("[/policy]") == 9
			rest = rest + s[:last+9]
			s = s[last+9:]
			continue
		}
		policies = append(policies, s[first:last+9])
		rest = rest + s[:first]
		s = s[last+9:]
	}
	return
}

var fieldRegExp, _ = regexp.Compile(`^[A-Za-z][A-Za-z0-9_\.]*$`)
var valueRegExp, _ = regexp.Compile(`^[A-Za-z0-9_\.,:@/]*$`)
