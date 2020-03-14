package Tree

func (node Node) IsNullable() bool {
	v := node.GetValue()
	if v == "*" {
		return true
	} else if v == "?" {
		return true
	} else if v == "" {
		return true
	}
	return false
}
