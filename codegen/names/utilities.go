package names

// ToEntry converts a function name into its entry name.
func ToEntry(s string) string {
	return s + ".entry"
}

// ToUpdatedEntry converts a function name into its updated entry name.
func ToUpdatedEntry(s string) string {
	return s + ".updated-entry"
}
