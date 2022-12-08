package cache

import "testing"

func TestBadgerCache_Ping(t *testing.T) {
	res, err := testBadgerCache.Ping()
	if err != nil {
		t.Error(err)
	}

	if res != "PONG" {
		t.Error("did not receive PONG when expected")
	}
}

func TestBadgerCache_Has(t *testing.T) {
	err := testBadgerCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}

	// get a value that does not exist
	inCache, err := testBadgerCache.Has("foo")
	if err == nil {
		t.Error("no error returned when error should be key not found")
	}
	if inCache {
		t.Error("foo found in cache and shouldn't be there")
	}

	// get a value that does exist
	_ = testBadgerCache.Set("foo", "bar")
	inCache, err = testBadgerCache.Has("foo")
	if err != nil {
		t.Error(err)
	}
	if !inCache {
		t.Error("foo not found in cache and should be there")
	}

	// forget test value
	err = testBadgerCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}
}

func TestBadgerCache_Get(t *testing.T) {
	err := testBadgerCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	x, err := testBadgerCache.Get("foo")
	if err != nil {
		t.Error(err)
	}
	if x != "bar" {
		t.Error("did not get correct value from cache")
	}

	// forget test value
	err = testBadgerCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}
}

func TestBadgerCache_Forget(t *testing.T) {
	err := testBadgerCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Forget("alpha")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("alpha")
	if err == nil {
		t.Error("no error returned when error should be key not found")
	}
	if inCache {
		t.Error("found value in cache that shouldn't be there")
	}
}

func TestBadgerCache_Empty(t *testing.T) {
	err := testBadgerCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Empty()
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("foo")
	if err == nil {
		t.Error("no error returned when error should be key not found")
	}

	if inCache {
		t.Error("foo found in cache and it should not be there")
	}
}

func TestBadgerCache_EmptyByMatch(t *testing.T) {
	err := testBadgerCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Set("foo", "foo")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Set("bar", "foo")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.EmptyByMatch("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("foo")
	if err == nil {
		t.Error("no error returned when error should be key not found")
	}
	if inCache {
		t.Error("foo found in cache and it should not be there")
	}

	inCache, err = testBadgerCache.Has("bar")
	if err != nil {
		t.Error(err)
	}
	if !inCache {
		t.Error("bar not found and it should be there")
	}
}
