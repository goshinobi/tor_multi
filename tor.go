package tor

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"time"
)

var (
	HOME_DIR        = os.Getenv("HOME")
	TOR_MULTI_DIR   = HOME_DIR + "/.tor_multi"
	TOR_LIB         = TOR_MULTI_DIR + "/lib"
	LATEST_PROXY_ID = 0
	PROXY_LIST      = map[int]*ProxyInfo{}
)

type ProxyInfo struct {
	Conf   *TorConf
	IP     string
	UpTime time.Time
	Cmd    *exec.Cmd
}

type TorConf struct {
	ID            *int   `json:"id"`
	SocksPort     int    `json:"socks_port"`
	ControlPort   int    `json:"control_port"`
	DataDirectory string `json:"data_directory"`
	ConfPath      string `conf_path`
}

func init() {
	//create tor conf dir
	if !isExist(TOR_MULTI_DIR) {
		err := os.Mkdir(TOR_MULTI_DIR, 0755)
		if err != nil {
			panic(err)
		}
	}
	if !isExist(TOR_LIB) {
		err := os.Mkdir(TOR_LIB, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func isExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func CreateTorConf(n int) (*TorConf, error) {
	c := &TorConf{
		nil,
		getPort(),
		getPort(),
		fmt.Sprintf("%s/tor%d", TOR_LIB, n),
		fmt.Sprintf("%s/torrc.%d", TOR_MULTI_DIR, n),
	}
	conf := fmt.Sprintf("SocksPort %d\n", c.SocksPort)
	conf += fmt.Sprintf("ControlPort %d\n", c.ControlPort)
	conf += fmt.Sprintf("DataDirectory %s\n", c.DataDirectory)
	err := ioutil.WriteFile(fmt.Sprintf("%s/torrc.%d", TOR_MULTI_DIR, n), []byte(conf), 0644)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func StartProxy() (*TorConf, error) {
	proxyID := LATEST_PROXY_ID
	LATEST_PROXY_ID++

	conf, err := CreateTorConf(proxyID)
	conf.ID = &proxyID
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("tor", "-f", conf.ConfPath)
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	PROXY_LIST[proxyID] = &ProxyInfo{
		conf,
		"",
		time.Now(),
		cmd,
	}
	return conf, nil
}

func KillTorProxy(n int) error {
	target, ok := PROXY_LIST[n]
	if !ok {
		return errors.New(fmt.Sprintf("%d is not existed", n))
	}
	err := target.Cmd.Process.Kill()
	if err != nil {
		return err
	}
	err = DeleteTorConf(target.Conf)
	if err != nil {
		return err
	}
	delete(PROXY_LIST, n)
	return nil
}

func GetWorkProxyList() map[int]*ProxyInfo {
	return PROXY_LIST
}

func DeleteTorConf(c *TorConf) (err error) {
	err = os.RemoveAll(c.DataDirectory)
	if err != nil {
		return err
	}
	return os.RemoveAll(c.ConfPath)
}

func getPort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}
