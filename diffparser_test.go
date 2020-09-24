// Copyright (c) 2015 Jesse Meek <https://github.com/waigani>
// This program is Free Software see LICENSE file for details.

package diffparser

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO(waigani) tests are missing more creative names (spaces, special
// chars), and diffed files that are not in the current directory.

func setup(t *testing.T) *Diff {
	byt, err := ioutil.ReadFile("example.diff")
	require.NoError(t, err)

	diff, err := Parse(string(byt))
	require.NoError(t, err)
	require.Equal(t, len(diff.Files), 6)

	return diff
}
func TestFileModeAndNaming(t *testing.T) {
	diff := setup(t)
	for i, expected := range []struct {
		mode     FileMode
		origName string
		newName  string
	}{
		{
			mode:     Modified,
			origName: "file1",
			newName:  "file1",
		},
		{
			mode:     Deleted,
			origName: "file2",
			newName:  "",
		},
		{
			mode:     Deleted,
			origName: "file3",
			newName:  "",
		},
		{
			mode:     New,
			origName: "",
			newName:  "file4",
		},
		{
			mode:     New,
			origName: "",
			newName:  "newname",
		},
		{
			mode:     Deleted,
			origName: "symlink",
			newName:  "",
		},
	} {
		file := diff.Files[i]
		t.Logf("testing file: %v", file)
		require.Equal(t, expected.mode, file.Mode)
		require.Equal(t, expected.origName, file.OrigName)
		require.Equal(t, expected.newName, file.NewName)
	}
}

func TestHunk(t *testing.T) {
	diff := setup(t)
	expectedOrigLines := []DiffLine{
		{
			Mode:     Unchanged,
			Number:   1,
			Content:  "some",
			Position: 2,
		}, {
			Mode:     Unchanged,
			Number:   2,
			Content:  "lines",
			Position: 3,
		}, {
			Mode:     Removed,
			Number:   3,
			Content:  "in",
			Position: 4,
		}, {
			Mode:     Unchanged,
			Number:   4,
			Content:  "file1",
			Position: 5,
		},
	}

	expectedNewLines := []DiffLine{
		{
			Mode:     Added,
			Number:   1,
			Content:  "add a line",
			Position: 1,
		}, {
			Mode:     Unchanged,
			Number:   2,
			Content:  "some",
			Position: 2,
		}, {
			Mode:     Unchanged,
			Number:   3,
			Content:  "lines",
			Position: 3,
		}, {
			Mode:     Unchanged,
			Number:   4,
			Content:  "file1",
			Position: 5,
		},
	}

	file := diff.Files[0]
	origRange := file.Chunks[0].OrigRange
	newRange := file.Chunks[0].NewRange

	require.Equal(t, 1, origRange.Start)
	require.Equal(t, 4, origRange.Length)
	require.Equal(t, 1, newRange.Start)
	require.Equal(t, 4, newRange.Length)

	for i, line := range expectedOrigLines {
		require.Equal(t, line, *origRange.Lines[i])
	}
	for i, line := range expectedNewLines {
		require.Equal(t, line, *newRange.Lines[i])
	}
}
