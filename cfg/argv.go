package cfg

func getFirstArg(argv []string) (arg string, err error) {
	if len(argv) <= 1 {
		err = ErrEnoughArgs
		return
	}
	arg = argv[1]
	return
}
