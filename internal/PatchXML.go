/*
Copyright © 2023 Miha miha.kralj@outlook.com

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
package internal

import (
	"strings"

	"github.com/beevik/etree"
)

func PatchElements(patchEl *etree.Element, newDoc *etree.Document) {
	if patchEl == nil {
		return
	}

	// Process elements
	switch patchEl.Space {
	case "add", "del":
		InjectElementAtPath(patchEl, newDoc)
	}

	// Process attributes
	for _, attr := range patchEl.Attr {
		switch attr.Space {
		case "add", "del":
			InjectElementAtPath(patchEl, newDoc)
		}
	}
	// Recursively process child elements
	for _, child := range patchEl.ChildElements() {
		PatchElements(child, newDoc)
	}
}

func InjectElementAtPath(el *etree.Element, doc *etree.Document) {

	match := doc.FindElement(el.GetPath())

	if match == nil {
		current := doc.Root()
		parts := strings.Split(el.GetPath(), "/")
		// No match found in doc, we need to create the new path
		for i := 2; i < len(parts); i++ {
			match = current.SelectElement(parts[i])
			if match == nil {
				match = current.CreateElement(parts[i])
			}
			current = match
		}
		if el.Space == "add" {
			match.SetText(el.Text())
		}
	} else {
		if el.Space == "add" {
			match.SetText(el.Text())
		}
		if el.Space == "del" {
			match.Parent().RemoveChild(match)
		}
	}

	for _, attr := range el.Attr {
		if attr.Space == "add" {
			match.CreateAttr(attr.Key, attr.Value)
		}
		if attr.Space == "del" {
			match.RemoveAttr(attr.Key)
		}

	}

}
