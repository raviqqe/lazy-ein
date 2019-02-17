package build

// Build builds an executable file or module.
func Build(f, runtimeDir, rootDir, cacheDir string) error {
	return newBuilder(runtimeDir, rootDir, cacheDir).Build(f)
}
