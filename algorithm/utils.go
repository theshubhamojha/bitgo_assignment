package algorithm

// used in DFS for maintaining the queue out of a slice
func Enqueue(queue []string, element string) []string {
	queue = append(queue, element)
	return queue
}

func Dequeue(queue []string) []string {
	return queue[1:]
}
