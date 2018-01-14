package models

import (
	"fmt"
	"github.com/weilaihui/fdfs_client"
)

//fastdfs根据文件名
func FDFSUploadByFileName(filename string) (groupName string, fileId string, err error) {
	fdfsClient, err := fdfs_client.NewFdfsClient("./conf/client.conf")
	if err != nil {
		fmt.Printf("New FdfsClient error %s", err.Error())
		return
	}

	uploadResponse, err := fdfsClient.UploadByFilename(filename)
	if err != nil {
		fmt.Printf("UploadByfilename error %s", err.Error())
		return "", "", err
	}
	fmt.Println(uploadResponse.GroupName)
	fmt.Println(uploadResponse.RemoteFileId)

	return uploadResponse.GroupName, uploadResponse.RemoteFileId, nil
}

//fastdfs 根据buffer上传文件
func FDFSUploadByBuffer(buffer []byte, suffix string) (groupName string, fileId string, err error) {

	fdfsClient, err := fdfs_client.NewFdfsClient("./conf/client.conf")
	if err != nil {
		fmt.Printf("New FdfsClient error %s", err.Error())
		return
	}

	uploadResponse, err := fdfsClient.UploadByBuffer(buffer, suffix)
	if err != nil {
		fmt.Printf("TestUploadByBuffer error %s", err.Error())
		return "", "", err
	}

	fmt.Println(uploadResponse.GroupName)
	fmt.Println(uploadResponse.RemoteFileId)

	return uploadResponse.GroupName, uploadResponse.RemoteFileId, nil
}
