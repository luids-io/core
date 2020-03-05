// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package xlist

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
)

// Resource stores the type of resource used by the RBLs
type Resource int

// List of valid resources at the current time
const (
	IPv4 Resource = iota
	IPv6
	Domain
	MD5
	SHA1
	SHA256
)

// Resources is an ordered vector that constains all valid resource values.
// Warning: It's a variable for simplicity, Do not modify the value!
var Resources = []Resource{IPv4, IPv6, Domain, MD5, SHA1, SHA256}

// IsValid returns true if the resource value is a valid
func (r Resource) IsValid() bool {
	v := int(r)
	if v >= int(IPv4) && v <= int(SHA256) {
		return true
	}
	return false
}

// ValidFormat returns true if format passed is valid
func (r Resource) ValidFormat(f Format) bool {
	if f == Plain {
		return true
	}
	if (r == IPv4 || r == IPv6) && f == CIDR {
		return true
	}
	if r == Domain && f == Sub {
		return true
	}
	return false
}

// InArray returns true if the resource value exists in the array passed
// as parameter
func (r Resource) InArray(array []Resource) bool {
	for _, c := range array {
		if r == c {
			return true
		}
	}
	return false
}

func (r Resource) string() string {
	switch r {
	case IPv4:
		return "ip4"
	case IPv6:
		return "ip6"
	case Domain:
		return "domain"
	case MD5:
		return "md5"
	case SHA1:
		return "sha1"
	case SHA256:
		return "sha256"
	default:
		return ""
	}
}

// String implements stringer interface
func (r Resource) String() string {
	s := r.string()
	if s == "" {
		return fmt.Sprintf("unkown(%d)", r)
	}
	return s
}

// ToResource returns the resource type from its string representation
func ToResource(s string) (Resource, error) {
	switch strings.ToLower(s) {
	case "ip4":
		return IPv4, nil
	case "ip6":
		return IPv6, nil
	case "domain":
		return Domain, nil
	case "md5":
		return MD5, nil
	case "sha1":
		return SHA1, nil
	case "sha256":
		return SHA256, nil
	default:
		return Resource(-1), fmt.Errorf("invalid resource %s", s)
	}
}

// MarshalJSON implements interface for struct marshalling
func (r Resource) MarshalJSON() ([]byte, error) {
	s := r.string()
	if s == "" {
		return nil, fmt.Errorf("invalid value %v for resource", r)
	}
	return json.Marshal(s)
}

// UnmarshalJSON implements interface for struct unmarshalling
func (r *Resource) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "ip4":
		*r = IPv4
		return nil
	case "ip6":
		*r = IPv6
		return nil
	case "domain":
		*r = Domain
		return nil
	case "md5":
		*r = MD5
		return nil
	case "sha1":
		*r = SHA1
		return nil
	case "sha256":
		*r = SHA256
		return nil
	default:
		return fmt.Errorf("cannot unmarshal resource %s", s)
	}
}

// ValidResource returns true if the name is of the resource type
func ValidResource(name string, resource Resource) bool {
	switch resource {
	case IPv4:
		return isIPv4(name)
	case IPv6:
		return isIPv6(name)
	case Domain:
		return isDomain(name)
	case MD5:
		return isMD5(name)
	case SHA1:
		return isSHA1(name)
	case SHA256:
		return isSHA256(name)
	default:
		return false
	}
}

// Canonicalize returns true if the name is of the resource type and returns
// a string with name canonicalized to the resource type
func Canonicalize(name string, resource Resource) (string, bool) {
	switch resource {
	case IPv4:
		return canonIPv4(name)
	case IPv6:
		return canonIPv6(name)
	case Domain:
		return canonDomain(name)
	case MD5:
		return canonMD5(name)
	case SHA1:
		return canonSHA1(name)
	case SHA256:
		return canonSHA256(name)
	default:
		return name, false
	}
}

// ResourceType tries to guess the resource type of the string passed.
// Returns an error if the type can't be guessed.
func ResourceType(s string, matchOrder []Resource) (Resource, error) {
	for _, r := range matchOrder {
		if ValidResource(s, r) {
			return r, nil
		}
	}
	return Resource(-1), errors.New("resource type can't be guessed")
}

// defined for context values
type key int

const (
	keyValidated key = iota
)

