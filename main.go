package main

import (
	"fmt"
	"path"
)

func main() {
	makeContent, makeFErr := readFile("./TagMake")
	if makeFErr != nil {
		fmt.Println(makeFErr.Error())
		return
	}

	madeTag, outpath, madeTagErr := interpretTagmake(makeContent)
	if madeTagErr != nil {
		fmt.Println(madeTagErr.Error())
		return
	}

	if len(outpath) >= 5 {
		_, f := path.Split(outpath)
		if outpath[len(outpath)-5:] != ".json" {
			outpath += ".json"
		} else if f == ".json" {
			fmt.Println("error: output path cannot be empty")
			return
		}
	} else {
		madeTag = "o_" + madeTag + ".json"
	}

	tagWritingErr := writeFile(outpath, madeTag)
	if tagWritingErr != nil {
		fmt.Println(tagWritingErr)
		return
	}
}
