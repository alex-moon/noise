 package main

import (
    "io/ioutil"
)

var text_dir string = Config().Text.Dir

type Reader struct {
    text_id string
    file_contents []byte
}

func NewReader(text_id string) Reader {
    // files, err := ioutil.ReadDir(text_dir)
    // if err != nil { panic(err) }

    var filepath string = text_dir + "/" + text_id

    file_contents, err := ioutil.ReadFile(filepath)
    if err != nil { panic(err) }

    return Reader{text_id, file_contents}
}

func (r Reader) ReadAll() []byte {
    return r.file_contents
}