package main

import (
	"github.com/civet148/compress/gzip"
	"github.com/civet148/log"
)

const (
	txtMetricsFile = "metrics.txt"
	gzMetricsFile  = "metrics.gz"
	gzPlainFile    = "plain.gz"
	decFile        = "metrics-decompressed.txt"
	plainText      = `go_gc_duration_seconds{quantile="0"} 0.000114827
go_gc_duration_seconds{quantile="0.25"} 0.000134637
go_gc_duration_seconds{quantile="0.5"} 0.000141811
go_gc_duration_seconds{quantile="0.75"} 0.000153749
go_gc_duration_seconds{quantile="1"} 0.001029093
go_gc_duration_seconds_sum 2.240228923
go_gc_duration_seconds_count 14908`
)

func main() {
	var err error
	com := gzip.NewCompressor()
	dec := gzip.NewDecompressor()

	err = GzipBytes2Bytes(com, dec)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	err = GzipBytes2Base64(com, dec)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	err = GzipFile2Bytes(com, dec)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	err = GzipBytes2File(com, dec)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	err = GzipFile2File(com, dec)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
}

func GzipBytes2Bytes(c *gzip.Compressor, d *gzip.Decompressor) (err error) {
	log.Infof("----------------------------------------------------------------")
	var cd, dd []byte
	cd, err = c.CompressBytes2Bytes([]byte(plainText))
	if err != nil {
		return log.Errorf(err.Error())
	}
	log.Infof("plain text size %d bytes, compressed size %v", len(plainText), len(cd))
	dd, err = d.DecompressBytes2Bytes(cd)
	log.Infof("decompressed text size [%v]", len(dd))
	return
}

func GzipBytes2Base64(c *gzip.Compressor, d *gzip.Decompressor) (err error) {
	log.Infof("----------------------------------------------------------------")
	var cd string
	cd, err = c.CompressBytes2Base64([]byte(plainText))
	if err != nil {
		return log.Errorf(err.Error())
	}
	log.Infof("plain text size %d bytes, compressed size %v", len(plainText), len(cd))
	return
}

func GzipFile2Bytes(c *gzip.Compressor, d *gzip.Decompressor) (err error) {
	log.Infof("----------------------------------------------------------------")
	var cd, dd []byte
	cd, err = c.CompressFile2Bytes(txtMetricsFile)
	if err != nil {
		return log.Errorf(err.Error())
	}
	log.Infof("file %s compressed size %d bytes", txtMetricsFile, len(cd))
	dd, err = d.DecompressBytes2Bytes(cd)
	log.Infof("file %s decompressed size [%v]", txtMetricsFile, len(dd))
	return
}

func GzipBytes2File(c *gzip.Compressor, d *gzip.Decompressor) (err error) {
	var n int64
	log.Infof("----------------------------------------------------------------")
	log.Infof("plain text size [%v]", len(plainText))
	n, err = c.CompressBytes2File([]byte(plainText), gzPlainFile)
	if err != nil {
		return log.Errorf(err.Error())
	}
	log.Infof("write to %s compressed size %d bytes", gzPlainFile, n)
	var data []byte
	data, err = d.DecompressFile2Bytes(gzPlainFile)
	if err != nil {
		return log.Errorf(err.Error())
	}
	log.Infof("read from %s decompressed size %d bytes", gzPlainFile, len(data))
	return
}

func GzipFile2File(c *gzip.Compressor, d *gzip.Decompressor) (err error) {
	log.Infof("----------------------------------------------------------------")
	var n int64
	n, err = c.CompressFile2File(txtMetricsFile, gzMetricsFile)
	if err != nil {
		return log.Errorf(err.Error())
	}
	log.Infof("write to %s compressed size %d bytes", gzMetricsFile, n)
	n, err = d.DecompressFile2File(gzMetricsFile, decFile)
	log.Infof("write to %s compressed size %d bytes", decFile, n)
	return
}
