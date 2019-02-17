package build

// Executable builds an executable file.
func Executable(f, runtimeDir, rootDir, cacheDir string) error {
	return newBuilder(runtimeDir, rootDir, cacheDir).BuildExecutable(f)
}
