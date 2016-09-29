package inspector

import (
	"errors"
	"reflect"
	"testing"
)

type fakeCheck struct {
	err error
}

func (c fakeCheck) Check() error {
	return c.err
}
func (c fakeCheck) Name() string {
	return "mockCheck"
}

type fakeRule struct {
	RuleMeta
	name string
}

func (r fakeRule) Name() string { return r.name }

type fakeRuleCheckMapper struct {
	check Check
	err   error
}

func (m fakeRuleCheckMapper) GetCheckForRule(Rule) (Check, error) {
	return m.check, m.err
}

func TestEngine(t *testing.T) {
	dummyError := errors.New("dummy error...")
	tests := []struct {
		mapper          fakeRuleCheckMapper
		rule            fakeRule
		ruleWhen        []string
		facts           []string
		expectedResults []RuleResult
		expectErr       bool
	}{
		// Single rule that passes
		{
			mapper: fakeRuleCheckMapper{
				check: fakeCheck{err: nil},
			},
			rule: fakeRule{
				name: "SuccessRule",
			},
			facts: []string{},
			expectedResults: []RuleResult{
				{
					Name:    "SuccessRule",
					Success: true,
				},
			},
		},
		// Single rule that fails
		{
			mapper: fakeRuleCheckMapper{
				check: fakeCheck{err: dummyError},
			},
			rule: fakeRule{
				name: "FailRule",
			},
			facts: []string{},
			expectedResults: []RuleResult{
				{
					Name:    "FailRule",
					Success: false,
					Error:   dummyError,
				},
			},
		},
		// Single rule that should run due to facts
		{
			mapper: fakeRuleCheckMapper{
				check: fakeCheck{err: dummyError},
			},
			rule: fakeRule{
				name: "FailRule",
			},
			ruleWhen: []string{"ubuntu", "worker"},
			facts:    []string{"ubuntu", "worker", "otherFact"},
			expectedResults: []RuleResult{
				{
					Name:    "FailRule",
					Success: false,
					Error:   dummyError,
				},
			},
		},
		// Single rule that should not run due to facts
		{
			mapper: fakeRuleCheckMapper{
				check: fakeCheck{err: dummyError},
			},
			rule: fakeRule{
				name: "FailRule",
			},
			ruleWhen:        []string{"ubuntu"},
			facts:           []string{"otherFact"},
			expectedResults: []RuleResult{},
		},
		// Single rule that should run regardless of facts
		{
			mapper: fakeRuleCheckMapper{
				check: fakeCheck{err: dummyError},
			},
			rule: fakeRule{
				name: "FailRule",
			},
			ruleWhen: []string{},
			facts:    []string{"ubuntu"},
			expectedResults: []RuleResult{
				{
					Name:    "FailRule",
					Success: false,
					Error:   dummyError,
				},
			},
		},
		// Mapper returns an error, engine should return error
		{
			mapper: fakeRuleCheckMapper{
				err: dummyError,
			},
			rule:      fakeRule{},
			expectErr: true,
		},
	}

	for _, test := range tests {
		// Set the rule when conditions here because of embedding :(
		if test.ruleWhen != nil {
			test.rule.When = test.ruleWhen
		}
		e := Engine{
			RuleCheckMapper: test.mapper,
		}
		result, err := e.ExecuteRules([]Rule{test.rule}, test.facts)
		if test.expectErr && err == nil {
			t.Errorf("expected an error, but didn't get one")
			continue
		}
		if err != nil && !test.expectErr {
			t.Errorf("got an unexpected error: %v", err)
			continue
		}

		if !reflect.DeepEqual(test.expectedResults, result) {
			t.Errorf("expected %+v, but got %+v", test.expectedResults, result)
		}
	}
}

type fakeClosableCheck struct {
	closeCalled bool
}

func (*fakeClosableCheck) Check() error { return nil }
func (*fakeClosableCheck) Name() string { return "" }

func (c *fakeClosableCheck) Close() error {
	c.closeCalled = true
	return nil
}

func TestEngineClosableCheck(t *testing.T) {
	fakeCheck := &fakeClosableCheck{}
	mapper := fakeRuleCheckMapper{
		check: fakeCheck,
	}
	e := Engine{
		RuleCheckMapper: mapper,
	}
	rule := fakeRule{}
	_, err := e.ExecuteRules([]Rule{rule}, []string{})
	if err != nil {
		t.Errorf("unexpected error when executing closable check: %v", err)
	}

	if err := e.CloseChecks(); err != nil {
		t.Errorf("unexpected error when closing checks: %v", err)
	}

	if !fakeCheck.closeCalled {
		t.Errorf("The check was not closed")
	}
}
