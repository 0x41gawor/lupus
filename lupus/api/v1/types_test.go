package v1

import (
	"errors"
	"strings"
	"testing"
)

func TestCompare(t *testing.T) {
	tests := []struct {
		name      string
		condition Condition
		field     interface{}
		expected  bool
		expectErr error
	}{
		{
			name: "equal integers",
			condition: Condition{
				Operator: "eq",
				Value:    42,
			},
			field:    42,
			expected: true,
		},
		{
			name: "not equal integers",
			condition: Condition{
				Operator: "ne",
				Value:    42,
			},
			field:    30,
			expected: true,
		},
		{
			name: "greater than integer",
			condition: Condition{
				Operator: "gt",
				Value:    10,
			},
			field:    20,
			expected: true,
		},
		{
			name: "less than integer",
			condition: Condition{
				Operator: "ls",
				Value:    20,
			},
			field:    10,
			expected: true,
		},
		{
			name: "equal strings",
			condition: Condition{
				Operator: "eq",
				Value:    "hello",
			},
			field:    "hello",
			expected: true,
		},
		{
			name: "not equal strings",
			condition: Condition{
				Operator: "ne",
				Value:    "hello",
			},
			field:    "world",
			expected: true,
		},
		{
			name: "equal booleans",
			condition: Condition{
				Operator: "eq",
				Value:    true,
			},
			field:    true,
			expected: true,
		},
		{
			name: "not equal booleans",
			condition: Condition{
				Operator: "ne",
				Value:    false,
			},
			field:    true,
			expected: true,
		},
		{
			name: "type mismatch with conversion: string to int",
			condition: Condition{
				Operator: "eq",
				Value:    "42",
			},
			field:    42,
			expected: true,
		},
		{
			name: "type mismatch with invalid conversion",
			condition: Condition{
				Operator: "eq",
				Value:    "invalid",
			},
			field:     42,
			expectErr: errors.New("type mismatch and conversion failed"),
		},
		{
			name: "unsupported operator",
			condition: Condition{
				Operator: "unsupported",
				Value:    42,
			},
			field:     42,
			expectErr: errors.New("unsupported operator: unsupported"),
		},
		{
			name: "unsupported field type",
			condition: Condition{
				Operator: "eq",
				Value:    42,
			},
			field:     []int{1, 2, 3},
			expectErr: errors.New("unsupported field type: []int"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.condition.Compare(test.field)

			if test.expectErr != nil {
				if test.expectErr != nil {
					if test.expectErr != nil {
						if err == nil {
							t.Fatalf("expected error: %v, got: nil", test.expectErr)
						}
						if !containsError(err, test.expectErr) {
							t.Fatalf("expected error containing: %v, got: %v", test.expectErr, err)
						}
					}

				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result != test.expected {
					t.Fatalf("expected result: %v, got: %v", test.expected, result)
				}
			}
		})
	}
}

func containsError(actual error, expected error) bool {
	return actual != nil && expected != nil &&
		strings.Contains(actual.Error(), expected.Error())
}
