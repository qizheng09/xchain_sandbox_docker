package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

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
		}
		fmt.Println("Initialize successful!")
	},
}

var (
	nodeNumber   int64
	minerNumber  int64
	withInitConf bool
	xRoot        string
	sandRoot     string
)

func initFlags() {
	xRoot = os.Getenv("XCHAIN_ROOT")
	sandRoot = os.Getenv("XCHAIN_SAND_ROOT")
	initCmd.Flags().Int64VarP(&nodeNumber, "nodeNumber", "N", 5, "The number of nodes to start")
	initCmd.Flags().Int64VarP(&minerNumber, "minerNumber", "M", 3, "The number of nodes to start")
	initCmd.Flags().BoolVarP(&withInitConf, "withInitConf", "", false, "The flag whether to init `xchain.yaml` and `xuper.json`")
}

// initEnv will init the multi xchain sandbox
// update binary
// init NodesFiles
// init configs: xchain.yaml, xuper.json
// init docker-compose.yml
func initEnv() error {
	if xRoot == "" {
		return errors.New("The XCHAIN_ROOT environment variable have not been set!")
	}
	if sandRoot == "" {
		return errors.New("The XCHAIN_SAND_ROOT environment variable have not been set!")
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
	// update plugins
	sorcPluginsPath := xRoot + "/output/plugins"
	destPluginsPath := sandRoot + "/plugins/"
	fs, _ := ioutil.ReadDir(sorcPluginsPath)
	for _, v := range fs {
		if v.IsDir() {
			fss, _ := ioutil.ReadDir(sorcPluginsPath + "/" + v.Name())
			for _, vs := range fss {
				srcFile := sorcPluginsPath + "/" + v.Name() + "/" + vs.Name()
				dstFile := destPluginsPath + v.Name() + "/" + vs.Name()
				if err := copyFile(srcFile, dstFile); err != nil {
					return err
				}
			}
		} else {
			srcFile := sorcPluginsPath + "/" + v.Name()
			dstFile := destPluginsPath + v.Name()
			if err := copyFile(srcFile, dstFile); err != nil {
				return err
			}
		}
	}
	return nil
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

// initNodesFiles will init nodes date/keys, data/netKeys and create filefolder needed
func initNodesFiles(nodeNums int64) error {
	if nodeNums < 0 {
		return errors.New("Node number can not less than 0!")
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

func initConf() error {
	if !withInitConf {
		srcPath := sandRoot + "/conf/"
		for i := 1; i <= int(nodeNumber); i++ {
			dstNodePath := sandRoot + "/nodes/node" + strconv.Itoa(i)
			if err := copyFile(srcPath+"xchain.yaml.tpl", dstNodePath+"/conf/xchain.yaml"); err != nil {
				return err
			}
			if err := copyFile(srcPath+"plugins.conf.tpl", dstNodePath+"/conf/plugins.conf"); err != nil {
				return err
			}
			if err := copyFile(srcPath+"xuper.json.tpl", dstNodePath+"/data/config/xuper.json"); err != nil {
				return err
			}
		}
	}
	return initSetConf()
}

func initSetConf() error {

	return nil
}

func initDockerCompose() error {
	return nil
}
