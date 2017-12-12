package gziptest

import (
	"compress/gzip"
	"io"
	"os"
)

//压缩文件Src到Dst
func CompressFile(Dst string, Src string) error {
	newfile, err := os.Create(Dst)
	if err != nil {
		return err
	}
	defer newfile.Close()

	file, err := os.Open(Src)
	if err != nil {
		return err
	}

	zw := gzip.NewWriter(newfile)

	filestat, err := file.Stat()
	if err != nil {
		return nil
	}

	zw.Name = filestat.Name()
	zw.ModTime = filestat.ModTime()
	_, err = io.Copy(zw, file)
	if err != nil {
		return nil
	}

	zw.Flush()
	if err := zw.Close(); err != nil {
		return nil
	}
	return nil
}

//解压文件Src到Dst
func DeCompressFile(Dst string, Src string) error {
	file, err := os.Open(Src)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	newfile, err := os.Create(Dst)
	if err != nil {
		panic(err)
	}
	defer newfile.Close()

	zr, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}

	filestat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	zr.Name = filestat.Name()
	zr.ModTime = filestat.ModTime()
	_, err = io.Copy(newfile, zr)
	if err != nil {
		panic(err)
	}

	if err := zr.Close(); err != nil {
		panic(err)
	}
	return nil
}
