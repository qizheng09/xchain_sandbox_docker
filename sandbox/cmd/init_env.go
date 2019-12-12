package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
	initFlags()
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init the sandbox enviroment",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Now initializing the sandbox enviroment.Please wait...")
		if err := initEnv(); err != nil {
			fmt.Println("Initialize failed! err = ", err.Error())
			return
		}
		fmt.Println("Initialize successful!")
	},
}

var (
	nodeNumber      int64
	minerNumber     int64
	withoutInitConf bool
	xRoot           string
	sandRoot        string
	initRpcPort     int64
	initP2pPort     int64
	userCurr        string
)

func initFlags() {
	xRoot = os.Getenv("XCHAIN_ROOT")
	sandRoot = os.Getenv("XCHAIN_SAND_ROOT")
	userC, _ := user.Current()
	userCurr = userC.Username
	initCmd.Flags().Int64VarP(&nodeNumber, "nodeNumber", "N", 5, "The number of nodes to start")
	initCmd.Flags().Int64VarP(&minerNumber, "minerNumber", "M", 3, "The number of nodes to start")
	initCmd.Flags().BoolVarP(&withoutInitConf, "withoutInitConf", "", false, "The flag whether to init `xchain.yaml` and `xuper.json`")
	initCmd.Flags().Int64VarP(&initRpcPort, "initRpcPort", "", 37101, "The init rpc mapping port")
	initCmd.Flags().Int64VarP(&initP2pPort, "initP2pPort", "", 47101, "The init p2p mapping port")
}

// initEnv will init the multi xchain sandbox
// update binary
// init NodesFiles
// init configs: xchain.yaml, xuper.json
// init docker-compose.yml
// TODO: @DhunterAO Now the network mode of xchain container is bridge,
// the performance of bridge is pretty poor. If the nodes of your network more than 30,
// you'd better choose host mode of your container. The tool will support this mode in near feature.
func initEnv() error {
	if xRoot == "" {
		return errors.New("The XCHAIN_ROOT environment variable have not been set")
	}
	if sandRoot == "" {
		return errors.New("The XCHAIN_SAND_ROOT environment variable have not been set")
	}
	if nodeNumber > 100 {
		return errors.New("The nodeNumber can not bigger than 50 in one machine")
	}
	if userCurr == "" {
		return errors.New("Get current username error")
	}
	if err := updateBinary(); err != nil {
		return err
	}

	if err := initNodesFiles(nodeNumber); err != nil {
		return err
	}

	if err := initConf(); err != nil {
		return err
	}

	if err := initDockerCompose(); err != nil {
		return err
	}
	return nil
}

// updateBinary will update all binary from {{XCHAIN_ROOT}} environment
func updateBinary() error {
	// update bin
	sorcXchainPath := xRoot + "/output/xchain"
	destXchainPath := sandRoot + "/bin/xchain"
	if err := copyFile(sorcXchainPath, destXchainPath); err != nil {
		return err
	}

	sorcCliPath := xRoot + "/output/xchain-cli"
	destCliPath := sandRoot + "/bin/xchain-cli"
	if err := copyFile(sorcCliPath, destCliPath); err != nil {
		return err
	}

	sorcWasmPath := xRoot + "/output/wasm2c"
	destWasmPath := sandRoot + "/bin/wasm2c"
	if err := copyFile(sorcWasmPath, destWasmPath); err != nil {
		return err
	}
	// update plugins
	sorcPluginsPath := xRoot + "/output/plugins"
	destPluginsPath := sandRoot + "/plugins/"
	return copyDir(sorcPluginsPath, destPluginsPath)
}

// copyFile copy src file to dst, will overwrite the dst file
func copyFile(src, dst string) error {
	fmt.Println("copy", src, "to", dst)
	if _, err := os.Stat(dst); err == nil {
		if err = os.Remove(dst); err != nil {
			return err
		}
	}

	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, input, 0755)
	if err != nil {
		fmt.Println("Error creating", dst)
		return err
	}
	return nil
}

// copyDir copy dir
func copyDir(src, dst string) error {
	fmt.Println("copyDir", src, "to", dst)
	fs, _ := ioutil.ReadDir(src)
	for _, v := range fs {
		if v.IsDir() {
			fss, _ := ioutil.ReadDir(src + "/" + v.Name())
			for _, vs := range fss {
				os.Mkdir(dst+v.Name(), 0755)
				srcFile := src + "/" + v.Name() + "/" + vs.Name()
				dstFile := dst + v.Name() + "/" + vs.Name()
				if err := copyFile(srcFile, dstFile); err != nil {
					return err
				}
			}
		} else {
			srcFile := src + "/" + v.Name()
			dstFile := dst + v.Name()
			if err := copyFile(srcFile, dstFile); err != nil {
				return err
			}
		}
	}
	return nil
}

