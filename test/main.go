package main

import (
	"fatty0860/ini"
	"fmt"
)

func main() {
	obj, _ := ini.ReadIniFile("./test.ini")
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
