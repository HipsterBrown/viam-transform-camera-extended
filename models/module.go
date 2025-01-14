package models

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"

	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/components/camera/rtppassthrough"
	"go.viam.com/rdk/gostream"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/pointcloud"
	"go.viam.com/rdk/resource"
	"go.viam.com/utils/rpc"

	imaging "github.com/disintegration/imaging"

	"transform-camera-extended/transforms"
)

var (
	TransformExtended = resource.NewModel("hipsterbrown", "camera", "transform-extended")
	errUnimplemented  = errors.New("unimplemented")
)

func init() {
	resource.RegisterComponent(camera.API, TransformExtended,
		resource.Registration[camera.Camera, *Config]{
			Constructor: newTransformCameraExtendedTransformExtended,
		},
	)
}

type Config struct {
	Source   string                 `json:"source"`
	Pipeline []transforms.Transform `json:"pipeline"`
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns implicit dependencies based on the config.
// The path is the JSON path in your robot's config (not the `Config` struct) to the
// resource being validated; e.g. "components.0".
func (cfg *Config) Validate(path string) ([]string, error) {
	var deps []string
	if len(cfg.Source) == 0 {
		return nil, resource.NewConfigValidationFieldRequiredError(path, "source")
	}

	deps = append(deps, cfg.Source)
	return deps, nil
}

type transformCameraExtendedTransformExtended struct {
	name resource.Name

	logger logging.Logger
	cfg    *Config

	cancelCtx  context.Context
	cancelFunc func()

	/* Uncomment this if your model does not need to reconfigure. */
	resource.TriviallyReconfigurable

	// Uncomment this if the model does not have any goroutines that
	// need to be shut down while closing.
	resource.TriviallyCloseable

	source camera.Camera
}

func newTransformCameraExtendedTransformExtended(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (camera.Camera, error) {

	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	s := &transformCameraExtendedTransformExtended{
		name:       rawConf.ResourceName(),
		logger:     logger,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
	}

	if err := s.Reconfigure(ctx, deps, rawConf); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *transformCameraExtendedTransformExtended) Name() resource.Name {
	return s.name
}

func (s *transformCameraExtendedTransformExtended) Reconfigure(ctx context.Context, deps resource.Dependencies, rawConf resource.Config) error {
	// Put reconfigure code here
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return err
	}
	s.cfg = conf

	source, err := camera.FromDependencies(deps, conf.Source)
	if err != nil {
		return fmt.Errorf("no source camera for transform pipeline  (%s): %w", conf.Source, err)
	}
	s.source = source
	return nil
}

func (s *transformCameraExtendedTransformExtended) NewClientFromConn(ctx context.Context, conn rpc.ClientConn, remoteName string, name resource.Name, logger logging.Logger) (camera.Camera, error) {
	return nil, errUnimplemented
}

func (s *transformCameraExtendedTransformExtended) Stream(ctx context.Context, errHandlers ...gostream.ErrorHandler) (gostream.VideoStream, error) {
	return s.source.Stream(ctx, errHandlers...)
}

func (s *transformCameraExtendedTransformExtended) Image(ctx context.Context, mimeType string, extra map[string]interface{}) ([]byte, camera.ImageMetadata, error) {
	srcImage, metadata, err := s.source.Image(ctx, mimeType, extra)
	if err != nil {
		return nil, camera.ImageMetadata{}, err
	}
	s.logger.Info(metadata)
	reader := bytes.NewReader(srcImage)
	decodedImage, err := jpeg.Decode(reader)
	if err != nil {
		return nil, camera.ImageMetadata{}, err
	}

	adjustedImage, err := s.processImage(decodedImage)
	if err != nil {
		return nil, camera.ImageMetadata{}, err
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, adjustedImage, &jpeg.Options{Quality: 90})
	if err != nil {
		return nil, camera.ImageMetadata{}, err
	}
	return buf.Bytes(), metadata, nil

}

func (s *transformCameraExtendedTransformExtended) Images(ctx context.Context) ([]camera.NamedImage, resource.ResponseMetadata, error) {
	return s.source.Images(ctx)
}

func (s *transformCameraExtendedTransformExtended) NextPointCloud(ctx context.Context) (pointcloud.PointCloud, error) {
	return nil, errUnimplemented
}

func (s *transformCameraExtendedTransformExtended) Properties(ctx context.Context) (camera.Properties, error) {
	return s.source.Properties(ctx)
}

func (s *transformCameraExtendedTransformExtended) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, errUnimplemented
}

func (s *transformCameraExtendedTransformExtended) SubscribeRTP(ctx context.Context, bufferSize int, packetsCB rtppassthrough.PacketCallback) (rtppassthrough.Subscription, error) {
	panic("not implemented")
}

func (s *transformCameraExtendedTransformExtended) Unsubscribe(ctx context.Context, id rtppassthrough.SubscriptionID) error {
	panic("not implemented")
}

func (s *transformCameraExtendedTransformExtended) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}

func (s transformCameraExtendedTransformExtended) processImage(img image.Image) (*image.NRGBA, error) {
	result := imaging.Clone(img)

	for _, transform := range s.cfg.Pipeline {
		var err error
		result, err = transforms.ApplyTransform(result, transform)
		if err != nil {
			return nil, fmt.Errorf("failed to apply transform %s: %v", transform.Type, err)
		}
	}
	return result, nil
}
