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

// PatchElements processes patch instructions in an XML element to modify a target XML document.
func PatchElements(patchEl *etree.Element, newDoc *etree.Document) {
	// Check for nil inputs and log if detected
	if patchEl == nil {
		return
	}

	// Process the patch element based on its namespace (either "add" or "del")
	if patchEl.Space == "add" || patchEl.Space == "del" {
		InjectElementAtPath(patchEl, newDoc)
	}

	// Process patch attributes in the same way
	for _, attr := range patchEl.Attr {
		if attr.Space == "add" || attr.Space == "del" {
			InjectElementAtPath(patchEl, newDoc)
		}
	}

	// Recursively handle child elements
	for _, child := range patchEl.ChildElements() {
		PatchElements(child, newDoc)
	}
}

// InjectElementAtPath either adds or deletes elements and attributes from the target document
// based on the patch instruction.
func InjectElementAtPath(el *etree.Element, doc *etree.Document) {
	// Try to find the element in the target document
	match := doc.FindElement(el.GetPath())

	// If no matching element found, create necessary path
	if match == nil {
		current := doc.Root()
		parts := strings.Split(el.GetPath(), "/")

		// Traverse or create the path in the target document
		for i := 2; i < len(parts); i++ {
			match = current.SelectElement(parts[i])
			if match == nil {
				match = current.CreateElement(parts[i])
			}
			current = match
		}

		// If adding, set the text content of the element
		if el.Space == "add" {
			match.SetText(el.Text())
		}
	} else {
		// If element is found and is being added, set its text content
		if el.Space == "add" {
			match.SetText(el.Text())
		} else if el.Space == "del" {
			// If element is being deleted, remove it from its parent
			match.Parent().RemoveChild(match)
		}
	}

	// Apply attribute patches
	for _, attr := range el.Attr {
		if attr.Space == "add" {
			match.CreateAttr(attr.Key, attr.Value)
		} else if attr.Space == "del" {
			match.RemoveAttr(attr.Key)
		}
	}
}
