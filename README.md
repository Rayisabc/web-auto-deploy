# Web Dist Auto Deploy

The program automatically monitors the `dist` folder of a local web frontend project.
When the modification time or size of the folder changes, it automatically packages the `dist` folder and uploads it to the Docker server.
It then invokes a shell script on the server to build the Docker image and push it to the Docker image registry.

## Description

The program automatically monitors the `dist` folder of a local web frontend project.
When the modification time or size of the folder changes, it automatically packages the `dist` folder and uploads it to the {ServerIp} server.
It then invokes a shell script on the server to build the Docker image and push it to the Docker image registry.
of course,you can disabled the configuration of the shell by setting the `ExecuteShell` variable to `false` in the `config.toml` file.

## Getting Started

1、make sure you have installed the Golang SDK on your computer.  
2、Clone the project to your local disk.  
3、Update the configuration in the `config.toml` file.  
4、Run the project with the following command:  

```bash
go run main.go
```

## Notice
This program is only for personal learning and research, please do not use it for commercial purposes.   
This program is for reference only, and its legality, accuracy, and effectiveness are not guaranteed. Please judge according to the situation.   
Due to the use of third-party libraries, there may be security risks. Please read the relevant documentation carefully before using it.

## License

"MIT"
