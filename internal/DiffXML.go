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
	"fmt"
	"strings"

	"github.com/beevik/etree"
)

// DiffXML compares two etree documents and returns a new document with only the changed elements.
func DiffXML(oldDoc, newDoc *etree.Document, fulltree bool) *etree.Document {
	diffDoc := oldDoc.Copy()

	EnumerateListElements(newDoc.Root())
	EnumerateListElements(diffDoc.Root())

	addMissingElements(newDoc.Root(), diffDoc)
	checkElements(diffDoc.Root(), newDoc)

	ReverseEnumerateListElements(diffDoc.Root())
	ReverseEnumerateListElements(newDoc.Root())

	if !fulltree {
		removeNodesWithoutSpace(diffDoc.Root())
		removeAttSpace(diffDoc.Root())
	}
	return diffDoc
}

// removeNodesWithoutSpace recursively removes elements without a "Space" attribute
func removeNodesWithoutSpace(el *etree.Element) {
	for i := 0; i < len(el.Child); i++ {
		child, ok := el.Child[i].(*etree.Element)
		if !ok {
			continue
		}

		// Check if any attribute has Space defined
		hasAttrWithSpace := false
		for _, attr := range child.Attr {
			if attr.Space != "" {
				hasAttrWithSpace = true
				break
			}
		}

		// Remove the child element only if it doesn't have a Space and none of its attributes have a Space
		if child.Space == "" && !hasAttrWithSpace {
			el.RemoveChildAt(i)
			i-- // Adjust index because we've removed an item
			continue
		}

		removeNodesWithoutSpace(child)
	}
}

func removeAttSpace(el *etree.Element) {
	if el == nil {
		return
	}

	// Remove or unset the "Space" attribute if it is set to "att"
	if el.Space == "att" {
		el.Space = ""
	}

	// Recursively process children
	for i := 0; i < len(el.Child); i++ {
		child, ok := el.Child[i].(*etree.Element)
		if !ok {
			continue // Skip if this child is not an Element
		}

		removeAttSpace(child)
	}
}

func RemoveChgSpace(el *etree.Element) {
	if el == nil {
		return
	}

	// Remove or unset the "Space" attribute if it is set to "att"
	if el.Space == "chg" {
		parts := strings.Split(el.Text(), "|||")
		if len(parts) > 1 {
			el.SetText(parts[1])
		}
		el.Space = "add"
	}
    // Process attributes
	for i := range el.Attr {
		// Check if the attribute space is "chg"
		if el.Attr[i].Space == "chg" {
			parts := strings.Split(el.Attr[i].Value, "|||")
			if len(parts) > 1 {
				el.Attr[i].Value = parts[1]
			}
			el.Attr[i].Space = "add"
		}
	}

	// Recursively process children
	for i := 0; i < len(el.Child); i++ {
		child, ok := el.Child[i].(*etree.Element)
		if !ok {
			continue // Skip if this child is not an Element
		}

		RemoveChgSpace(child)
	}
}

func checkElements(oldEl *etree.Element, newDoc *etree.Document) {
	newEl := newDoc.FindElement(oldEl.GetPath())
	if newEl != nil {
		// Element found in newDoc
		newElText := strings.TrimSpace(newEl.Text())
		oldElText := strings.TrimSpace(oldEl.Text())

		if newElText != oldElText {
			if newElText != "" && oldElText != "" {
				oldEl.Space = "chg"
				oldEl.SetText(fmt.Sprintf("%s|||%s", oldElText, newElText))
				markParentSpace(oldEl)
			} else if newElText != "" && oldElText == ""{
				oldEl.Space = "chg"
				oldEl.SetText("N/A|||"+newEl.Text())
				markParentSpace(oldEl)
			}
		}
		copyAttributes(newEl, oldEl)

		// Check comments
		checkComments(oldEl, newEl)
	} else {
		oldEl.Space = "del"
		markParentSpace(oldEl)
	}

	// Recursively check all child elements
	for _, child := range oldEl.ChildElements() {
		checkElements(child, newDoc)
	}
}

func checkComments(oldEl, newEl *etree.Element) {
	oldComments := getComments(oldEl)
	newComments := getComments(newEl)

	for _, oldComment := range oldComments {
		if !containsComment(newComments, oldComment) {
			updateComment(oldEl, "del:"+oldComment)
		}
	}

	for i, newComment := range newComments {
		if !containsComment(oldComments, newComment) {
			newCommentNode := etree.NewComment("new:" + newComment)
			oldEl.InsertChildAt(i, newCommentNode)
		}
	}
}

