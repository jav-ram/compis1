package Tree

func (node Node) IsNullable() bool {
	v := node.GetValue()
	if v == "*" {
		return true
	} else if v == "?" {
		return true
	} else if v == "`" { //TODO: epsilon
		return true
	}
	return false
}
