package configsTest

func Run() int {
	CreateJsonFiles()
	TestV2Ray := RunTestV2RayProcesses()
	// RemoveTestsDir()
	return TestV2Ray
}
