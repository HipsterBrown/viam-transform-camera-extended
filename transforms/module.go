package transforms

import (
	"fmt"
	"go.viam.com/rdk/utils"
	"image"
	"image/color"

	imaging "github.com/disintegration/imaging"
)

type TransformType string

const (
	Resize     TransformType = "resize"
	Brightness TransformType = "brightness"
	Contrast   TransformType = "contrast"
	Blur       TransformType = "blur"
	Sharpen    TransformType = "sharpen"
	Grayscale  TransformType = "grayscale"
	FlipH      TransformType = "fliph"
	FlipV      TransformType = "flipv"
	Rotate     TransformType = "rotate"
)

type Transform struct {
	Type   TransformType      `json:"type"`
	Params utils.AttributeMap `json:"params"`
}

// Specific parameter structs for each transformation type
type ResizeParams struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type BrightnessParams struct {
	Amount float64 `json:"amount"`
}

type ContrastParams struct {
	Amount float64 `json:"amount"`
}

type BlurParams struct {
	Sigma float64 `json:"sigma"`
}

type SharpenParams struct {
	Sigma float64 `json:"sigma"`
}

type RotateParams struct {
	Angle float64 `json:"angle"`
}

func ApplyTransform(img *image.NRGBA, transform Transform) (*image.NRGBA, error) {
	switch transform.Type {
	case Resize:
		width, ok := transform.Params["width"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid width parameter")
		}
		height, ok := transform.Params["height"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid height parameter")
		}
		return imaging.Resize(img, int(width), int(height), imaging.Lanczos), nil

	case Brightness:
		amount, ok := transform.Params["amount"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid amount parameter")
		}
		return imaging.AdjustBrightness(img, amount), nil

	case Contrast:
		amount, ok := transform.Params["amount"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid amount parameter")
		}
		return imaging.AdjustContrast(img, amount), nil

	case Blur:
		sigma, ok := transform.Params["sigma"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid sigma parameter")
		}
		return imaging.Blur(img, sigma), nil

	case Sharpen:
		sigma, ok := transform.Params["sigma"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid sigma parameter")
		}
		return imaging.Sharpen(img, sigma), nil

	case Grayscale:
		return imaging.Grayscale(img), nil

	case FlipH:
		return imaging.FlipH(img), nil

	case FlipV:
		return imaging.FlipV(img), nil

	case Rotate:
		angle, ok := transform.Params["angle"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid angle parameter")
		}
		return imaging.Rotate(img, angle, color.Transparent), nil

	default:
		return nil, fmt.Errorf("unsupported transform type: %s", transform.Type)

	}
}
