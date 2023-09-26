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

// FocusEtree returns an etree.Element that represents the specified path in the document.
// It removes all other branches above the element that are not on the path
// If the path does not exist, it returns nil.
func FocusEtree(doc *etree.Document, path string) *etree.Element {
	path = strings.TrimPrefix(path, "/")

	// Find all elements that match the path
	foundElements := doc.FindElements(path)
	if len(foundElements) == 0 {
		Log(1, "Xpath element \"%s\" does not exist", path)
		return nil
	}

	// Create a new element to represent the focused path
	parts := strings.Split(path, "/")
	focused := etree.NewElement(parts[0])

	// Get the space of the found element
	space := foundElements[0].Space
	depth := len(parts)
	if depth > 1 {
		// Create child elements for each part of the path
		parts = parts[:depth-1]
		current := focused
		for i := 1; i < len(parts); i++ {
			newElem := current.CreateElement(parts[i])
			// Find the element in the document and copy its attributes
			element := doc.FindElement(strings.Join(parts[:i+1], "/"))
			space = element.Space
			if space != "" {
				newElem.Space = space
			}
			if element != nil {
				for _, attr := range element.Attr {
					newElem.CreateAttr(attr.Key, attr.Value)
				}
			}
			current = newElem
		}
		// Add all found elements as children of the last child element
		for _, foundElement := range foundElements {
			current.AddChild(foundElement.Copy())
		}
	} else {
		// If the path is just the root element, return the root element of the document
		focused = doc.Root()
	}
	if space != "" {
		// Set the space of the focused element to "att"
		focused.Space = "att"
		Log(5, "element maked with attention flag: %s", focused.GetPath())
	}
	return focused
}
