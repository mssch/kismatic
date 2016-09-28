package inspector

type mockCheck struct {
	err error
}

func (c mockCheck) Check() error {
	return c.err
}
func (c mockCheck) Name() string {
	return "mockCheck"
}

//
// func TestEngine(t *testing.T) {
// 	// Setup some rules for the tests
// 	failRule := Rule{
// 		Name:  "failRule",
// 		check: mockCheck{errors.New("error occurred")},
// 	}
// 	passingRule := Rule{
// 		Name:  "successRule",
// 		check: mockCheck{nil},
// 	}
// 	// Setup test cases
// 	tests := []struct {
// 		rules           []Rule
// 		expectedResults []RuleResult
// 		labels          []string
// 	}{
// 		// A single rule that passes
// 		{
// 			rules: []Rule{
// 				passingRule,
// 			},
// 			expectedResults: []RuleResult{
// 				{
// 					Rule:    passingRule,
// 					Success: true,
// 					Error:   nil,
// 				},
// 			},
// 		},
// 		// A single rule that fails
// 		{
// 			rules: []Rule{
// 				failRule,
// 			},
// 			expectedResults: []RuleResult{
// 				{
// 					Rule:    failRule,
// 					Success: false,
// 					Error:   failRule.check.Check(),
// 				},
// 			},
// 		},
// 		// One passing rule, one failing rule
// 		{
// 			rules: []Rule{
// 				passingRule,
// 				failRule,
// 			},
// 			expectedResults: []RuleResult{
// 				{
// 					Rule:    passingRule,
// 					Success: true,
// 					Error:   nil,
// 				},
// 				{
// 					Rule:    failRule,
// 					Success: false,
// 					Error:   failRule.check.Check(),
// 				},
// 			},
// 		},
// 		// Single rule that should run because of labels
// 		{
// 			labels: []string{"centos", "master", "otherLabel"},
// 			rules: []Rule{
// 				{
// 					Name:   "PassWithLabels",
// 					check:  mockCheck{},
// 					Labels: []string{"centos"},
// 				},
// 			},
// 			expectedResults: []RuleResult{
// 				{
// 					Rule: Rule{
// 						Name:   "PassWithLabels",
// 						check:  mockCheck{},
// 						Labels: []string{"centos"},
// 					},
// 					Success: true,
// 					Error:   nil,
// 				},
// 			},
// 		},
// 		// Single rule that should not run because of labels
// 		{
// 			labels: []string{"master", "otherLabel"},
// 			rules: []Rule{
// 				{
// 					Name:   "PassWithLabels",
// 					check:  mockCheck{},
// 					Labels: []string{"centos"},
// 				},
// 			},
// 			expectedResults: []RuleResult{},
// 		},
// 		// Single rule that should run because it has no labels
// 		{
// 			labels: []string{"master", "centos"},
// 			rules: []Rule{
// 				{
// 					Name:   "passNoLabels",
// 					check:  mockCheck{},
// 					Labels: []string{},
// 				},
// 			},
// 			expectedResults: []RuleResult{
// 				{
// 					Rule: Rule{
// 						Name:   "passNoLabels",
// 						check:  mockCheck{},
// 						Labels: []string{},
// 					},
// 					Success: true,
// 				},
// 			},
// 		},
// 	}
//
// 	for _, test := range tests {
// 		e := engine{}
// 		res := e.executeRules(test.rules, test.labels)
// 		if !reflect.DeepEqual(test.expectedResults, res) {
// 			t.Errorf("Expected results not equal to obtained results. Expected:\n%+v. Got:\n%+v", test.expectedResults, res)
// 		}
// 	}
// }
