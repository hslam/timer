package timer

type Less func(i Score,j Score)bool
type Score int64
type Value timerFunc

var LessInt64= func(i Score,j Score)bool{
	if i<j {return true} else {return false}
}

type Node struct {
	score		Score
	value		Value
	prev		*Node
	next		*Node
}

func (n *Node) Score()  Score {
	return n.score
}

func (n *Node) Value()  Value {
	return n.value
}

func (n *Node) Set(value Value) {
	n.value = value
}

func (n *Node) Prev() *Node {
	return n.prev
}

func (n *Node) Next() *Node {
	return n.next
}

type SortedList struct {
	asc 		bool
	less		Less
	front		*Node
	rear		*Node
	length		int
}

func NewSortedList() (*SortedList) {
	front := &Node{
		value:	nil,
		prev: 	nil,
	}
	rear := &Node{
		value:	nil,
		prev:	nil,
	}
	front.next = rear
	rear.prev=front
	return &SortedList{
		asc:		true,
		less:		LessInt64,
		front:		front,
		rear:		rear,
	}
}

func (l *SortedList) Length() int {
	return l.length
}

//read front data node
func (l *SortedList) Front() *Node {
	if l.length == 0 {
		return nil
	}
	return l.front.next
}

//Insert
func (l *SortedList) Insert(score Score,value Value) bool {
	if value == nil {
		return false
	}
	node := &Node{
		score:score,
		value: value,
	}
	if l.length == 0 {
		node.next=l.rear
		node.prev=l.front
		l.front.next = node
		l.rear.prev=node
		l.length=1
		return true
	}
	cur := l.front.next
	if l.asc{
		for l.less(cur.score,score){
			cur = cur.next
			if cur == l.rear{
				break
			}
		}
	}else {
		for l.less(score,cur.score)&& cur != l.rear {
			cur = cur.next
		}
	}
	prev:=cur.prev
	node.next=cur
	node.prev=prev
	prev.next=node
	cur.prev=node
	l.length++
	return true
}

//read and remove
func (l *SortedList) Top() *Node{
	if l.length == 0 {
		return nil
	}
	result := l.front.next
	l.front.next = result.next
	result.next.prev=l.front
	result.next = nil
	result.prev = nil
	l.length--
	return result
}