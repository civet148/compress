package gzip

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

type Decompressor struct {
}

func NewDecompressor() *Decompressor {
	return &Decompressor{}
}

func (m *Decompressor) NewReader(r io.Reader) (*gzip.Reader, error) {
	return gzip.NewReader(r)
}

func (m *Decompressor) DecompressBytes2Bytes(data []byte) ([]byte, error) {
	var err error
	var gz *gzip.Reader
	var in bytes.Buffer
	var out bytes.Buffer
	_, err = in.Write(data)
	if err != nil {
		return nil, err
	}
	gz, err = m.NewReader(&in)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(&out, gz)
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
	return out.Bytes(), nil
}

func (m *Decompressor) DecompressBytes2File(data []byte, dst string) (int64, error) {
	compr, err := m.DecompressBytes2Bytes(data)
	//create dest file
	fd, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	reader := bytes.NewReader(compr)
	_, err = io.Copy(fd, reader)
	if err != nil {
		return 0, err
	}
	var fi os.FileInfo
	fi, err = fd.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func (m *Decompressor) DecompressFile2Bytes(src string) ([]byte, error) {
	f, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	var fi os.FileInfo
	fi, err = f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return nil, fmt.Errorf("file %s is a directory", src)
	}
	var out bytes.Buffer
	var gz *gzip.Reader
	gz, err = m.NewReader(f)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(&out, gz)
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
	return out.Bytes(), nil
}

func (m *Decompressor) DecompressFile2File(src, dst string) (int64, error) {
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

	//create gzip reader
	var gz *gzip.Reader
	gz, err = m.NewReader(fs)
	if err != nil {
		return 0, err
	}
	_, err = io.Copy(fd, gz)
	if err != nil {
		_ = gz.Close()
		return 0, err
	}
	_ = gz.Close()
	fi, err = fd.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func (m *Decompressor) DecompressBytes2Base64(data []byte) (string, error) {
	dec, err := m.DecompressBytes2Bytes(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(dec), nil
}

func (m *Decompressor) DecompressFile2Base64(src string) (string, error) {
	dec, err := m.DecompressFile2Bytes(src)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(dec), nil
}
