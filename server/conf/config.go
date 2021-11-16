package conf

type Config struct {
	AppName        string `required:"true"`
	Port           int    `required:"true"`
	Domain         string `required:"true"`
	BaseUrl        string `required:"true"`
	ConsoleBaseUrl string `required:"true"`
	Template       string `required:"true"`
	FaceID         FaceID `required:"true"`
	LogPath        string
	Debug          bool
}

type FaceID struct {
	DetectorModelPath   string
	RecognizerModelPath string
	DatabasePath        string
}