// initNodesFiles will init nodes date/keys, data/netKeys and create filefolder needed
func initNodesFiles(nodeNums int64) error {
	if nodeNums < 0 {
		return errors.New("Node number can not less than 0")
	}
	dstNodesPath := sandRoot + "/nodes"
	fmt.Println("initNodesFiles Path=", dstNodesPath)
	os.RemoveAll(dstNodesPath)
	if err := os.Mkdir(dstNodesPath, 0755); err != nil {
		return err
	}
	for i := 1; i <= int(nodeNums); i++ {
		if err := initNode(i); err != nil {
			return err
		}
	}
	return nil
}

func initNode(i int) error {
	dstNodePath := sandRoot + "/nodes/node" + strconv.Itoa(i)
	dstNodeDataPath := dstNodePath + "/data"

	// create file folder
	if err := os.Mkdir(dstNodePath, 0755); err != nil {
		return err
	}
	if err := os.Mkdir(dstNodeDataPath, 0755); err != nil {
		return err
	}
	if err := os.Mkdir(dstNodePath+"/conf", 0755); err != nil {
		return err
	}
	if err := os.Mkdir(dstNodeDataPath+"/blockchain", 0755); err != nil {
		return err
	}
	if err := os.Mkdir(dstNodeDataPath+"/config", 0755); err != nil {
		return err
	}
	if err := os.Mkdir(dstNodeDataPath+"/keys", 0755); err != nil {
		return err
	}
	if err := os.Mkdir(dstNodeDataPath+"/netkeys", 0755); err != nil {
		return err
	}

	if err := os.Mkdir(dstNodePath+"/plugins", 0755); err != nil {
		return err
	}

	if err := os.Mkdir(dstNodePath+"/plugins/consensus", 0755); err != nil {
		return err
	}

	if err := os.Mkdir(dstNodePath+"/plugins/contract", 0755); err != nil {
		return err
	}

	if err := os.Mkdir(dstNodePath+"/plugins/crypto", 0755); err != nil {
		return err
	}

	if err := os.Mkdir(dstNodePath+"/plugins/kv", 0755); err != nil {
		return err
	}

	srcPluginsPath := sandRoot + "/plugins"
	dstPluginsPath := dstNodePath + "/plugins/"
	if err := copyDir(srcPluginsPath, dstPluginsPath); err != nil {
		return err
	}

	// init keys
	cmdKeyStr := sandRoot + "/bin/xchain-cli account newkeys -f -o " + dstNodeDataPath + "/keys"
	println("cmdKeyStr", cmdKeyStr)
	cmdKey := exec.Command("bash", "-c", cmdKeyStr)
	if err := cmdKey.Run(); err != nil {
		return err
	}
	cmdNetKeyStr := sandRoot + "/bin/xchain-cli netURL gen --path " + dstNodeDataPath + "/netkeys/"
	println("cmdNetKeyStr", cmdNetKeyStr)
	cmdNetKey := exec.Command("bash", "-c", cmdNetKeyStr)
	if err := cmdNetKey.Run(); err != nil {
		return err
	}
	return nil
}

func copySeedNodeConf() error {
	srcPath := sandRoot + "/conf/"
	dstNodePath := sandRoot + "/nodes/node1"
	if err := copyFile(srcPath+"xchain.yaml.seed.tpl", dstNodePath+"/conf/xchain.yaml"); err != nil {
		return err
	}
	return nil
}

func copyNodeConf(start, end int) error {
	srcPath := sandRoot + "/conf/"
	for i := start; i <= end; i++ {
		dstNodePath := sandRoot + "/nodes/node" + strconv.Itoa(i)
		if err := copyFile(srcPath+"tmp/xchain.yaml", dstNodePath+"/conf/xchain.yaml"); err != nil {
			return err
		}
	}
	return nil
}

func copyChainConf(start, end int) error {
	srcPath := sandRoot + "/conf/"
	for i := start; i <= end; i++ {
		dstNodePath := sandRoot + "/nodes/node" + strconv.Itoa(i)
		if err := copyFile(srcPath+"tmp/xuper.json", dstNodePath+"/data/config/xuper.json"); err != nil {
			return err
		}
	}
	return nil
}

func copyPluginConf(start, end int) error {
	srcPath := sandRoot + "/conf/"
	for i := start; i <= end; i++ {
		dstNodePath := sandRoot + "/nodes/node" + strconv.Itoa(i)
		if err := copyFile(srcPath+"plugins.conf.tpl", dstNodePath+"/conf/plugins.conf"); err != nil {
			return err
		}
	}
	return nil
}

