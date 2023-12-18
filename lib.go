package ges

/***************************
    @author: tiansheng.ren
    @date: 2022/12/22
    @desc:

***************************/

type doc struct {
	Id       string
	Document interface{}
}

func (d doc) ID() string {
	return d.Id
}
func (d doc) Doc() interface{} {
	return d.Document
}

func (d doc) Item() (string, interface{}) {
	return d.Id, d.Document
}

func NewDoc(id string, d interface{}) Document {
	return doc{Id: id, Document: d}
}

func DocsFromMap(docs map[string]interface{}) []Document {
	res := make([]Document, 0, len(docs))
	for id, d := range docs {
		res = append(res, doc{Id: id, Document: d})
	}
	return res
}

type mapStrAny map[string]interface{}

