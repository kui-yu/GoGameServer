package main

func VecDelMulti(d []byte, cs []byte) ([]byte, bool) {
	for _, v := range cs {
		isDel := false
		for i := len(d) - 1; i >= 0; i-- {
			if d[i] == v {
				d = append(d[:i], d[i+1:]...)
				isDel = true
				break
			}
		}
		//
		if !isDel {
			return d, isDel
		}
	}
	//
	return d, true
}
