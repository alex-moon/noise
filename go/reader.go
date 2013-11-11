 package main

import (
    "fmt"
    "io/ioutil"
)

// range chan as per https://sites.google.com/site/gopatterns/object-oriented/iterators
type Reader struct {
    texts chan TextReader
}

type TextReader interface {
    Uuid() string
    Bytes() []byte
}


// FILESYSTEM
type fsTextReader struct {
    uuid string
}

func (text fsTextReader) Uuid() string {
    return text.uuid
}

func (text fsTextReader) Bytes() []byte {
    fmt.Printf("Bytes() called!\n")
    filepath := Config().Text.Dir + "/" + text.Uuid()
    file_contents, err := ioutil.ReadFile(filepath)
    if err != nil { panic(err) }
    return file_contents
}

func fsReader() Reader {
    text_files, err := ioutil.ReadDir(Config().Text.Dir)
    if err != nil { panic(err) }
    texts := make(chan TextReader)
    go func () {
        for _, filestat := range text_files {
            uuid := filestat.Name()
            texts <- fsTextReader{string(uuid)}
        }
    } ()
    return Reader{texts}
}


// REDIS
/*
            text_uuid, err := c.conn.Do("RPOP", c.channel)
            if err == nil && text_uuid != nil {
                uuid, ok := text_uuid.([]byte)
                if ok {
                    string(uuid)
*/

func NewReader() Reader {
    // @todo read config for type of reader (Redis vs SQL vs Mongo vs Solr vs FS vs RSS vs REST vs scraping vs whatever)
    return fsReader()
}