func addMissingElements(newEl *etree.Element, oldDoc *etree.Document) {
	oldEl := oldDoc.FindElement(newEl.GetPath())
	if oldEl == nil {
		// Element not found in oldDoc
		parentPath := newEl.Parent().GetPath()
		parentInOldDoc := oldDoc.FindElement(parentPath)
		if parentInOldDoc != nil {

			oldEl := etree.NewElement(fmt.Sprintf("add:%s", newEl.Tag))
			oldEl.SetText(newEl.Text())
			copyAttributes(newEl, oldEl)

			parentInOldDoc.AddChild(oldEl)
			addedchild := parentInOldDoc.Child[len(parentInOldDoc.Child)-1]

			markParentSpace(addedchild.(*etree.Element))
		}
	}

	// Recursively check all child elements
	for _, child := range newEl.ChildElements() {
		addMissingElements(child, oldDoc)
	}
}

func copyAttributes(oldEl, newEl *etree.Element) {
	// Check if oldEl or newEl is nil
	if oldEl == nil || newEl == nil {
		return
	}

	// Check attributes in oldEl
	for _, oldAttr := range oldEl.Attr {
		newAttr := newEl.SelectAttr(oldAttr.Key)
		if newAttr != nil {
			// Attribute exists in newEl
			if newAttr.Value != oldAttr.Value {
				// Different value, add chg: in front of attribute name
				newEl.RemoveAttr(oldAttr.Key)
				newEl.CreateAttr(fmt.Sprintf("chg:%s", oldAttr.Key), fmt.Sprintf("%s|||%s", newAttr.Value, oldAttr.Value))
				markParentSpace(newEl)
			}
			// If same value, do nothing
		} else {
			// Attribute does not exist in newEl, add with namespace del:
			newEl.CreateAttr(fmt.Sprintf("add:%s", oldAttr.Key), oldAttr.Value)
			markParentSpace(newEl)
		}
	}

	// Create a copy of newEl.Attr
	newAttrs := make([]etree.Attr, len(newEl.Attr))
	copy(newAttrs, newEl.Attr)

	// Check attributes in newEl
	for _, newAttr := range newAttrs {
		oldAttr := oldEl.SelectAttr(newAttr.Key)
		if oldAttr == nil {
			// Attribute does not exist in oldEl, add with namespace new:
			newEl.RemoveAttr(newAttr.Key)
			newEl.CreateAttr(fmt.Sprintf("del:%s", newAttr.Key), strings.TrimSpace(newAttr.Value))
			markParentSpace(newEl)
		}
		// If attribute exists in oldEl, it has already been handled
	}
}

func EnumerateListElements(el *etree.Element) {
	childElementCounts := make(map[string]int)
	childElements := el.ChildElements()

	// Count occurrences of each tag
	for _, child := range childElements {
		childElementCounts[child.Tag]++
	}

	// Rename elements with duplicate tags
	for tag, count := range childElementCounts {
		if count > 1 {
			var index = 1
			for _, child := range childElements {
				if child.Tag == tag {
					child.Tag = fmt.Sprintf("%s.%d", tag, index)
					index++
				}
				EnumerateListElements(child)
			}
		} else {
			for _, child := range childElements {
				if child.Tag == tag {
					EnumerateListElements(child)
				}
			}
		}
	}
}

func ReverseEnumerateListElements(el *etree.Element) {
	childElements := el.ChildElements()

	// Iterate over child elements
	for _, child := range childElements {
		// Check if the tag contains a dot
		if strings.Contains(child.Tag, ".") {
			// Split the tag on the dot and take the first part
			parts := strings.Split(child.Tag, ".")
			child.Tag = parts[0]
		}
		// Recursively call the function on the child
		ReverseEnumerateListElements(child)
	}
}

func getComments(el *etree.Element) []string {
	var comments []string
	for _, token := range el.Child {
		if comment, ok := token.(*etree.Comment); ok {
			comments = append(comments, comment.Data)
		}
	}
	return comments
}

func containsComment(comments []string, comment string) bool {
	for _, c := range comments {
		if c == comment {
			return true
		}
	}
	return false
}

func getCommentString(el *etree.Element) string {
	commentStr := ""
	for _, token := range el.Child {
		if comment, ok := token.(*etree.Comment); ok {
			commentStr = comment.Data
		}
	}
	return commentStr
}

func updateComment(el *etree.Element, newCommentStr string) {
	for _, child := range el.Child {
		if comment, ok := child.(*etree.Comment); ok {
			comment.Data = newCommentStr
			break
		}
	}
}

func markParentSpace(el *etree.Element) {
	if el == nil {
		return
	}
	parent := el.Parent()
	if parent != nil && parent.Space == "" {
		parent.Space = "att"
		markParentSpace(parent)
	}
}
