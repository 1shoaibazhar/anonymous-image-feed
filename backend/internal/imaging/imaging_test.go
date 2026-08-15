package imaging

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"strings"
	"testing"
)

func TestNormalizeResizesToFixedSize(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 400, 200))
	for x := 0; x < 400; x++ {
		for y := 0; y < 200; y++ {
			src.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, src); err != nil {
		t.Fatalf("failed to encode test png: %v", err)
	}

	data, err := Normalize(&buf)
	if err != nil {
		t.Fatalf("Normalize returned error: %v", err)
	}

	out, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode normalized output: %v", err)
	}

	bounds := out.Bounds()
	if bounds.Dx() != targetWidth || bounds.Dy() != targetHeight {
		t.Errorf("got size %dx%d, want %dx%d", bounds.Dx(), bounds.Dy(), targetWidth, targetHeight)
	}
}

func TestNormalizeInvalidImageReturnsError(t *testing.T) {
	_, err := Normalize(strings.NewReader("not an image"))
	if err == nil {
		t.Error("expected an error for invalid image data, got nil")
	}
}
