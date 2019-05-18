# TarWriter
This package was created send files to a docker container via [Official Go SDK for Docker](https://github.com/docker/go-docker) more easily.

## Usage with Official Go SDK for Docker
```
var tarArchive bytes.Buffer
tarWriter := tarwriter.NewTarWriter(&tarArchive)

err := tarWriter.AddFolder("somedirectory", tarwriter.AddFolderOptions{IncludeRootFolder: true})
if err != nil{
	return err
}

err := tarWriter.AddFile("somefile.txt", tarwriter.AddFileOptions{})
if err != nil{
	return err
}

err := tarWriter.AddContentFile("someotherfile.txt", []byte("content"), tarwriter.AddFileOptions{})
if err != nil{
	return err
}

tarWriter.Close()

err = cli.CopyToContainer(ctx,
    createResp.ID,
    "/home/user",
    &tarArchive,
    types.CopyToContainerOptions{AllowOverwriteDirWithFile: true},
)
if err != nil {
    return err
}
```

## Usage with HTTP Request
```
var tarArchive bytes.Buffer
tarWriter := tarwriter.NewTarWriter(&tarArchive)

err := tarWriter.AddFolder("somedirectory", tarwriter.AddFolderOptions{IncludeRootFolder: true})
if err != nil{
	return err
}

err := tarWriter.AddFile("somefile.txt", tarwriter.AddFileOptions{})
if err != nil{
	return err
}

err := tarWriter.AddContentFile("someotherfile.txt", []byte("content"), tarwriter.AddFileOptions{})
if err != nil{
	return err
}

tarWriter.Close()

req, err := http.NewRequest(http.MethodPut, "https://someurl.org/archive.tar", &tarArchive)
if err != nil{
	return err
}

```

## Usage with Files
```
var tarArchive bytes.Buffer
tarWriter := tarwriter.NewTarWriter(&tarArchive)

err := tarWriter.AddFolder("somedirectory", tarwriter.AddFolderOptions{IncludeRootFolder: true})
if err != nil{
	return err
}

err := tarWriter.AddFile("somefile.txt", tarwriter.AddFileOptions{})
if err != nil{
	return err
}

err := tarWriter.AddContentFile("someotherfile.txt", []byte("content"), tarwriter.AddFileOptions{})
if err != nil{
	return err
}

tarWriter.Close()

file, err := os.Create("archive.tar")
if err != nil{
	return err
}
defer file.Close()

_, err := io.Copy(file, &tarArchive)
if err != nil{
	return err
}

```
