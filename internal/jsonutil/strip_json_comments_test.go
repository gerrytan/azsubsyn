package jsonutil_test

import (
	"encoding/json"
	"testing"

	"github.com/gerrytan/azsubsyn/internal/jsonutil"
)

func TestStripJSONComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "No comments",
			input: `{"name": "test", "value": 123}`,
			expected: `{"name": "test", "value": 123}
`,
		},
		{
			name: "Single line comment at end",
			input: `{
  "name": "test", // This is a comment
  "value": 123
}`,
			expected: `{
"name": "test",
"value": 123
}
`,
		},
		{
			name: "Single line comment on separate line",
			input: `{
  "name": "test",
  // This is a comment line
  "value": 123
}`,
			expected: `{
"name": "test",
"value": 123
}
`,
		},
		{
			name: "Multi-line comment",
			input: `{
  "name": "test",
  /* This is a
     multi-line comment */
  "value": 123
}`,
			expected: `{
"name": "test",
"value": 123
}
`,
		},
		{
			name: "Multi-line comment on single line",
			input: `{
  "name": "test", /* inline comment */
  "value": 123
}`,
			expected: `{
"name": "test",
"value": 123
}
`,
		},
		{
			name: "Comment inside string should be preserved",
			input: `{
  "name": "test with // comment",
  "description": "This has /* comment */ inside",
  "value": 123
}`,
			expected: `{
"name": "test with // comment",
"description": "This has /* comment */ inside",
"value": 123
}
`,
		},
		{
			name: "Mixed comments",
			input: `{
  // Top level comment
  "name": "test", // End of line comment
  /* Multi-line
     comment here */
  "value": 123, /* Another inline comment */
  "description": "String with // comment inside"
}`,
			expected: `{
"name": "test",
"value": 123,
"description": "String with // comment inside"
}
`,
		},
		{
			name: "Empty lines and whitespace",
			input: `{

  "name": "test",   // Comment with spaces
  
  "value": 123
  
}`,
			expected: `{
"name": "test",
"value": 123
}
`,
		},
		{
			name: "Escaped quotes in strings",
			input: `{
  "name": "test with \"quotes\" and // comment",
  "value": 123 // This is a real comment
}`,
			expected: `{
"name": "test with \"quotes\" and // comment",
"value": 123
}
`,
		},
		{
			name: "Complex nested JSON with comments",
			input: `{
  // Configuration object
  "config": {
    "name": "app", // Application name
    "settings": {
      /* Database configuration
         with multiple lines */
      "database": {
        "host": "localhost", // Default host
        "port": 5432
      }
    }
  },
  // Array with comments
  "items": [
    "first", // First item
    "second" /* Second item */
  ]
}`,
			expected: `{
"config": {
"name": "app",
"settings": {
"database": {
"host": "localhost",
"port": 5432
}
}
},
"items": [
"first",
"second"
]
}
`,
		},
		{
			name:     "Only comments",
			input:    `// Just a comment\n/* Another comment */`,
			expected: ``,
		},
		{
			name: "URL with double slashes in string",
			input: `{
  "url": "https://example.com/path", // This is a comment
  "protocol": "http://" // Another comment  
}`,
			expected: `{
"url": "https://example.com/path",
"protocol": "http://"
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := jsonutil.StripJSONComments([]byte(tt.input))

			// Normalize whitespace for comparison
			resultStr := string(result)

			if resultStr != tt.expected {
				t.Errorf("StripJSONComments() = %q, expected %q", resultStr, tt.expected)
			}

			// Additional test: verify the result is valid JSON if expected is not empty
			if tt.expected != "" {
				var jsonObj interface{}
				if err := json.Unmarshal(result, &jsonObj); err != nil {
					t.Errorf("StripJSONComments() produced invalid JSON: %v\nResult: %s", err, resultStr)
				}
			}
		})
	}
}

func TestStripJSONComments_ValidJSONOutput(t *testing.T) {
	// Test that common JSONC patterns produce valid JSON
	jsonc := `{
  // This is a configuration file
  "name": "azsubsyn-plan",
  "version": "1.0.0", // Version comment
  
  /* Resource provider registrations
     These will be applied to the target subscription */
  "rpRegistrations": [
    {
      "namespace": "Microsoft.Compute", // For VMs
      "reason": "NotRegisteredInTarget"
    },
    {
      "namespace": "Microsoft.Storage", /* For storage accounts */
      "reason": "NotFoundInTarget"
    }
  ]
}`

	result := jsonutil.StripJSONComments([]byte(jsonc))

	// Verify it's valid JSON
	var parsed map[string]interface{}
	err := json.Unmarshal(result, &parsed)
	if err != nil {
		t.Fatalf("Failed to parse result as JSON: %v\nResult: %s", err, string(result))
	}

	// Verify structure
	if parsed["name"] != "azsubsyn-plan" {
		t.Errorf("Expected name to be 'azsubsyn-plan', got %v", parsed["name"])
	}

	if parsed["version"] != "1.0.0" {
		t.Errorf("Expected version to be '1.0.0', got %v", parsed["version"])
	}

	rpRegistrations, ok := parsed["rpRegistrations"].([]interface{})
	if !ok {
		t.Fatalf("Expected rpRegistrations to be an array")
	}

	if len(rpRegistrations) != 2 {
		t.Errorf("Expected 2 rpRegistrations, got %d", len(rpRegistrations))
	}
}
