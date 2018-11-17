package names

// ToEntry converts a function name into its entry name.
func ToEntry(s string) string {
	return s + ".entry"
}

// ToUpdatedEntry converts a function name into its updated entry name.
// TODO: Change this function name to "normal form entry".
func ToUpdatedEntry(s string) string {
	return s + ".updated-entry"
}
