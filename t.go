"*+?.|"

if c == "*" {
	l := ev.getTree([]interface{}{input[i-1]})
	l.SetParent(n)
	n.AddLeftChild(l)
	input[i-1] = n
	input = remove(input, i)
} else if c == "." {
	l := ev.getTree([]interface{}{input[i-1]})
	l.SetParent(n)
	r := ev.getTree([]interface{}{input[i+1]})
	r.SetParent(n)
	n.AddChilds(l, r)
	input[i-1] = n
	input = remove(input, i+1)
	input = remove(input, i)
} else if c == "|" {
	l := ev.getTree([]interface{}{input[i-1]})
	l.SetParent(n)
	r := ev.getTree([]interface{}{input[i+1]})
	r.SetParent(n)
	n.AddChilds(l, r)
	fmt.Printf("input%v\n", input)
	input[i-1] = n
	fmt.Printf("input%v\n", input)
	input = remove(input, i+1)
	input = remove(input, i)
	fmt.Printf("input%v\n", input)
}
i--