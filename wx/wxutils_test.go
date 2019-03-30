package wx

import (
	"fmt"
	"testing"
)

func TestOpen(t *testing.T) {
	err := Open("C://Users//Administrator//Desktop//repo//pub//v1.0//yegoo-marking-mp")
	// err := OpenToolsOrProject("     ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("success")

}

func TestLogin(t *testing.T) {
	qrCode, err := Login()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("success: %s\n", qrCode)
}

func TestPreview(t *testing.T) {
	qrCode, err := Preview("C://Users//Administrator//Desktop//yegoo-marking-mp", "pages/index/index")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("success: %s\n", qrCode)
}

func TestUpload(t *testing.T) {
	err := Upload("C://Users//Administrator//Desktop//yegoo-marking-mp", "v1.0", "", "")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestClose(t *testing.T) {
	err := Close("C://Users//Administrator//Desktop//repo//pri//10000//v1.0//yegoo-marking-mp")
	if err != nil {
		fmt.Println(err)
	}
	// 实测失败 返回200, ide unknown command (empty parse result): /close?projectpath=xxx
}

func TestQuit(t *testing.T) {
	err := Quit()
	if err != nil {
		fmt.Println(err)
	}
}

func TestBuildnpm(t *testing.T) {
	err := Buildnpm("C://Users//Administrator//Desktop//repo//pri//10000//v1.0//yegoo-marking-mp")
	if err != nil {
		fmt.Println(err)
	}
}

func TestAutoTest(t *testing.T) {
	result, err := AutoTest("C://Users//Administrator//Desktop//repo//pri//10000//v1.0//yegoo-marking-mp")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestAutoPreview(t *testing.T) {
	err := AutoPreview("C://Users//Administrator//Desktop//repo//pri//10000//v1.0//yegoo-marking-mp")
	if err != nil {
		fmt.Println(err)
	}
}
