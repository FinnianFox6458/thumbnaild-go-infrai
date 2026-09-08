package thumb

import "fmt"

type Request struct {
	Image  string `json:"image"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Fit    string `json:"fit"`
	Format string `json:"format"`
}

func Validate(r Request) error {
	if r.Image == "" {
		return fmt.Errorf("image is required")
	}
	if r.Width <= 0 || r.Height <= 0 {
		return fmt.Errorf("width and height must be positive")
	}
	if r.Fit == "" {
		return fmt.Errorf("fit is required")
	}
	return nil
}
