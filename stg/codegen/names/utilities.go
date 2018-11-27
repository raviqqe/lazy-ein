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

// ToTag converts a constructor name into its tag name.
func ToTag(s string) string {
	return s + ".tag"
}

// ToUnionify converts a constructor name into a name of a function which
// converts structs into unions.
func ToUnionify(s string) string {
	return s + ".unionify"
}

// ToStructify converts a constructor name into a name of a function which
// converts unions into structs.
func ToStructify(s string) string {
	return s + ".structify"
}
