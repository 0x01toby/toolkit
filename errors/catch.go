package errors

func Capture() {
	x := recover()
	switch t := x.(type) {
	case nil:
		return
	case error:
		panic(t)
	default:
		//rvalString := fmt.Sprint(t)
		panic(t)
	}

}
