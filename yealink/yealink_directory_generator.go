/*
 * Copyright 2019. Improcom Inc
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 *  IMPROCOM, Inc
 *  1803 Gravesend Neck road
 *  Brooklyn, NY, 11229 USA
 *  Phone: +1-718-769-3400
 *  @author Alex Romanchuk - alexr@improcom.com - Improcom INC
 */

// Package yealink provides generators for yealink phones directory config
package yealink

import (
	"fmt"
	"io"
	"strings"
)

type PhoneDirectoryEntry interface {
}

//type DirectoryEntry interface {
//	Name() string
//	Phone1() string
//	Phone2() string
//	Phone3() string
//}

type PhoneDirectoryConfig struct {
	MenuName string
	Entries  []ConfigDirectoryEntry
}

type ConfigDirectoryEntry struct {
	Name   string
	Phone1 string
	Phone2 string
	Phone3 string
}

func GenerateXml(phConfig PhoneDirectoryConfig) (xmlContent strings.Builder) {
	fmt.Fprintf(&xmlContent, `<?xml version="1.0" encoding="UTF-8"?>
<YealinkIPPhoneBook>
  <Title>Yealink</Title>
  <Menu Name="Sales">
`)

	for _, contactEntry := range phConfig.Entries {
		fmt.Fprintf(&xmlContent, `<Unit Name="%v" Phone1="%v" Phone2="%v" Phone3="%v" default_photo="Resource:"/>
`, contactEntry.Name, contactEntry.Phone1, contactEntry.Phone2, contactEntry.Phone3)
	}

	fmt.Fprintf(&xmlContent, `</Menu>
</YealinkIPPhoneBook>`)

	return xmlContent
}

func processTopic(w io.Writer, id string, properties map[string][]string) {
	fmt.Fprintf(w, "<card entity=\"%s\">\n", id)
	fmt.Fprintln(w, "  <facts>")
	for k, v := range properties {
		for _, value := range v {
			fmt.Fprintf(w, "    <fact property=\"%s\">%s</fact>\n", k, value)
		}
	}
	fmt.Fprintln(w, "  </facts>")
	fmt.Fprintln(w, "</card>")
}
