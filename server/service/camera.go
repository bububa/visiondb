package service

import (
	"image"
	"time"

	"github.com/bububa/camera"
	"go.uber.org/atomic"

	"github.com/bububa/visiondb/logger"
	"github.com/bububa/visiondb/server/service/cam"
	"github.com/bububa/visiondb/utils"
)

type cameraService struct {
	cam                   *camera.Camera
	config                camera.Options
	converter             utils.ImageConverter
	imagesClassifier      utils.ImagesClassifier
	frames                []image.Image
	delay                 int
	startCh               chan chan error
	stopCh                chan chan error
	readFrameCh           chan struct{}
	frameCh               chan image.Image
	closeCh               chan struct{}
	exitCh                chan struct{}
	setConverterCh        chan utils.ImageConverter
	setImagesClassifierCh chan utils.ImagesClassifier
	isStarted             *atomic.Bool
}

var cameraInstance *cameraService

func Camera() *cameraService {
	return cameraInstance
}

func closeCamera() {
	if cameraInstance == nil {
		return
	}
	cameraInstance.Close()
}

func initCamera(config camera.Options) (err error) {
	cameraInstance = &cameraService{
		config:                config,
		delay:                 config.Delay,
		isStarted:             atomic.NewBool(false),
		startCh:               make(chan chan error, 1),
		stopCh:                make(chan chan error, 1),
		readFrameCh:           make(chan struct{}, 1),
		closeCh:               make(chan struct{}, 1),
		exitCh:                make(chan struct{}, 1),
		setConverterCh:        make(chan utils.ImageConverter, 1),
		setImagesClassifierCh: make(chan utils.ImagesClassifier, 1),
	}
	go cameraInstance.startWorker()
	return err
}

func (s *cameraService) startWorker() {
	for {
		select {
		case syncCh := <-s.startCh:
			syncCh <- s.doStart()
		case <-s.readFrameCh:
			s.doRead()
		case syncCh := <-s.stopCh:
			logger.Warn().Msg("stopping")
			syncCh <- s.doStop()
			logger.Warn().Msg("stopped")
		case <-s.closeCh:
			logger.Warn().Msg("closing")
			s.doClose()
			logger.Warn().Msg("closed")
		case converter := <-s.setConverterCh:
			s.converter = converter
		case classifier := <-s.setImagesClassifierCh:
			s.imagesClassifier = classifier
		case <-s.exitCh:
			return
		}
	}
}

func (s *cameraService) doStart() (err error) {
	if s.cam != nil {
		if err := s.doStop(); err != nil {
			logger.Error().Err(err).Send()
			return err
		}
	}
	if s.cam, err = cam.New(s.config); err != nil {
		logger.Error().Err(err).Send()
		return err
	}
	if err := s.cam.Start(); err != nil {
		logger.Error().Err(err).Send()
		return err
	}
	s.frames = []image.Image{}
	s.frameCh = make(chan image.Image, 120)
	s.readFrameCh <- struct{}{}
	s.isStarted.Store(true)
	return nil
}

func (s *cameraService) doStop() error {
	s.isStarted.CAS(true, false)
	if s.imagesClassifier != nil {
		s.imagesClassifier(s.frames)
	}
	if s.cam != nil {
		if err := s.cam.Close(); err != nil {
			logger.Error().Err(err).Send()
			return err
		}
	}
	for len(s.readFrameCh) > 0 {
		<-s.readFrameCh
	}
	if s.frameCh != nil {
		close(s.frameCh)
		s.frameCh = nil
	}
	s.converter = nil
	s.imagesClassifier = nil
	s.frames = []image.Image{}
	s.delay = s.config.Delay
	s.cam = nil
	return nil
}

func (s *cameraService) doRead() {
	defer func() {
		if !s.isStarted.Load() {
			return
		}
		time.Sleep(time.Duration(s.Delay()) * time.Millisecond)
		s.readFrameCh <- struct{}{}
	}()
	if !s.isStarted.Load() {
		return
	}
	if in, err := s.cam.Read(); err != nil {
		s.frameCh <- in
		return
	} else {
		if s.imagesClassifier != nil {
			s.frames = append(s.frames, in)
		}
		if s.converter != nil {
			if out, err := s.converter(in); err == nil {
				s.frameCh <- out
				return
			}
		} else {
			s.frameCh <- in
		}
	}
}

// Start start camera
func (s *cameraService) Start() error {
	syncCh := make(chan error, 1)
	s.startCh <- syncCh
	return <-syncCh
}

// Stop stop camera
func (s *cameraService) Stop() error {
	syncCh := make(chan error, 1)
	s.stopCh <- syncCh
	return <-syncCh
}

func (s *cameraService) doClose() {
	s.doStop()
	close(s.startCh)
	close(s.stopCh)
	close(s.readFrameCh)
	s.exitCh <- struct{}{}
}

// Close close camera
func (s *cameraService) Close() {
	s.closeCh <- struct{}{}
	<-s.exitCh
}

// Delay returns delay between frames
func (s *cameraService) Delay() int {
	return s.delay
}

// SetDelay set delay between frames
func (s *cameraService) SetDelay(delay int) {
	s.delay = delay
}

// SetConverter set image converter
func (s *cameraService) SetConverter(fn utils.ImageConverter) {
	s.setConverterCh <- fn
}

// SetImagesClassifier set images classifier
func (s *cameraService) SetImagesClassifier(fn utils.ImagesClassifier) {
	s.setImagesClassifierCh <- fn
}

// Read read image from camera
func (s *cameraService) Read() <-chan image.Image {
	return s.frameCh
}
