package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

var ip_a = flag.String("ip", "127.0.0.1", "Enter the ip")
var algoritm_pick = flag.String("alg", "3", "enter the alg")
var filename = "node_exporter.service"
var url = "https://github.com/prometheus/node_exporter/releases/download/v1.9.1/node_exporter-1.9.1.linux-amd64.tar.gz"
var filelocbin = "/usr/local/bin/"
var list_filename_algoritm = []string{"id_ecdsa", "id_rsa", "id_ed25519"}

func sshkeysetting(number *string, list []string) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(errors.New("error occured while detecting home directory"))
	}
	value, err := strconv.Atoi(*number)
	if err != nil {
		fmt.Println(errors.New("error has occured while getting number"))
	}
	value_alg := list[value-1]
	ssh_keygen := exec.Command("ssh-keygen", "-N", " ", "-f", fmt.Sprintf("%s/.ssh/%s", homedir, value_alg))
	ssh_keygen.Stdout = os.Stdout
	if err = ssh_keygen.Run(); err != nil {
		ssh_keygen.Stderr = os.Stderr
		print(value_alg)
		print(homedir)
		fmt.Println(errors.New("error while running command"))
	}
	ssh_copy_id := exec.Command("ssh-copy-id", "-i") //fmt.Sprintf("%s/.ssh/%s.pub", homedir, value_alg)
	ssh_copy_id.Stdout = os.Stdout
	if err := ssh_copy_id.Run(); err != nil {
		fmt.Println(errors.New("error while copy public key"))
		ssh_copy_id.Stderr = os.Stderr
	}
}

func GetFile(ip *string) {

	_, err := os.Stat("./node_exporter-1.9.1.linux-amd64.tar.gz")
	if os.IsNotExist(err) {
		get_url := exec.Command("wget", url)
		if err := get_url.Run(); err != nil {
			fmt.Println(errors.New("error has occured while getting an url"))
		}
		extraxt_file := exec.Command("sudo", "tar", "-xvf", "node_exporter-1.9.1.linux-amd64.tar.gz")
		extraxt_file.Stdout = os.Stdout
		if err := extraxt_file.Run(); err != nil {
			fmt.Println(err.Error())
		}
	}
	sendbinnary := exec.Command("scp", "node_exporter-1.9.1.linux-amd64/node_exporter", fmt.Sprintf("root@%s:/tmp", *ip))
	sendbinnary.Stdout = os.Stdout
	if err := sendbinnary.Run(); err != nil {
		sendbinnary.Stderr = os.Stderr
		fmt.Println(err.Error())
	}

	setting_permission := exec.Command("ssh", fmt.Sprintf("root@%s", *ip), "mv", "/tmp/node_exporter", filelocbin) //, "&&", "chmod", "+x", "/usr/local/bin/node_exporter")
	if err := setting_permission.Run(); err != nil {
		fmt.Println(errors.New("error while transfer a file"))
		setting_permission.Stderr = os.Stderr
		os.Exit(1)
	}
	changingFileDir := exec.Command("ssh", fmt.Sprintf("root@%s", *ip), "chmod", "+x", "/usr/local/bin/node_exporter")
	if err := changingFileDir.Run(); err != nil {
		fmt.Println(errors.New("error while setting permission"))
		changingFileDir.Stderr = os.Stderr
		os.Exit(1)
	}

}
func SendFIle(ip *string) {
	sendfile_cmd := exec.Command("scp", filename, fmt.Sprintf("root@%s:/etc/systemd/system/", *ip))
	sendfile_cmd.Stdout = os.Stdout
	sendfile_cmd.Run()
}

func CreateServiceFile() {
	serviceconfig := `[Unit]
Description=Prometheus exporter for machine metrics

[Service]
Restart=always
User=prometheus
ExecStart=/usr/local/bin/node_exporter
ExecReload=/bin/kill -HUP $MAINPID
TimeoutStopSec=20s
SendSIGKILL=no

[Install]
WantedBy=multi-user.target`

	create_file, err := os.OpenFile("node_exporter.service", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)

	if err != nil {
		fmt.Println(errors.New("error while creating a file"))
	}
	defer create_file.Close()
	create_file.WriteString(serviceconfig)
	fmt.Println("File service success")
}

func restart_cmd(ip *string) {
	enable := exec.Command("ssh", fmt.Sprintf("root@%s", *ip), "systemctl", "enable", "--now", filename)
	start := exec.Command("ssh", fmt.Sprintf("root@%s", *ip), "systemctl", "start", filename)
	enable.Stdout = os.Stdout
	start.Stdout = os.Stdout
	enable.Stderr = os.Stderr
	start.Stderr = os.Stderr
	if err := enable.Run(); err != nil {
		fmt.Println(errors.New("error has occured while enabling the service"))
	}
	if err := start.Run(); err != nil {
		fmt.Println(errors.New("error has occured while starting the service"))
	}
	//enable.Run()
}

func create_user(ip *string) {
	checkUser := exec.Command("ssh", fmt.Sprintf("root@%s", *ip), "id", "prometheus")
	if err := checkUser.Run(); err != nil {
		create_user_cmd := exec.Command("ssh", fmt.Sprintf("root@%s", *ip), "useradd", "--system", "-s", "/sbin/nologin", "prometheus")
		create_user_cmd.Stdout = os.Stdout
		create_user_cmd.Stderr = os.Stderr
		if err := create_user_cmd.Run(); err != nil {
			fmt.Println(errors.New("error has occured while creating user"))
		} else {
			fmt.Println("succes")
		}
	} else {
		fmt.Println("User Exist")
	}

}
func main() {
	flag.Parse()
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err.Error())
	}

	_, error_home := os.Stat(homedir)

	if os.IsNotExist(error_home) {
		sshkeysetting(algoritm_pick, list_filename_algoritm)
	}

	//download file from url
	GetFile(ip_a)

	//create file node_exporter
	CreateServiceFile()

	//sendfile to server
	SendFIle(ip_a)

	//create user
	create_user(ip_a)

	//restart command
	restart_cmd(ip_a)

}
