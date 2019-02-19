package ast

func replaceVariable(v string, vs map[string]string) string {
	if v, ok := vs[v]; ok {
		return v
	}

	return v
}

func removeVariables(vs map[string]string, ss ...string) map[string]string {
	vvs := make(map[string]string, len(vs))

	for k, v := range vs {
		vvs[k] = v
	}

	for _, s := range ss {
		delete(vvs, s)
	}

	return vvs
}
