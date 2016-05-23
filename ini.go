package ini

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type IniFile struct {
	data map[string]map[string]string
}

// public function
func (o *IniFile) AddSection(section string) bool {
	if _, ok := o.data[section]; ok {
		return false
	} else {
		o.data[section] = make(map[string]string)
	}
	return true
}

func (o *IniFile) RemoveSection(section string) bool {
	if _, ok := o.data[section]; ok {
		return false
	} else {
		// remove key- value
		for s := range o.data[section] {
			delete(o.data[section], s)
		}
		//remove section
		delete(o.data, section)
	}

	return true
}

func (o *IniFile) SetKey(section, key, value string) bool {
	if _, ok := o.data[section][key]; ok {
		return false
	} else {
		o.data[section][key] = value
	}
	return true
}

func (o *IniFile) RemoveKey(section, key string) bool {
	if _, ok := o.data[section][key]; ok {
		return false
	} else {
		delete(o.data[section], key)
	}
	return true
}

// read ini file
func (o *IniFile) ReadIniStream(reader io.Reader) (err error) {
	var r *bufio.Reader
	var line, section, key, value string
	var l, i int

	fmt.Printf("start to read ini file \n")

	r = bufio.NewReader(reader)
	for {
		line, err = r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		line = strings.TrimSpace(line)
		l = len(line)

		switch {
		case l == 0:
			continue
		case line[0] == '#':
			continue
		case line[0] == ';':
			continue
		case line[0] == '[' && line[l-1] == ']':
			section = strings.TrimSpace(line[1:(l - 1)])
			if section != "" {
				o.AddSection(section)
			}
			continue
		default:
			i = strings.IndexAny(line, "=")
			key = strings.TrimSpace(line[0:i])
			value = strings.TrimSpace(line[(i + 1):])
			if section != "" && key != "" {
				o.SetKey(section, key, value)
				//fmt.Printf("================== DEBUG[%s]\n", o.data[section][key])
			}
			continue
		}
	}
	return nil
}

func ReadIniFile(fname string) (c *IniFile, err error) {
	var f *os.File
	if f, err = os.Open(fname); err != nil {
		return nil, err
	} else {
		defer f.Close()
	}

	c = new(IniFile)
	c.data = make(map[string]map[string]string)

	if err = c.ReadIniStream(f); err != nil {
		return nil, err
	}

	return c, err
}

func (o *IniFile) GetValue(section, key string) string {
	s, ok := o.data[section][key]
	if !ok {
		return ""
	}
	return s
}

func (o *IniFile) GetAllSection() []string {
	var sections []string
	for key := range o.data {
		sections = append(sections, key)
	}
	return sections
}

func (o *IniFile) GetAllKeys(section string) []string {
	var keys []string
	for key := range o.data[section] {
		keys = append(keys, key)
	}
	return keys
}

/* Write to file */
func (o *IniFile) WriteToIniStream(writer io.Writer) (err error) {
	buf := bytes.NewBuffer(nil)

	for section, section_map := range o.data {
		_, err = buf.WriteString("\n[" + section + "]\n")
		if err != nil {
			return err
		}

		for key, value := range section_map {
			_, err = buf.WriteString(key + " = " + value + "\n")
			if err != nil {
				return err
			}
		}
	}

	buf.WriteTo(writer)

	return nil
}

func (o *IniFile) WriteToIniFile(fname string) (err error) {
	var f *os.File

	f, err = os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	err = o.WriteToIniStream(f)
	if err != nil {
		return err
	}

	return nil
}

/* test main function */
/*
func main() {
	obj, _ := ReadIniFile("./Bus.ini")

	obj.AddSection("test")
	obj.SetKey("test", "key", "value")

	for _, k := range obj.GetAllSection() {
		fmt.Printf("\n[%s]\n", k)

		for _, key := range obj.GetAllKeys(k) {
			fmt.Printf("%s = %s\n", key, obj.GetValue(k, key))
		}
	}

	obj.WriteToIniFile("./test.ini")

}
*/
