package main

import (
	"example.com/thumbnaild/internal/thumb"
	"testing"
)

func TestImageProcessPayloadUsesContractFields(t *testing.T) {
	payload := imageProcessPayload(thumb.Request{
		Image: "img_123", Width: 320, Height: 180, Fit: "cover", Format: "webp",
	})

	for _, key := range []string{"image", "ops", "format", "store"} {
		if _, ok := payload[key]; !ok {
			t.Fatalf("missing image.process field %q", key)
		}
	}
	for _, key := range []string{"width", "height", "fit", "enlarge"} {
		if _, ok := payload[key]; ok {
			t.Fatalf("invalid top-level image.process field %q", key)
		}
	}
}
