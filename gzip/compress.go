package gzip

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

type Option struct {
	Level int // range -2 ~ 9
}

type Compressor struct {
	opt *Option
}

func NewCompressor(options ...*Option) *Compressor {
	var opt *Option
	if len(options) != 0 {
		opt = options[0]
	}
	return &Compressor{
		opt: opt,
	}
}

func (m *Compressor) NewWriter(w io.Writer) (writer *gzip.Writer, err error) {
	if m.opt != nil {
		writer, err = gzip.NewWriterLevel(w, m.opt.Level)
		if err != nil {
			return nil, err
		}
	} else {
		writer = gzip.NewWriter(w)
	}
	return writer, nil
}

func (m *Compressor) CompressBytes2Bytes(data []byte) ([]byte, error) {
	var err error
	var gz *gzip.Writer
	var in bytes.Buffer
	gz, err = m.NewWriter(&in)
	if err != nil {
		return nil, err
	}

	_, err = gz.Write(data)
	if err != nil {
		err = gz.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	err = gz.Close()
	if err != nil {
		return nil, err
	}
	return in.Bytes(), nil
}

func (m *Compressor) CompressBytes2File(data []byte, dst string) (int64, error) {
	bs, err := m.CompressBytes2Bytes(data)
	if err != nil {
		return 0, err
	}
	fd, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer fd.Close()
	buf := bytes.NewReader(bs)
	return io.Copy(fd, buf)
}

func (m *Compressor) CompressFile2Bytes(src string) ([]byte, error) {
	fs, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer fs.Close()

	var fi os.FileInfo
	fi, err = fs.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return nil, fmt.Errorf("file %s is a directory", src)
	}
	var in bytes.Buffer
	var gz *gzip.Writer
	reader := bufio.NewReader(fs)
	gz, err = m.NewWriter(&in)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(gz, reader)
	if err != nil {
		err = gz.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	err = gz.Close()
	if err != nil {
		return nil, err
	}
	return in.Bytes(), nil
}

func (m *Compressor) CompressFile2File(src, dst string) (int64, error) {
	//open source file
	fs, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer fs.Close()

	var fi os.FileInfo
	fi, err = fs.Stat()
	if err != nil {
		return 0, err
	}
	if fi.IsDir() {
		return 0, fmt.Errorf("file %s is a directory", src)
	}
	//create dest file
	fd, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_RDWR, fi.Mode())
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	//create gzip writer
	var gz *gzip.Writer
	reader := bufio.NewReader(fs)
	gz, err = m.NewWriter(fd)
	if err != nil {
		return 0, err
	}

	_, err = io.Copy(gz, reader)
	if err != nil {
		err = gz.Close()
		if err != nil {
			return 0, err
		}
		return 0, err
	}
	err = gz.Close()
	if err != nil {
		return 0, err
	}
	fi, err = fd.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func (m *Compressor) CompressBytes2Base64(data []byte) (string, error) {
	comp, err := m.CompressBytes2Bytes(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(comp), nil
}
