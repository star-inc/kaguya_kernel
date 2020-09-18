/*
Package kaguya : The library for kaguya

    Copyright 2020 Star Inc.(https://starinc.xyz)

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/
package kaguya

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// HTTPGet : Get WWW resources from Internet
func HTTPGet(url string, recovery int) string {
	var output string
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	DeBug("GenRequest", err)
	resp, err := client.Do(req)
	if err == nil {
		body, err := ioutil.ReadAll(resp.Body)
		DeBug("ReadHTML", err)
		output = string(body)
	} else {
		if recovery == 0 {
			time.Sleep(time.Duration(2) * time.Second)
			HTTPGet(url, 1)
		} else {
			DeBug("GetHTTP", err)
		}
	}
	resp.Body.Close()
	return output
}

// DeBug : Print errors for debug and report
func DeBug(where string, err error) bool {
	if err != nil {
		log.Printf("Kaguya Error #%s\nReason:\n%s\n\n", where, err)
		return false
	}
	return true
}

// ReplaceSyntaxs : Remove space and syntax
func ReplaceSyntaxs(rawString string, filled string) string {
	var output bytes.Buffer
	rawString = strings.ReplaceAll(rawString, " ", "\x1e")
	rawString = strings.ReplaceAll(rawString, "\t", "\x1e")
	rawString = strings.ReplaceAll(rawString, "\n", "\x1e")
	stringSlice := strings.Split(rawString, "\x1e")
	for _, word := range stringSlice {
		if word != "" {
			output.WriteString(word + filled)
		}
	}
	return output.String()
}

// RemoveChildNode : Remove all child html node selected
func RemoveChildNode(rootNode *html.Node, removeNode *html.Node) {
	foundNode := false
	checkNode := make(map[int]*html.Node)

	i := 0
	for n := rootNode.FirstChild; n != nil; n = n.NextSibling {
		if n == removeNode {
			foundNode = true
			n.Parent.RemoveChild(n)
		}

		checkNode[i] = n
		i++
	}

	if !foundNode {
		for _, item := range checkNode {
			RemoveChildNode(item, removeNode)
		}
	}
}
