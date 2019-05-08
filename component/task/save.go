package task

import (
	"fmt"

	"github.com/safeie/spider/component/url"
)

// BeforeSaveFunc 存储前置方法
type BeforeSaveFunc func(rule *Rule, uri *url.URI, val map[string]interface{}) map[string]interface{}

// AfterSaveFunc 存储后置方法
type AfterSaveFunc func(rule *Rule, uri *url.URI)

// SaveFunc 存储方法
type SaveFunc func(taskID, pk string, val map[string]interface{}) error

// defaultSaveFunc 默认存储费方法，将数据打印出来
func defaultSaveFunc(taskID, pk string, val map[string]interface{}) error {
	fmt.Println("----- savedata -----")
	fmt.Printf("TaskID: %s\n", taskID)
	fmt.Printf("DataPK: %s\n", pk)
	fmt.Printf("Data Value:\n%v\n", val)
	fmt.Println("")
	return nil
}
