package directory_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"source/directory"
)

func TestList(t *testing.T) {
	// Creating a library with a non-existent directory should return an error.
	_, err := directory.New("./inexistent-testdir", []string{".txt"})
	require.Error(t, err)

	// Creating a library with a file should return an error.
	_, err = directory.New("./testdir/0.txt", []string{".txt"})
	require.Error(t, err)

	// Creating a library with a valid directory should work.
	l, err := directory.New("./testdir", []string{".txt"})
	require.NoError(t, err)

	// Listing the library should return all the files in the directory.
	media, err := l.List()
	require.NoError(t, err)

	require.Len(t, media, 4)
	require.Equal(t, "./testdir/0.txt", media[0].URI())
	require.Equal(t, "./testdir/a/1.txt", media[1].URI())
	require.Equal(t, "./testdir/b/c/2.txt", media[2].URI())
	require.Equal(t, "./testdir/b/d/e/3.txt", media[3].URI())
}
