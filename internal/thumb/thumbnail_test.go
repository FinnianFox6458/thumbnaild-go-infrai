package thumb

import "testing"

func TestValidateThumbnailRequest(t *testing.T) {
	tests := []struct {
		name string
		req  Request
		want bool
	}{{"complete", Request{Image: "img_123", Width: 320, Height: 180, Fit: "cover", Format: "webp"}, true}, {"missing image", Request{Width: 320, Height: 180, Fit: "cover"}, false}, {"zero size", Request{Image: "img_123", Width: 0, Height: 180, Fit: "cover"}, false}}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if (Validate(tc.req) == nil) != tc.want {
				t.Fatal("validation result mismatch")
			}
		})
	}
}
