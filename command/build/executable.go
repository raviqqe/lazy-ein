package build

// Executable builds an executable file.
func Executable(f, runtimePath, moduleRootPath, cacheDir string) error {
	return newBuilder(runtimePath, moduleRootPath, cacheDir).BuildExecutable(f)
}
