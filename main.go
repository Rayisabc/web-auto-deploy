// Author: Rui
// Date: 2023/01/05 16:22
// Description: compress the local web dist,and upload to the docker the server,invoke the remote the shell deploy docker images

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
	"web-auto-deploy/model"
	"web-auto-deploy/util"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var (
	user, password, scriptName, host, localFilePath, remoteFilePath string
	executeShell                                                    bool
	sleepSeconds                                                    int = 10
)

const fileName string = "dist.zip"
const serviceName string = "WEB-Deploy"

// do some init things
func init() {
	initLogFile()

	initParams()
}

func main() {
	loopMonitor(localFilePath)
}

func loopMonitor(localDir string) {
	var fileSize int64 = 0
	var lastModifiedTime time.Time
	for {
		fileInfo, err := os.Stat(localDir)
		if err != nil {
			util.Error.Println("Failed to get directory info")
			continue
		}

		if fileSize == 0 || lastModifiedTime.IsZero() {
			lastModifiedTime = fileInfo.ModTime()
			fileSize, _ = util.GetDirSize(localDir)
			continue
		}

		newFileSize, _ := util.GetDirSize(localDir)

		if fileSize != newFileSize || lastModifiedTime.Before(fileInfo.ModTime()) {
			fileSize = newFileSize
			lastModifiedTime = fileInfo.ModTime()
			go doSync()
			sleepSeconds = 60
		} else {
			log.Println("### no change found")
			sleepSeconds = 10
		}

		time.Sleep(time.Duration(sleepSeconds) * time.Second)
	}
}

func doSync() {
	util.Info.Printf("### %s start \n", serviceName)

	// zip file
	localZipFile := strings.TrimSuffix(localFilePath, "dist/") + fileName

	CompressUtil := &util.CompressUtil{}
	err := CompressUtil.Compress(localFilePath, localZipFile)
	if err != nil {
		panic(err)
	}

	util.Info.Println("### Compress Finished")

	// create connection
	client, sftpClient := connectServer(host, user, password)
	defer client.Close()
	defer sftpClient.Close()

	// copy file
	fileOperation(localZipFile, fileName, sftpClient, remoteFilePath)

	util.Info.Printf("### File[%s]upload success \n", fileName)

	executeCommand(remoteFilePath, scriptName, client)

	// remove local zip file
	os.Remove(localZipFile)

	util.Info.Printf("### %s Finished", serviceName)
}

func initLogFile() {
	file, err := os.OpenFile("web-deploy.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		util.Error.Fatal("Failed to open log file:", err)
	}

	util.Info = newLog("INFO: ", file)
	util.Warn = newLog("Warn: ", file)
	util.Error = newLog("Error: ", file)
	util.Fatal = newLog("Fatal: ", file)

	fmt.Fprintln(file, "")
}

func newLog(level string, file *os.File) *log.Logger {
	return log.New(file, level, log.Ldate|log.Ltime|log.Lshortfile)
}

func initParams() {
	cfg := &util.CfgUtil{}
	var config model.Config = cfg.ReadConfig()

	host = config.ServerIp
	user = config.Username
	password = config.Password
	scriptName = config.ScriptName
	remoteFilePath = config.ServerDir
	localFilePath = config.LocalDirectory
	executeShell = config.ExecuteShell

	// Make sure the path end with "/"
	if !strings.HasSuffix(localFilePath, "/") {
		localFilePath += "/"
	}

	if !strings.HasSuffix(remoteFilePath, "/") {
		remoteFilePath += "/"
	}
}

func fileOperation(localZipFile string, fileName string, sftpClient *sftp.Client, remoteFilePath string) {
	// open local file
	localFile, err := os.Open(localZipFile)
	if err != nil {
		panic(err)
	}
	defer localFile.Close()

	// create remote file
	remoteFile, err := sftpClient.Create(remoteFilePath + fileName)
	if err != nil {
		panic(err)
	}
	defer remoteFile.Close()

	// copy file from local to remoute
	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		panic(err)
	}

	_, err = remoteFile.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}
}

func connectServer(host string, user string, password string) (c *ssh.Client, s *sftp.Client) {
	// connect ssh
	client, err := ssh.Dial("tcp", host+":22", &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		util.Error.Printf("network error: %s", err)
		panic(err)
	}

	// create sftp client
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		util.Error.Printf("sftp network error: %s", err)
		panic(err)
	}

	return client, sftpClient
}

func executeCommand(remoteFilePath string, scriptName string, client *ssh.Client) {
	// create session
	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var commands = []string{
		fmt.Sprintf("mkdir -p %sdist", remoteFilePath),
		fmt.Sprintf("unzip %sdist.zip -d %sdist", remoteFilePath, remoteFilePath),
		fmt.Sprintf("rm -rf %sdist.zip", remoteFilePath),
	}
	// need execute shell or not
	if executeShell {
		util.Info.Println("### Execute remote shell script")
		commands = append(commands, fmt.Sprintf("sh %s%s", remoteFilePath, scriptName))
	} else {
		util.Info.Println("### Dont't execute remote shell script")
	}

	// join commands
	cmdStr := strings.Join(commands, "; ")

	// execute commands
	output, err := session.CombinedOutput(cmdStr)
	if err != nil {
		util.Error.Printf("[Error]Command finished with error: %v\n", err)
	} else {
		util.Info.Printf("Output of '%s':\n%s\n", cmdStr, string(output))
	}

}
