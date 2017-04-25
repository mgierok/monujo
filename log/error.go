package log

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
