package acceptance

import "context"

var (
	failures int
	success  int
	features []*FeatureImpl
)

// FeatureFunc is the func type run to test features
type FeatureFunc func(ctx context.Context)

// FeatureImpl is a feature implementation.
type FeatureImpl struct {
	label string
	name  string

	fn FeatureFunc

	hasRun bool
	failed bool

	teardown *tearDownImpl

	before string
	inside string

	requiredFlags []string
}

type tearDownImpl struct {
	name string

	fn FeatureFunc
}

// Feature represents a testable feature
func Feature(label, name string, fn FeatureFunc) *FeatureImpl {
	fi := &FeatureImpl{
		label: label,
		name:  name,
		fn:    fn,
	}
	features = append(features, fi)

	return fi
}

// TearDown represents the optional teardown step for a feature
func (f *FeatureImpl) TearDown(name string, fn FeatureFunc) interface{} {
	td := &tearDownImpl{
		name: name,
		fn:   fn,
	}

	f.teardown = td

	return td
}

// RunsBefore expresses a dependency between two features. A feature that runs
// before another feature will execute its tests, and tear down before the
// other is run.
func (f *FeatureImpl) RunsBefore(label string) interface{} {
	f.before = label
	return nil
}

// RunsInside expresses a dependency between two features. A feature that runs
// inside another feature will execute its tests, and tear down after the
// other is run, but before it is torn down.
func (f *FeatureImpl) RunsInside(label string) interface{} {
	f.inside = label
	return nil
}

// RequiredFlags marks a set of flags required for a feature. We will use this
// to see that all flags are provided when executing a test.
func (f *FeatureImpl) RequiredFlags(flags ...string) interface{} {
	for _, flag := range flags {
		if !f.NeedsFlag(flag) {
			f.requiredFlags = append(f.requiredFlags, flag)
		}
	}

	return nil
}

// NeedsFlag checks if a Feature Implementation needs a specific flag or not.
func (f *FeatureImpl) NeedsFlag(flag string) bool {
	for _, rflag := range f.requiredFlags {
		if rflag == flag {
			return true
		}
	}

	return false
}

// Default represents the default test case of a feature.
func Default(fn func()) {
	block("Default case", fn)
}

// ErrorCase represents an error case for a feature. These are optionally tested
// based on command line flags.
//
// Failure of an error case does not prevent further error cases or RunsInside
// features from running.
func ErrorCase(name string, fn func()) {
	if !shouldRunErrorCases {
		return
	}

	defer func() { recover() }()
	block("Error case: "+name, fn)
}

func block(name string, fn func()) {
	res := fail
	enter(name)
	defer func() {
		result(name, res)
		exit()
	}()

	fn()
	res = pass
}
