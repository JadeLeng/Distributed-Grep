package grepRPC

import (
	"io/ioutil"
	"os/exec"
)

// pattern - grep pattern, file - search file range
type GrepArgs struct {
	Pattern, File string
}

// GreoRes...the return result of grep
type GrepRes string

func (t *GrepRes) GetGrep(args *GrepArgs, reply *string) error {
	cmd := "grep " + args.Pattern + " " + args.File
	grepCmd := exec.Command("bash", "-c", cmd)
	grepIn, _ := grepCmd.StdinPipe()
	grepOut, _ := grepCmd.StdoutPipe()
	grepCmd.Start()
	grepIn.Write([]byte("hello grep\ngoodbye grep"))
	grepIn.Close()
	grepBytes, _ := ioutil.ReadAll(grepOut)
	grepCmd.Wait()

	//fmt.Println("> grep argu")
	//fmt.Println(string(grepBytes))
	*reply = string(grepBytes)
	return nil
}
