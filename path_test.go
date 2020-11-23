package httpjson

import "testing"

func TestPath(t *testing.T) {
	assertEqual(t, uint64(0), RetrievePrefixedPathInteger("", ""))
	assertEqual(t, uint64(0), RetrievePrefixedPathInteger("/", ""))
	assertEqual(t, uint64(0), RetrievePrefixedPathInteger("/", "hello"))
	assertEqual(t, uint64(0), RetrievePrefixedPathInteger("/hello", "hello"))
	assertEqual(t, uint64(123), RetrievePrefixedPathInteger("/hello/123", "hello"))
	assertEqual(t, uint64(0), RetrievePrefixedPathInteger("/hello/123/goodbye", "goodbye"))
	assertEqual(t, uint64(456), RetrievePrefixedPathInteger("/hello/123/goodbye/456", "goodbye"))
}
