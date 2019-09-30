package main

//根据数值和花色排序
func CardSort(cards []int) {
	for i := 0; i < len(cards)-1; i++ {
		for j := 0; j < len(cards)-1-i; j++ {
			if GetCradValue(cards[j]) < GetCradValue(cards[j+1]) {
				cards[j], cards[j+1] = cards[j+1], cards[j]
			} else if GetCradValue(cards[j]) == GetCradValue(cards[j+1]) {
				if GetCardColr(cards[j]) < GetCardColr(cards[j+1]) {
					cards[j], cards[j+1] = cards[j+1], cards[j]
				}
			}
		}
	}
}
