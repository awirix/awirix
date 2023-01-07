package html

import (
	"github.com/awirix/awirix/luadoc"
)

func Lib() *luadoc.Lib {
	selector := &luadoc.Param{
		Name:        "selector",
		Type:        luadoc.String,
		Description: `The CSS selector to use to find the elements.`,
	}

	selection := &luadoc.Param{
		Name:        "selection",
		Type:        selectionTypeName,
		Description: "A selection object.",
	}

	classSelection := &luadoc.Class{
		Name: selectionTypeName,
		Methods: []*luadoc.Method{
			{
				Name:        "find",
				Description: `Finds all elements matching the given selector.`,
				Value:       selectionFind,
				Params:      []*luadoc.Param{selector},
				Returns:     []*luadoc.Param{selection},
			},
			{
				Name:        "each",
				Description: `Iterates over all elements in the selection, calling the given function for each one.`,
				Value:       selectionEach,
				Params: []*luadoc.Param{
					{
						Name:        "fn",
						Description: `The function to call for each element.`,
						Type: luadoc.Func{
							Params: []*luadoc.Param{
								selection,
								{
									Name: "index",
									Type: luadoc.Number,
								},
							},
						}.AsType(),
					},
				},
			},
			{
				Name:        "text",
				Description: `Returns the text contents of the first element in the selection.`,
				Value:       selectionText,
				Returns: []*luadoc.Param{
					{
						Name:        "text",
						Description: `The text contents of the first element in the selection.`,
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "html",
				Description: `Returns the HTML contents of the first element in the selection.`,
				Value:       selectionHtml,
				Returns: []*luadoc.Param{
					{
						Name:        "html",
						Description: `Gets the HTML contents of the first element in the set of matched elements. It includes text and comment nodes.`,
						Type:        luadoc.String,
					},
					{
						Name:        "error",
						Description: `An error if one occurred.`,
						Type:        luadoc.String,
						Optional:    true,
					},
				},
			},
			{
				Name:        "first",
				Description: `Reduces the set of matched elements to the first in the set. It returns a new selection object, and an empty selection object if the the selection is empty.`,
				Value:       selectionFirst,
				Returns: []*luadoc.Param{
					selection,
				},
			},
			{
				Name:        "last",
				Description: `Reduces the set of matched elements to the last in the set. It returns a new selection object, and an empty selection object if the the selection is empty.`,
				Value:       selectionLast,
				Returns: []*luadoc.Param{
					selection,
				},
			},
			{
				Name:        "parent",
				Description: `Gets the parent of each element in the Selection. It returns a new Selection object containing the matched elements.`,
				Value:       selectionParent,
				Returns: []*luadoc.Param{
					selection,
				},
			},
			{
				Name:        "eq",
				Description: `Reduces the set of matched elements to the one at the specified index. If a negative index is given, it counts backwards starting at the end of the set. It returns a new Selection object, and an empty Selection object if the index is invalid.`,
				Value:       selectionEq,
				Params: []*luadoc.Param{
					{
						Name:        "index",
						Description: `The index of the element to select.`,
						Type:        luadoc.Number,
					},
				},
				Returns: []*luadoc.Param{
					selection,
				},
			},
			{
				Name:        "attr",
				Description: `Returns the value of the given attribute of the first element in the selection.`,
				Value:       selectionAttr,
				Params: []*luadoc.Param{
					{
						Name:        "name",
						Description: `The name of the attribute to get.`,
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "value",
						Description: `The value of the attribute, or nil if the attribute is not present.`,
						Type:        luadoc.String,
					},
					{
						Name:        "ok",
						Description: `Whether the attribute was present.`,
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "attr_or",
				Description: `Returns the value of the given attribute of the first element in the selection, or a default value if the attribute is not present.`,
				Value:       selectionAttrOr,
				Params: []*luadoc.Param{
					{
						Name:        "name",
						Description: `The name of the attribute to get.`,
						Type:        luadoc.String,
					},
					{
						Name:        "default",
						Description: `The default value to return if the attribute is not present.`,
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "value",
						Description: `The value of the attribute, or the default value if the attribute is absent.`,
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "has_class",
				Description: `determines whether any of the matched elements are assigned the given class.`,
				Value:       selectionHasClass,
				Params: []*luadoc.Param{
					{
						Name:        "class",
						Description: `The class to check for.`,
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "ok",
						Description: `Whether any of the matched elements have the given class.`,
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "add_class",
				Description: `Adds the specified class(es) to each of the set of matched elements.`,
				Value:       selectionAddClass,
				Params: []*luadoc.Param{
					{
						Name:        "class",
						Description: `One or more class names to be added to the class attribute of each matched element.`,
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "remove_class",
				Description: `Removes the given class(es) from each element in the set of matched elements. Multiple class names can be specified, separated by a space or via multiple arguments. If no class name is provided, all classes are removed.`,
				Value:       selectionRemoveClass,
				Params: []*luadoc.Param{
					{
						Name:        "class",
						Description: "The class to remove. If not provided, all classes are removed.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "toggle_class",
				Description: `Adds or removes the given class(es) for each element in the set of matched elements. Multiple class names can be specified, separated by a space or via multiple arguments.`,
				Value:       selectionToggleClass,
				Params: []*luadoc.Param{
					{
						Name:        "class",
						Description: `The class to toggle.`,
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "next",
				Description: `Gets the immediately following sibling of each element in the set of matched elements, optionally filtered by a selector. It returns a new Selection object containing the matched elements.`,
				Value:       selectionNext,
				Returns:     []*luadoc.Param{selection},
			},
			{
				Name:        "next_all",
				Description: `Gets all the following siblings of each element in the Selection. It returns a new Selection object containing the matched elements.`,
				Value:       selectionNextAll,
			},
			{
				Name:        "next_until",
				Description: `gets all following siblings of each element up to but not including the element matched by the selector. It returns a new Selection object containing the matched elements.`,
				Value:       selectionNextUntil,
				Params: []*luadoc.Param{
					selector,
				},
				Returns: []*luadoc.Param{
					selection,
				},
			},
			{
				Name:        "prev",
				Description: `Gets the immediately preceding sibling of each element in the Selection. It returns a new selection object containing the matched elements.`,
				Value:       selectionPrev,
				Returns: []*luadoc.Param{
					selection,
				},
			},
			{
				Name:        "prev_all",
				Description: `Gets all the preceding siblings of each element in the Selection. It returns a new selection object containing the matched elements.`,
				Value:       selectionPrevAll,
				Returns: []*luadoc.Param{
					selection,
				},
			},
			{
				Name:        "prev_until",
				Description: `Gets all preceding siblings of each element up to but not including the element matched by the selector. It returns a new selection object containing the matched elements.`,
				Value:       selectionPrevUntil,
				Params: []*luadoc.Param{
					selector,
				},
				Returns: []*luadoc.Param{
					selection,
				},
			},
			{
				Name:        "siblings",
				Description: `Gets all sibling elements of each element in the Selection. It returns a new selection object containing the matched elements.`,
				Value:       selectionSiblings,
				Returns: []*luadoc.Param{
					selection,
				},
			},
			{
				Name:        "children",
				Description: `Gets all direct child elements of each element in the Selection. It returns a new selection object containing the matched elements.`,
				Value:       selectionChildren,
				Returns: []*luadoc.Param{
					selection,
				},
			},
			{
				Name:        "contents",
				Description: `Contents gets the children of each element in the selection, including text and comment nodes. It returns a new selection object containing these elements`,
				Value:       selectionContents,
				Returns: []*luadoc.Param{
					selection,
				},
			},
			{
				Name:        "filter",
				Description: `Filter reduces the set of matched elements to those that match the selector string. It returns a new Selection object for this subset of matching elements.`,
				Value:       selectionFilter,
				Params:      []*luadoc.Param{selector},
				Returns:     []*luadoc.Param{selection},
			},
			{
				Name:        "remove",
				Description: `Removes elements from the selection that match the selector string. It returns a new selection object with the matching elements removed.`,
				Value:       selectionRemove,
				Params:      []*luadoc.Param{selector},
				Returns:     []*luadoc.Param{selection},
			},
			{
				Name:        "is",
				Description: `Checks the current matched set of elements against a selector and returns true if at least one of these elements matches.`,
				Value:       selectionIs,
				Params:      []*luadoc.Param{selector},
				Returns: []*luadoc.Param{
					{
						Name:        "ok",
						Description: `Whether any of the matched elements match the selector.`,
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "find_selection",
				Description: `gets the descendants of each element in the current selection, filtered by a selection. It returns a new selection object containing these matched elements.`,
				Value:       selectionFindSelection,
				Params:      []*luadoc.Param{selection},
				Returns:     []*luadoc.Param{selection},
			},
			{
				Name:        "map",
				Description: `Iterates over a selection, executing a function for each matched element. The function's return value is added to the returned table.`,
				Value:       selectionMap,
				Params: []*luadoc.Param{
					{
						Name:        "fn",
						Description: `The function to execute for each element. It receives the index of the element in the selection and the element as arguments.`,
						Type: luadoc.Func{
							Params: []*luadoc.Param{
								selection,
								{
									Name:        "index",
									Description: `The index of the element in the selection.`,
									Type:        luadoc.Number,
								},
							},
							Returns: []*luadoc.Param{
								{
									Name:        "value",
									Description: `The value to add to the returned table.`,
									Type:        luadoc.Any,
								},
							},
						}.AsType(),
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "results",
						Description: `A table containing the return values of the function for each element.`,
						Type:        luadoc.List(luadoc.Any),
					},
				},
			},
			{
				Name:        "closest",
				Description: `Gets the first element that matches the selector by testing the element itself and traversing up through its ancestors in the DOM tree.`,
				Value:       selectionClosest,
				Params:      []*luadoc.Param{selector},
				Returns:     []*luadoc.Param{selection},
			},
			{
				Name:        "parents",
				Description: `Gets the ancestors of each element in the current Selection. It returns a new Selection object with the matched elements.`,
				Value:       selectionParents,
				Returns:     []*luadoc.Param{selection},
			},
			{
				Name:        "parents_until",
				Description: `Gets the ancestors of each element in the current Selection, up to but not including the element matched by the selector. It returns a new Selection object with the matched elements.`,
				Value:       selectionParentsUntil,
				Params:      []*luadoc.Param{selector},
				Returns:     []*luadoc.Param{selection},
			},
			{
				Name:        "slice",
				Description: `Returns a selection containing a subset of the elements in the original selection. It returns a new selection object with the matched elements.`,
				Value:       selectionSlice,
				Params: []*luadoc.Param{
					{
						Name:        "start",
						Description: `The index of the first element to include in the new selection.`,
						Type:        luadoc.Number,
					},
					{
						Name:        "finish",
						Description: `The index of the first element to exclude from the new selection.`,
						Type:        luadoc.Number,
					},
				},
				Returns: []*luadoc.Param{selection},
			},
			{
				Name:        "terminate",
				Description: `Ends the most recent filtering operation in the current chain and returns the set of matched elements to its previous state.`,
				Value:       selectionTerminate,
				Returns:     []*luadoc.Param{selection},
			},
			{
				Name:        "add",
				Description: `Add adds the selector string's matching nodes to those in the current selection and returns a new selection object. The selector string is run in the context of the document of the current selection object.`,
				Value:       selectionAdd,
				Params:      []*luadoc.Param{selector},
				Returns:     []*luadoc.Param{selection},
			},
			{
				Name:        "length",
				Description: `Returns the number of elements in the selection.`,
				Value:       selectionLength,
				Returns: []*luadoc.Param{
					{
						Name:        "n",
						Description: `The number of elements in the selection.`,
						Type:        luadoc.Number,
					},
				},
			},
			{
				Name:        "add_back",
				Description: `Adds the specified Selection object's nodes to those in the current selection and returns a new Selection object.`,
				Value:       selectionAddBack,
				Params:      []*luadoc.Param{selection},
				Returns:     []*luadoc.Param{selection},
			},
			{
				Name:        "add_selection",
				Description: `Adds the specified selection object's nodes to those in the current selection and returns a new selection object.`,
				Value:       selectionAddSelection,
				Params:      []*luadoc.Param{selection},
				Returns:     []*luadoc.Param{selection},
			},
			{
				Name:        "markdown",
				Description: `Converts the selection to Markdown. Can be used to show the contents of a selection in info page`,
				Value:       selectionMarkdown,
				Returns: []*luadoc.Param{
					{
						Name:        "markdown",
						Description: `The Markdown representation of the selection.`,
						Type:        luadoc.String,
					},
				},
			},
		},
	}

	classDocument := &luadoc.Class{
		Name:        documentTypeName,
		Description: `Document represents an HTML document to be manipulated. Unlike jQuery, which is loaded as part of a DOM document, and thus acts upon its containing document, GoQuery doesn't know which HTML document to act upon. So it needs to be told, and that's what the document class is for. It holds the root document node to manipulate, and can make selections on this document.`,
		Methods: []*luadoc.Method{
			{
				Name:        "find",
				Description: `Finds all elements that match the selector string. It returns a new selection object with the matched elements.`,
				Value:       documentFind,
				Params:      []*luadoc.Param{selector},
				Returns:     []*luadoc.Param{selection},
			},
			{
				Name:        "selection",
				Description: `Converts document to a selection object.`,
				Value:       documentSelection,
				Returns:     []*luadoc.Param{selection},
			},
			{
				Name:        "html",
				Description: `Gets the HTML contents of the first element in the set of matched elements. It includes text and comment nodes.`,
				Value:       documentHtml,
				Returns: []*luadoc.Param{
					{
						Name:        "html",
						Description: `The HTML contents of the first element in the set of matched elements.`,
						Type:        luadoc.String,
					},
					{
						Name:        "error",
						Description: `An error if one occurred.`,
						Type:        luadoc.String,
						Optional:    true,
					},
				},
			},
			{
				Name:        "markdown",
				Description: `Converts the document to Markdown. Can be used to show the contents of a document in info page`,
				Value:       documentMarkdown,
				Returns: []*luadoc.Param{
					{
						Name:        "markdown",
						Description: `The Markdown representation of the document.`,
						Type:        luadoc.String,
					},
				},
			},
		},
	}

	return &luadoc.Lib{
		Name:        "html",
		Description: "This library provides functions for parsing HTML and querying it using CSS selectors. It is based on [goquery](https://github.com/PuerkitoBio/goquery).",
		Funcs: []*luadoc.Func{
			{
				Name:        "parse",
				Description: `Parses the given HTML and returns a selection containing the root element.`,
				Value:       parse,
				Params: []*luadoc.Param{
					{
						Name:        "html",
						Description: `The HTML to parse.`,
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "document",
						Description: `The document object.`,
						Type:        documentTypeName,
					},
					{
						Name:        "error",
						Description: `An error if one occurred.`,
						Type:        luadoc.String,
						Optional:    true,
					},
				},
			},
		},
		Classes: []*luadoc.Class{
			classSelection,
			classDocument,
		},
	}
}