func getNodeNetURL(nodeID int) string {
	cmdStr := sandRoot + "/bin/xchain-cli netURL preview --path " + sandRoot + "/nodes/node" + strconv.Itoa(nodeID) + "/data/netkeys/"
	netKey, _ := exec.Command("bash", "-c", cmdStr).Output()
	netURL := strings.Replace(string(netKey), "\n", "", -1)
	netURL = strings.Replace(netURL, "ip4", "dns4", -1)
	netURL = strings.Replace(netURL, "127.0.0.1", "node"+strconv.Itoa(nodeID), -1)
	return netURL
}

func renderNodeConf() error {
	bootURL := getNodeNetURL(1)
	nodeConfTplPath := sandRoot + "/conf/xchain.yaml.tpl"
	os.Mkdir(sandRoot+"/conf/tmp", 0755)
	nodeConfTmpPath := sandRoot + "/conf/tmp/xchain.yaml"
	chainConfTpl := map[string]interface{}{
		"SeedUrl": bootURL,
	}
	return renderFile(nodeConfTplPath, nodeConfTmpPath, chainConfTpl)
}

func getNodeAddress(nodeID int) string {
	path := sandRoot + "/nodes/node" + strconv.Itoa(nodeID) + "/data/keys/address"
	address, _ := ioutil.ReadFile(path)
	return string(address)
}

func renderChainConf() error {
	addrList := []string{}
	netURLList := []string{}
	for i := 1; i <= int(minerNumber); i++ {
		addr := getNodeAddress(i)
		addrList = append(addrList, addr)
		netURL := getNodeNetURL(i)
		netURLList = append(netURLList, netURL)
	}

	InitProposer := ""
	for i := 0; i <= len(addrList)-1; i++ {
		if i != len(addrList)-1 {
			InitProposer = InitProposer + "\"" + addrList[i] + "\",\n"
		} else {
			InitProposer = InitProposer + "\"" + addrList[i] + "\""
		}
	}

	InitProposerNeturl := ""
	for i := 0; i <= len(netURLList)-1; i++ {
		if i != len(netURLList)-1 {
			InitProposerNeturl = InitProposerNeturl + "\"" + netURLList[i] + "\",\n"
		} else {
			InitProposerNeturl = InitProposerNeturl + "\"" + netURLList[i] + "\""
		}
	}

	chainConfTpl := map[string]interface{}{
		"PredistributionAddr": addrList[0],
		"ProposerNum":         strconv.Itoa(int(minerNumber)),
		"InitProposer":        InitProposer,
		"InitProposerNeturl":  InitProposerNeturl,
	}
	chainConfTplPath := sandRoot + "/conf/xuper.json.tpl"
	chainConfTmpPath := sandRoot + "/conf/tmp/xuper.json"
	return renderFile(chainConfTplPath, chainConfTmpPath, chainConfTpl)
}

func renderFile(tplFile, tmpFile string, args interface{}) error {
	t, err := template.ParseFiles(tplFile)
	if err != nil {
		return err
	}

	f, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, args)
}

func initSetConf() error {
	if err := copySeedNodeConf(); err != nil {
		return err
	}
	if err := renderNodeConf(); err != nil {
		return err
	}
	if err := renderChainConf(); err != nil {
		return err
	}
	if err := copyNodeConf(2, int(nodeNumber)); err != nil {
		return err
	}
	if err := copyChainConf(1, int(nodeNumber)); err != nil {
		return err
	}
	if err := copyPluginConf(1, int(nodeNumber)); err != nil {
		return err
	}
	return nil
}

func initConf() error {
	if withoutInitConf {
		if err := copyNodeConf(1, int(nodeNumber)); err != nil {
			return err
		}
		if err := copyChainConf(1, int(nodeNumber)); err != nil {
			return err
		}
		if err := copyPluginConf(1, int(nodeNumber)); err != nil {
			return err
		}
	}
	return initSetConf()
}

func genSlice(num int) []int {
	res := []int{}
	for i := 1; i <= num; i++ {
		res = append(res, i)
	}
	return res
}

func renderComposeFile() error {
	composeConfTplPath := sandRoot + "/conf/docker-compose.yml.tpl"
	composeConfTmpPath := sandRoot + "/docker-compose.yml"
	composeFileTpl := map[string]interface{}{
		"Index":    genSlice(int(nodeNumber)),
		"SandRoot": sandRoot,
		"User":     userCurr,
		"PortSeg":  int(10),
	}
	return renderFile(composeConfTplPath, composeConfTmpPath, composeFileTpl)
}

func initDockerCompose() error {
	return renderComposeFile()
}
