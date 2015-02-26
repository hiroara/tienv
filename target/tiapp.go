package target

import (
	"errors"
	"github.com/moovweb/gokogiri/xml"
	"path/filepath"
)

type Tiapp struct {
	xmlFile   string
	directory string
	path      string
	Document  xml.Document
}

func GetTiappWithRestore(directory, backupExt string) (t *Tiapp, err error) {
	t = &Tiapp{xmlFile: "tiapp.xml", directory: directory}

	Restore(t, backupExt)

	t.Document, err = t.readDocument()
	return
}

func GetTiapp(directory string) (t *Tiapp, err error) {
	t = &Tiapp{xmlFile: "tiapp.xml", directory: directory}
	t.Document, err = t.readDocument()
	return
}

func (t *Tiapp) GetFilePath() (path string, err error) {
	if t.path != "" {
		path = t.path
		return
	}
	path, err = filepath.Abs(filepath.Join(t.directory, t.xmlFile))
	t.path = path
	return
}

func (t *Tiapp) Replace(xpath, value string) (err error) {
	node, err := t.searchWithXPath(xpath)
	if err != nil {
		return
	}

	err = node.SetContent(value)
	return
}

func (t *Tiapp) ReplaceWithConf(path string) (err error) {
	json, err := readJson(path)
	replaces := json.Get("replaces")
	for _, elem := range replaces.MustArray() {
		replace := elem.(map[string]interface{})
		xpath := replace["xpath"].(string)
		value := replace["value"].(string)
		err = t.Replace(xpath, value)
		if err != nil {
			return
		}
	}
	return
}

func (t *Tiapp) Append(xpath, content string) (err error) {
	node, err := t.searchWithXPath(xpath)
	if err != nil {
		return
	}

	err = node.AddChild(content)
	return
}

func (t *Tiapp) AppendWithConf(path string) (err error) {
	json, err := readJson(path)
	additions := json.Get("additions")
	for _, elem := range additions.MustArray() {
		addition := elem.(map[string]interface{})
		xpath := addition["xpath"].(string)
		content := addition["content"].(string)
		err = t.Append(xpath, content)
		if err != nil {
			return
		}
	}
	return
}

func (t *Tiapp) searchWithXPath(xpath string) (node xml.Node, err error) {
	root := t.Document.Root()
	result, err := root.EvalXPath(xpath, nil)
	if err != nil {
		return
	}

	nodes := result.([]xml.Node)
	if len(nodes) == 0 {
		return nil, errors.New("Xml node is not found: " + xpath)
	}
	return nodes[0], nil
}

func (t *Tiapp) Free() {
	if t.Document != nil {
		t.Document.Free()
	}
}

func (t *Tiapp) readDocument() (xml.Document, error) {
	if t.Document != nil {
		return t.Document, nil
	}
	path, err := t.GetFilePath()
	if err != nil {
		return nil, err
	}
	return xml.ReadFile(path, xml.DefaultParseOption)
}
