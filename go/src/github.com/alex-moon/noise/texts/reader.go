package texts

import (
    "io/ioutil"
)

type Reader interface {
    Bytes() []byte
    Uuid() string
}

type FileSystemReader struct {
    text_dir string
    uuid string
}

func NewFileSystemReader(text_dir, uuid string) FileSystemReader {
    return FileSystemReader{text_dir, uuid}
}

func (r FileSystemReader) Bytes() []byte {
    file_contents, err := ioutil.ReadFile(r.text_dir + "/" + r.uuid)
    if err != nil { panic(err) }
    return file_contents
}

func (r FileSystemReader) Uuid() string {
    return r.uuid
}