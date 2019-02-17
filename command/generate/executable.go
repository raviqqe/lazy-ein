package generate

// Executable builds an executable file.
func Executable(f, runtimePath, cacheDir string) error {
	return newBuilder(runtimePath, cacheDir).BuildExecutable(f)
}
