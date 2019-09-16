package util

import (
	"sort"
)

//结构体排序，必须重写数组Len() Swap() Less()函数
type body_wrapper struct {
	Bodys []interface{}
	by    func(p, q *interface{}) bool //内部Less()函数会用到
}
type SortBodyBy func(p, q *interface{}) bool //定义一个函数类型

//数组长度Len()
func (acw body_wrapper) Len() int {
	return len(acw.Bodys)
}

//元素交换
func (acw body_wrapper) Swap(i, j int) {
	acw.Bodys[i], acw.Bodys[j] = acw.Bodys[j], acw.Bodys[i]
}

//比较函数，使用外部传入的by比较函数
func (acw body_wrapper) Less(i, j int) bool {
	return acw.by(&acw.Bodys[i], &acw.Bodys[j])
}

// 自定义排序字段
func SortBody(bodys []interface{}, by SortBodyBy) {
	sort.Sort(body_wrapper{bodys, by})
}
