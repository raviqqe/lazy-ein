package names

// ToEntry converts a function name into its entry name.
func ToEntry(s string) string {
	return s + ".entry"
}

// ToNormalFormEntry converts a function name into its normal form entry name.
func ToNormalFormEntry(s string) string {
	return s + ".normal-form.entry"
}

// ToBlackHoleEntry converts a function name into its black hole entry name.
func ToBlackHoleEntry(s string) string {
	return s + ".black-hole.entry"
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
