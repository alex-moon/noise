 package main

import (
    "io/ioutil"
)

type Ranger struct {
    i int
    obj interface {}
}

func (r Ranger) Range() interface {} {
    switch r.obj.(type) {
    case chan string:
        next, ok := r.obj.(chan string)
        if !ok { panic("channel") }
        return <- next
    case []string:
        next, ok := r.obj.([]string)
        if !ok { panic("array") }
        r.i++  // ah to be young again
        return next[r.i]
    }
    return nil
}

func NewRanger(obj interface {}) Ranger {
    return Ranger{obj: obj, i: -1}  // so we can save some lines of code and just do r.i++
}

// range chan as per https://sites.google.com/site/gopatterns/object-oriented/iterators
// @todo really? A chan is the most efficient way to do this? Pretty sure this is the "generic" solution...
// - perhaps a more generic interface... wrap any rangeable in a Range() method that accounts for differences in range syntax
//  - in progress - see above
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
    filepath := Config().Text.Dir + "/" + text.Uuid()
    file_contents, err := ioutil.ReadFile(filepath)
    if err != nil { panic(err) }
    return file_contents
}

func fsReader() Reader {
    text_files, err := ioutil.ReadDir(Config().Text.Dir)
    if err != nil { panic(err) }
    /* WAS (chan version):
    texts := make(chan TextReader)
    go func () {
        for _, filestat := range text_files {
            uuid := filestat.Name()
            texts <- fsTextReader{string(uuid)}
        }
    } ()
    */
    texts := NewRanger(text_files)
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