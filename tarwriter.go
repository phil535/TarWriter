package tarwriter

import(
	"io"
	"os"
	"path/filepath"
	"archive/tar"
	"errors"
	"io/ioutil"
	"github.com/phayes/permbits"
	"time"
)


type TarWriter struct{
	tar.Writer
}
func NewTarWriter(w io.Writer)*TarWriter{
	return &TarWriter{Writer: *tar.NewWriter(w)}
}

type AddFileOptions struct {
	InternalPath string
	ForceExecutableFlags bool
}
func (tf *TarWriter)AddFile(filePath string, options AddFileOptions)error{

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return err
	}

	fileMode := permbits.FileMode(fileStat.Mode())
	if options.ForceExecutableFlags {
		fileMode.SetUserExecute(true)
		fileMode.SetGroupExecute(true)
		fileMode.SetOtherExecute(true)
	}

	fileHeader := &tar.Header{
		Name: filepath.Join(options.InternalPath, filepath.Base(filePath)),
		Mode: int64(fileMode),
		Size: fileStat.Size(),
		ModTime: fileStat.ModTime(),
	}

	err = tf.WriteHeader(fileHeader)
	if err != nil {
		return err
	}

	_, err = io.Copy(tf, file)
	if err != nil {
		return err
	}

	return nil
}
func (tf *TarWriter)AddFiles(filePaths []string, options AddFileOptions)error {
	for _, filePath := range filePaths{
		err := tf.AddFile(filePath, options)
		if err != nil{
			return err
		}
	}
	return nil
}

func (tf *TarWriter)AddContentFile(filename string, content []byte, options AddFileOptions)error{

	if len(filename) == 0{
		return errors.New("Invalid filename")
	}

	var fileMode permbits.PermissionBits

	fileMode.SetUserRead(true)
	fileMode.SetUserWrite(true)
	fileMode.SetGroupRead(true)
	fileMode.SetGroupWrite(true)
	fileMode.SetOtherRead(true)
	fileMode.SetOtherWrite(true)

	if options.ForceExecutableFlags {
		fileMode.SetUserExecute(true)
		fileMode.SetGroupExecute(true)
		fileMode.SetOtherExecute(true)
	}

	fileHeader := &tar.Header{
		Name: filepath.Join(options.InternalPath, filename),
		Mode: int64(fileMode),
		Size: int64(len(content)),
		ModTime: time.Now(),
	}

	err := tf.WriteHeader(fileHeader)
	if err != nil {
		return err
	}

	_, err = tf.Write(content)
	if err != nil {
		return err
	}

	return nil
}


type AddFolderOptions struct {
	InternalPath string
	ForceExecutableFlags bool
	NonRecursive bool
	IncludeRootFolder bool
}

func (tf *TarWriter)AddEmptyFolder(folderPath, internalPath string)error {

	folderStats, err := os.Stat(folderPath)
	if err != nil{
		return err
	}
	if !folderStats.IsDir(){
		return os.PathError{Op: "IsDir", Path: folderPath, Err: os.ErrInvalid}.Err
	}

	folderHeader := &tar.Header{
		Name: internalPath,
		Mode: int64(folderStats.Mode()),
		Size: 0, // directories have no contents
		ModTime: folderStats.ModTime(),
	}
	err = tf.WriteHeader(folderHeader)
	if err != nil {
		return err
	}
	return nil
}

func (tf *TarWriter)AddFolder(folderPath string, options AddFolderOptions)error{

	folderStats, err := ioutil.ReadDir(folderPath)
	if err != nil{
		return err
	}

	rootFolder := filepath.Base(folderPath)
	if options.IncludeRootFolder{
		options.InternalPath = filepath.Join(options.InternalPath, rootFolder)
	}else{
		if len(folderStats) == 0{
			err = tf.AddEmptyFolder(folderPath, filepath.Join(options.InternalPath, rootFolder))
			if err != nil{
				return err
			}
			return nil
		}
	}

	for _, fileStat := range folderStats{
		if fileStat.IsDir(){
			if options.NonRecursive{
				continue
			}
			err = tf.AddFolder(filepath.Join(folderPath, fileStat.Name()),
				AddFolderOptions{
					InternalPath: filepath.Join(options.InternalPath, fileStat.Name()),
					NonRecursive: false,
					ForceExecutableFlags: options.ForceExecutableFlags,
					IncludeRootFolder: false,
				})
			if err != nil{
				return err
			}
		}else{
			err = tf.AddFile(
				filepath.Join(folderPath, fileStat.Name()),
				AddFileOptions{
					InternalPath: options.InternalPath,
					ForceExecutableFlags: options.ForceExecutableFlags,
				})
			if err != nil{
				return err
			}
		}
	}
	return nil
}
func (tf *TarWriter)AddFolders(folderPaths []string, options AddFolderOptions)error {
	for _, folderPath := range folderPaths{
		err := tf.AddFolder(folderPath, options)
		if err != nil{
			return err
		}
	}
	return nil
}