// DoValidation validates the name and the resource type, and canonicalizes
// de name. If the resource
// is valid, it will return a context with a flag that indicates that
// future calls to the function should not validate the resource again,
// avoiding redundant validations. A validation can be forced if the force
// flag is set.
// If the validation is not successful, and error will be returned.
// This function must be used by the components that implement the interface
// Check and it should not be used outside of this context.
// To validate in any other use case, the function ValidResource must be used.
func DoValidation(ctx context.Context, name string, resource Resource, force bool) (string, context.Context, error) {
	if !force {
		validated, _ := ctx.Value(keyValidated).(bool)
		if validated {
			return name, ctx, nil
		}
	}
	if !resource.IsValid() {
		return name, ctx, ErrNotImplemented
	}
	canon, ok := Canonicalize(name, resource)
	if !ok {
		return name, ctx, ErrBadRequest
	}

	return canon, context.WithValue(ctx, keyValidated, true), nil
}

// ClearResourceDups returns an array with duplicate and invalid resource items
// removed
func ClearResourceDups(resources []Resource) []Resource {
	ret := make([]Resource, 0, len(Resources))
	copied := make([]bool, len(Resources), len(Resources))
	for _, r := range resources {
		if !r.IsValid() {
			continue
		}
		if !copied[int(r)] {
			ret = append(ret, r)
			copied[int(r)] = true
		}
	}
	return ret
}

//validation and canonicalization functions for resources

func isIPv4(s string) bool {
	ip := net.ParseIP(s)
	if ip == nil {
		return false
	}
	if ip.To4() == nil {
		return false
	}
	return true
}

func canonIPv4(s string) (string, bool) {
	ip := net.ParseIP(s)
	if ip == nil {
		return s, false
	}
	if ip.To4() == nil {
		return s, false
	}
	return ip.String(), true
}

func isCIDR4(s string) bool {
	ip, _, err := net.ParseCIDR(s)
	if err != nil {
		return false
	}
	if ip.To4() == nil {
		return false
	}
	return true
}

func isCIDR6(s string) bool {
	ip, _, err := net.ParseCIDR(s)
	if err != nil {
		return false
	}
	if ip.To4() != nil {
		return false
	}
	return true
}

func isIPv6(s string) bool {
	ip := net.ParseIP(s)
	if ip == nil {
		return false
	}
	if ip.To4() != nil {
		return false
	}
	return true
}

func canonIPv6(s string) (string, bool) {
	ip := net.ParseIP(s)
	if ip == nil {
		return s, false
	}
	if ip.To4() != nil {
		return s, false
	}
	return ip.String(), true
}

// note: we precompute for performance reasons
var validDomainRegexp, _ = regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9_][a-zA-Z0-9\-_]*[a-zA-Z0-9_])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)

func isDomain(s string) bool {
	return validDomainRegexp.MatchString(s)
}

func canonDomain(s string) (string, bool) {
	ok := validDomainRegexp.MatchString(s)
	if !ok {
		return s, false
	}
	return strings.ToLower(s), true
}

// note: we precompute for performance reasons
var validMD5Regexp, _ = regexp.Compile(`^[a-fA-F0-9]{32}$`)

func isMD5(s string) bool {
	return validMD5Regexp.MatchString(s)
}

func canonMD5(s string) (string, bool) {
	ok := validMD5Regexp.MatchString(s)
	if !ok {
		return s, false
	}
	return strings.ToLower(s), true
}

// note: we precompute for performance reasons
var validSHA1Regexp, _ = regexp.Compile(`^[a-fA-F0-9]{40}$`)

func isSHA1(s string) bool {
	return validSHA1Regexp.MatchString(s)
}

func canonSHA1(s string) (string, bool) {
	ok := validSHA1Regexp.MatchString(s)
	if !ok {
		return s, false
	}
	return strings.ToLower(s), true
}

// note: we precompute for performance reasons
var validSHA256Regexp, _ = regexp.Compile(`^[a-fA-F0-9]{64}$`)

func isSHA256(s string) bool {
	return validSHA256Regexp.MatchString(s)
}

func canonSHA256(s string) (string, bool) {
	ok := validSHA256Regexp.MatchString(s)
	if !ok {
		return s, false
	}
	return strings.ToLower(s), true
}
