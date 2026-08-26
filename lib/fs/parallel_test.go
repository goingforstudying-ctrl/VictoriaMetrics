package fs

import (
	"path/filepath"
	"testing"
)

func TestParallelReaderAtOpenerRandomAccessHint(t *testing.T) {
	dir := t.TempDir()
	defaultPath := filepath.Join(dir, "default")
	randomPath := filepath.Join(dir, "random")
	MustWriteSync(defaultPath, []byte("default"))
	MustWriteSync(randomPath, []byte("random"))

	var defaultReader MustReadAtCloser
	var defaultSize uint64
	var randomReader MustReadAtCloser
	var randomSize uint64

	var pro ParallelReaderAtOpener
	pro.Add(defaultPath, &defaultReader, &defaultSize)
	pro.AddRandomAccess(randomPath, &randomReader, &randomSize)
	pro.Run()
	defer defaultReader.MustClose()
	defer randomReader.MustClose()

	if defaultSize != uint64(len("default")) {
		t.Fatalf("unexpected default file size; got %d; want %d", defaultSize, len("default"))
	}
	if randomSize != uint64(len("random")) {
		t.Fatalf("unexpected random-access file size; got %d; want %d", randomSize, len("random"))
	}
	if defaultReader.(*ReaderAt).useRandomReadHint {
		t.Fatalf("unexpected random read hint for Add")
	}
	if !randomReader.(*ReaderAt).useRandomReadHint {
		t.Fatalf("missing random read hint for AddRandomAccess")
	}
}
