package sequence

func fRecover() {
	recover()
}
func safeSend(c chan<- interface{}, value interface{}) {
	defer fRecover()
	c <- value
}

func safeClose(c chan interface{}) {
	defer fRecover()
	close(c)
}
