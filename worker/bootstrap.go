package worker

// InitWorker is the entry point of init the worker, during which multiple components of
// worker need to be initialized.
func InitWorker(filePath string) error {
	if err := InitConfig(filePath); err != nil {
		return err
	}

	if err := InitRegister(); err != nil {
		return err
	}

	if err := InitLogSink(); err != nil {
		return err
	}

	if err := InitExecutor(); err != nil {
		return err
	}

	// todo init scheduler

	if err := InitJobManager(); err != nil {
		return err
	}

	return nil
}
