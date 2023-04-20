/*
	Copyright 2020 The pdfcpu Authors.

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

package pdfcpu

import (
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/pkg/errors"
)

var (
	errNoBookmarks        = errors.New("pdfcpu: no bookmarks available")
	errCorruptedBookmarks = errors.New("pdfcpu: corrupt bookmark")
	errExistingBookmarks  = errors.New("pdfcpu: existing bookmarks")
)

// Bookmark represents an outline item tree.
type Bookmark struct {
	Title    string
	PageFrom int
	PageThru int // for extraction only; >= pageFrom and reaches until before pageFrom of the next bookmark.
	Bold     bool
	Italic   bool
	Color    *color.SimpleColor
	Children []Bookmark
	Parent   *Bookmark
}

// Style returns an int corresponding to the bookmark style.
func (bm Bookmark) Style() int {
	var i int
	if bm.Bold {
		i += 2
	}
	if bm.Italic {
		i += 1
	}
	return i
}

func positionToOutlineTreeLevel1(ctx *model.Context) (types.Dict, *types.IndirectRef, error) {
	// Load Dests nametree.
	if err := ctx.LocateNameTree("Dests", false); err != nil {
		return nil, nil, err
	}

	ir, err := ctx.Outlines()
	if err != nil {
		return nil, nil, err
	}
	if ir == nil {
		return nil, nil, errNoBookmarks
	}

	d, err := ctx.DereferenceDict(*ir)
	if err != nil {
		return nil, nil, err
	}
	if d == nil {
		return nil, nil, errNoBookmarks
	}

	first := d.IndirectRefEntry("First")
	last := d.IndirectRefEntry("Last")

	// We consider Bookmarks at level 1 or 2 only.
	for *first == *last {
		if d, err = ctx.DereferenceDict(*first); err != nil {
			return nil, nil, err
		}
		first = d.IndirectRefEntry("First")
		last = d.IndirectRefEntry("Last")
	}

	return d, first, nil
}

func outlineItemTitle(s string) string {
	var sb strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b >= 32 {
			sb.WriteByte(b)
		}
	}
	return sb.String()
}

// PageObjFromDestinationArray return an IndirectRef for this destinations page object.
func PageObjFromDestinationArray(ctx *model.Context, dest types.Object) (*types.IndirectRef, error) {
	var (
		err error
		ir  types.IndirectRef
		arr types.Array
	)
	switch dest := dest.(type) {
	case types.Name:
		arr, err = ctx.DereferenceDestArray(dest.Value())
		if err == nil {
			ir = arr[0].(types.IndirectRef)
		}
	case types.StringLiteral:
		arr, err = ctx.DereferenceDestArray(dest.Value())
		if err == nil {
			ir = arr[0].(types.IndirectRef)
		}
	case types.HexLiteral:
		arr, err = ctx.DereferenceDestArray(dest.Value())
		if err == nil {
			ir = arr[0].(types.IndirectRef)
		}
	case types.Array:
		if dest[0] != nil {
			ir = dest[0].(types.IndirectRef)
		}
		// else skipping bookmarks that don't point to anything.
	}
	return &ir, err
}

// BookmarksForOutlineItem returns the bookmarks tree for an outline item.
func BookmarksForOutlineItem(ctx *model.Context, item *types.IndirectRef, parent *Bookmark) ([]Bookmark, error) {
	bms := []Bookmark{}

	var (
		d   types.Dict
		err error
	)

	// Process outline items.
	for ir := item; ir != nil; ir = d.IndirectRefEntry("Next") {

		if d, err = ctx.DereferenceDict(*ir); err != nil {
			return nil, err
		}

		s, _ := model.Text(d["Title"])
		title := outlineItemTitle(s)

		// Retrieve page number out of a destination via "Dest" or "Goto Action".
		dest, destFound := d["Dest"]
		if !destFound {
			act, actFound := d["A"]
			if !actFound {
				continue
			}
			act, _ = ctx.Dereference(act)
			actType := act.(types.Dict)["S"]
			if actType.String() != "GoTo" {
				continue
			}
			dest = act.(types.Dict)["D"]
		}

		dest, err := ctx.Dereference(dest)
		if err != nil {
			return nil, err
		}

		ir, err := PageObjFromDestinationArray(ctx, dest)
		if err != nil {
			return nil, err
		}
		if ir == nil {
			continue
		}

		pageFrom, err := ctx.PageNumber(ir.ObjectNumber.Value())
		if err != nil {
			return nil, err
		}

		if len(bms) > 0 {
			if pageFrom > bms[len(bms)-1].PageFrom {
				bms[len(bms)-1].PageThru = pageFrom - 1
			} else {
				bms[len(bms)-1].PageThru = bms[len(bms)-1].PageFrom
			}
		}

		newBookmark := Bookmark{
			Title:    title,
			PageFrom: pageFrom,
			Parent:   parent,
		}

		first := d["First"]
		if first != nil {
			indRef := first.(types.IndirectRef)
			children, _ := BookmarksForOutlineItem(ctx, &indRef, &newBookmark)
			newBookmark.Children = children
		}

		bms = append(bms, newBookmark)
	}

	return bms, nil
}

// BookmarksForOutline returns all ctx bookmark information recursively.
func BookmarksForOutline(ctx *model.Context) ([]Bookmark, error) {
	_, first, err := positionToOutlineTreeLevel1(ctx)
	if err != nil {
		return nil, err
	}

	return BookmarksForOutlineItem(ctx, first, nil)
}

func bmDict(bm Bookmark, pageIndRef, parent types.IndirectRef, s string) types.Dict {

	d := types.Dict(map[string]types.Object{
		"Dest":   types.Array{pageIndRef, types.Name("Fit")},
		"Title":  types.StringLiteral(s),
		"Parent": parent},
	)

	if bm.Color != nil {
		d["C"] = types.Array{types.Float(bm.Color.R), types.Float(bm.Color.G), types.Float(bm.Color.B)}
	}

	if style := bm.Style(); style > 0 {
		d["F"] = types.Integer(style)
	}

	return d
}

func createOutlineItemDict(ctx *model.Context, bms []Bookmark, parent *types.IndirectRef, parentPageNr *int) (*types.IndirectRef, *types.IndirectRef, int, error) {
	var (
		first  *types.IndirectRef
		irPrev *types.IndirectRef
		dPrev  types.Dict
		count  int
	)

	for i, bm := range bms {

		if i == 0 && parentPageNr != nil && bm.PageFrom < *parentPageNr {
			return nil, nil, 0, errCorruptedBookmarks
		}

		if i > 0 && bm.PageFrom < bms[i-1].PageFrom {
			return nil, nil, 0, errCorruptedBookmarks
		}

		_, pageIndRef, _, err := ctx.PageDict(bm.PageFrom, false)
		if err != nil {
			return nil, nil, 0, err
		}

		s, err := types.Escape(types.EncodeUTF16String(bm.Title))
		if err != nil {
			return nil, nil, 0, err
		}

		d := bmDict(bm, *pageIndRef, *parent, *s)
		ir, err := ctx.IndRefForNewObject(d)
		if err != nil {
			return nil, nil, 0, err
		}

		if first == nil {
			first = ir
		}

		if bm.Children != nil {
			first, last, c, err := createOutlineItemDict(ctx, bm.Children, ir, &bm.PageFrom)
			if err != nil {
				return nil, nil, 0, err
			}
			d["First"] = *first
			d["Last"] = *last

			d["Count"] = types.Integer(c + 1)
			count += c + 1
		} else {
			count++
		}

		if irPrev != nil {
			d["Prev"] = *irPrev
			dPrev["Next"] = *ir
		}

		dPrev = d
		irPrev = ir

	}

	return first, irPrev, count, nil
}

// AddBookmarks adds bms to ctx.
func AddBookmarks(ctx *model.Context, bms []Bookmark) error {

	rootDict, err := ctx.Catalog()
	if err != nil {
		return err
	}

	if _, ok := rootDict.Find("Outlines"); ok {
		return errExistingBookmarks
	}

	outlinesDict := types.Dict(map[string]types.Object{"Type": types.Name("Outlines")})
	outlinesir, err := ctx.IndRefForNewObject(outlinesDict)
	if err != nil {
		return err
	}

	first, last, count, err := createOutlineItemDict(ctx, bms, outlinesir, nil)
	if err != nil {
		return err
	}

	outlinesDict["First"] = *first
	outlinesDict["Last"] = *last
	outlinesDict["Count"] = types.Integer(count)

	rootDict["Outlines"] = *outlinesir

	return nil
}
