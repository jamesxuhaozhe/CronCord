package worker

func InitWorker(filePath string) error {
	if err := InitConfig(filePath); err != nil {
		return err
	}
	return nil
}
