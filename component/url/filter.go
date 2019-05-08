package url

import "strings"

// FilterGroupDefault 默认过滤组，包含多个过滤器，脚本，样式表，注释，空白
func FilterGroupDefault(f *Field) {
	FilterRemoveScript(f)
	FilterRemoveStyle(f)
	FilterRemoveNote(f)
	FilterRemoveBlank(f)
}

// FilterRemoveBlank 过滤器，过滤空行和文本两头的空白
func FilterRemoveBlank(f *Field) {
	v := f.String()
	if v == "" {
		return
	}
	v = strings.Replace(v, "\r\n", "\n", -1)
	v = strings.Replace(v, "\n\n", "", -1)
	v = strings.Replace(v, "\n\n", "", -1)
	v = strings.TrimSpace(v)
	f.SetValue(v)
}

// FilterRemoveScript 过滤器，删除内容中的 script脚本
func FilterRemoveScript(f *Field) {
	f.Remove("<script(*)</script>")
}

// FilterRemoveNote 过滤器，删除注释
func FilterRemoveNote(f *Field) {
	f.Remove("<!--(*)-->")
}

// FilterRemoveStyle 过滤器，删除style样式表
func FilterRemoveStyle(f *Field) {
	f.Remove("<style(*)</style>")
}

// FilterRemoveImgage 过滤器，删除所有图片
func FilterRemoveImgage(f *Field) {
	f.RemoveHTMLTags("img")
}

// FilterRemoveA 过滤器，删除所有A链接
func FilterRemoveA(f *Field) {
	f.RemoveHTMLTags("a")
}

// FilterRemoveXMLCDATA 过滤器，删除XML CDATA标记
func FilterRemoveXMLCDATA(f *Field) {
	v := f.String()
	if v == "" {
		return
	}
	v = strings.Replace(v, "<![CDATA[", "", 1)
	v = strings.Replace(v, "]]>", "", 1)
	v = strings.TrimSpace(v)
	f.SetValue(v)
}
