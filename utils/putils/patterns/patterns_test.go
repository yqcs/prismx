package patterns

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatterns(t *testing.T) {
	testCases := []struct {
		name     string
		pattern  *regexp.Regexp
		positive []string
		negative []string
	}{
		{
			name:    "URLRegexp",
			pattern: URLRegexp,
			positive: []string{
				"http://example.com",
				"https://example.net",
				"http://sub.example.com:8080",
				"https://www.example.com/path?query=param#fragment",
				"http://ex-ample.io",
				"https://123-example.net",
				"http://example.com/a_b/c~d#fragment",
			},
			negative: []string{
				// "http//example.com",
				"https:/example.com",
				"http:///example.com",
				// "http://ex_ample.com",
				"http://.com",
				"http://example..com",
				"http://example-.com",
				"http://-example.com",
				"http://example.com:-80",
				"http://example.com:800000",
				"http://",
			},
		},
		{
			name:    "DomainRegexp",
			pattern: DomainRegexp,
			positive: []string{
				"example.com",
				"mail.xyz123.org",
				"doc-s-2.example-4.gov",
				"m4u.museum",
				"golang.dev",
				"123-numbers-school.edu",
				"b1-d-l.net",
				"ai-learning.codes",
				"v3-0-coral.connect",
			},
			negative: []string{
				".example.com",
				"example.org.",
				"#oogle.z9*.net",
				"!example.com",
				"my_test@example.com",
				"example..com",
				"(.example{9}.org",
				"my+example4.travel",
				"abc&123.deg",
				"my:site.123.gov",
			},
		},
		{
			name:    "IPv4Regexp",
			pattern: IPv4Regexp,
			positive: []string{
				"123.45.67.89",
				"192.168.1.1",
				"10.0.0.1",
				"172.16.254.1",
				"1.2.3.4",
				"0.0.0.0",
				"255.255.255.255",
				"200.200.200.200",
				"50.238.2.98",
			},
			negative: []string{
				"256.1.1.1",
				"192.168.1.256",
				"999.999.999.999",
				"172.316.254.1",
				"abc.def.ghi.jkl",
				"192.168.1.",
				"192.168..1",
				"192..168.1",
				"..192.168",
			},
		},
		{
			name:    "URLPort80Regexp",
			pattern: URLPort80Regexp,
			positive: []string{
				"http://localhost:80",
				"http://example.com:80",
				"https://localhost:80",
				"localhost:80",
				"192.168.1.1:80",
				"www.example.com:80",
				"http://192.168.1.1:80",
				"https://192.168.1.1:80",
			},
			negative: []string{
				"http://localhost",
				"http://example.com:8080",
				"https://localhost:443",
				"example.com",
				"localhost:81",
				"192.168.1.1:1234",
				"www.example.com",
				"http://localhost:443",
				"http://192.168.1.1",
			},
		},
		{
			name:    "EmailRegexp",
			pattern: EmailRegexp,
			positive: []string{
				"john.doe@example.com",
				"mary_smith123@demo.net",
				"张三@测试.中国",
				"a@b.co",
				"first.last@sub.dom.ain",
				"email@sub-domain.com",
				"_user@domain.org",
				"123.456@numbers.domain",
				"user@domain.io",
			},
			negative: []string{
				"john.doe",
				"mary_smith123@.net",
				"@example.com",
				"张三测试.中国",
				// "a@b.c",
				// "first.last@sub..dom.ain",
				"email@.com",
				"_user@domain@org",
				"123.456@numbers@domain",
				"user@domain@io",
				"user..name@domain.io",
				"user@domain..io",
			},
		},
		{
			name:    "CIDRRegexp",
			pattern: CIDRRegexp,
			positive: []string{
				"192.168.1.1/24",
				"10.0.0.0/8",
				"172.16.0.0/12",
				"255.255.255.255/32",
				"192.0.2.0/24",
				"203.0.113.0/24",
				"198.51.100.0/24",
				"203.0.113.0/24",
			},
			negative: []string{
				"192.168.1.1",
				"10.0.0.0",
				"172.16.0.0",
				"255.255.255.255",
				"192.0.2.0.24",
				"203.0.113.0.24",
				"198.51.100.0.24",
				"203.0.113.0.24",
				"300.0.0.0/24",
				"192.168.1.1/33",
				"172.16.0.0/512",
			},
		},
		{
			name:    "CVEIdentifierRegexp",
			pattern: CVEIdentifierRegexp,
			positive: []string{
				"cve-2018-123456",
				"cve-2019-456789",
				"cve-2020-987654",
				"cve-2021-654321",
				"cve-2022-112233",
			},
			negative: []string{
				"cve-202022-654321", // invalid year format
				"cve2022-112233",    // missing hyphen
				"cve-2020-112233a",  // non-numeric ID
				"cve-2020-",         // missing ID
				"cv-2020-112233",    // invalid prefix
			},
		},
		{
			name:    "DualSpaceRegexp",
			pattern: DualSpaceRegexp,
			positive: []string{
				" ",
				"    ",
				"\t",
				"\n",
				" \t",
				"\t\n",
				"\n ",
				" \t\n",
				" hello world ",
				"\texample\t",
				"two  spaces",
				"\nnew\nline\n",
			},
			negative: []string{
				"",
				"a",
				"nospace",
				"period.after",
				"under_score",
				"dash-separated",
				"123",
				"abc123",
			},
		},
		{
			name:    "DecolorizerRegexp",
			pattern: DecolorizerRegexp,
			positive: []string{
				"\x1B[0m",       // Reset
				"\x1B[31m",      // Red
				"\x1B[1;31m",    // Bold Red
				"\x1B[4;33m",    // Underline Yellow
				"\x1B[44m",      // Background Blue
				"\x1B[1;44;33m", // Bold Yellow text on Blue background
				"\x1B[7;32m",    // Reverse video Green
				"\x1B[2J",       // Clear screen
				"\x1B[2K",       // Clear line
			},
			negative: []string{
				"Hello, World!",
				"12345",
				"\x1B",
				"\x1B[31",
				"ESC[1;32m",
				"\x1B[1;32",
				"\x1Bm",
				"\x1B[",
				"\x1B[;31",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, pos := range tc.positive {
				ok := tc.pattern.MatchString(pos)
				if !ok {
					t.Log(pos)
				}
				assert.True(t, tc.pattern.MatchString(pos))
			}
			for _, neg := range tc.negative {
				ok := tc.pattern.MatchString(neg)
				if ok {
					t.Log(neg)
				}
				assert.False(t, tc.pattern.MatchString(neg))
			}
		})
	}
}
