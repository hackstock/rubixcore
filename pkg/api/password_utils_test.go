package api

import "testing"

func TestHashPassword(t *testing.T) {
	testCases := []struct {
		tag           string
		plainPassword string
	}{
		{
			tag:           "valid case",
			plainPassword: "foobar",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tag, func(t *testing.T) {
			_, err := hashPassword(tc.plainPassword)
			if err != nil {
				t.Fatalf("expected no error hashing password, got %v", err)
			}
		})
	}
}

func TestComparePasswords(t *testing.T) {
	testCases := []struct {
		tag            string
		hashedPassword string
		plainPassword  string
		expect         bool
	}{
		{
			tag:            "valid case",
			hashedPassword: "$2a$10$pJofeBaFtdXo4RdRrKBJF.FW/ePvnS3.xgNpdC0N4FNt2S1H3QO2K",
			plainPassword:  "foobar",
			expect:         true,
		},
		{
			tag:            "invalid case",
			hashedPassword: "$2a$10$pJofeBaFtdXo4RdRrKBJF.FW/ePvnS3.xgNpdC0N4FNt2S1H3QO2K",
			plainPassword:  "somepassword",
			expect:         false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tag, func(t *testing.T) {
			got := comparePasswords(tc.hashedPassword, tc.plainPassword)
			if got != tc.expect {
				t.Fatalf("expected %v, got %v", tc.expect, got)
			}
		})
	}
}
