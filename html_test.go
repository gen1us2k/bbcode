// Copyright 2015 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

var urlTests = map[string]string{
	"http://example.com/path?query=value#fragment":         "http://example.com/path?query=value#fragment",
	"<script>http://example.com":                           "%3Cscript%3Ehttp://example.com",
	"http://example.com/path?query=value#fragment<script>": "http://example.com/path?query=value#fragment%3Cscript%3E",
	"http://example.com/path?query=<script>":               "http://example.com/path?query=<script>",
	"javascript:alert(1);":                                 "javascript:alert(1);",
}
