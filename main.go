package main

import (
	"encoding/json"
	"github.com/tidusant/c3m-common/c3mcommon"
	"github.com/tidusant/c3m-common/mycrypto"
	"github.com/tidusant/chadmin-repo/models"
	"log"

	"fmt"
	"os"
)

var testsession string

func decodeResponse(querystring string, data string) (rs models.RequestResult, err error) {
	//encode data
	rtstr := c3mcommon.RequestAPI(os.Getenv("API_URL"), querystring, data)
	err = json.Unmarshal([]byte(rtstr), &rs)
	if err != nil {
		rs.Error = rtstr
	}
	return
}

func doCall(testname, requesturl, queryData string) models.RequestResult {
	fmt.Println("\n\n==== " + testname + " ====")
	fmt.Printf("Data: url: %s - data:%s\n", requesturl, queryData)
	rs, err := decodeResponse(requesturl, queryData)
	if err != nil {
		log.Fatalf("Test fail: request error: %s", err.Error())
	}
	fmt.Printf("Request return: %+v\n", rs)
	return rs
}

func setup() {
	// Switch to test mode so you don't get such noisy output

}
func main() {
	setup()
	TestSpecialChar()
	TestCreateSex()
	TestCreateSex2()
	TestLoginWithouSession()
	TestLoginWrongUser()
	TestLoginSuccessUser()
	TestCallRPCWithoutSession()
	TestCallRPCWithoutAuth()
	TestCallUnkownRPCWithAuth()
	TestCallRPCWithUnknownAction()
	TestCallRPCWithAuth()
	fmt.Println("===== ALL PASSED =======")
}

//test special char
func TestSpecialChar() {
	rs := doCall("TestSpecialChar", c3mcommon.GetSpecialChar(), "")
	//check test data
	if rs.Status == 1 {
		log.Fatalf("Test fail")
	}
	fmt.Println("PASS")
}

func TestCreateSex() {
	rs := doCall("TestCreateSex", "CreateSex", "")
	//check test data
	if rs.Status != 1 {
		log.Fatalf("Cannot create Sex")
	}
	fmt.Println("PASS")
}

//double test create session
func TestCreateSex2() {
	rs := doCall("TestCreateSex2", "CreateSex", "")
	//check test data
	if rs.Status != 1 {
		log.Fatalf("Cannot create Sex after create sex")
	}
	testsession = rs.Data
	fmt.Println("PASS")
}

//test login
func TestLoginWithouSession() {
	rs := doCall("TestLoginWithouSession", "aut", "l|admin,123456")
	//check test data
	if rs.Status == 1 {
		log.Fatalf("Test fail: user logged in without session")
	}
	fmt.Println("PASS")

}
func TestLoginWrongUser() {
	rs := doCall("TestLoginWrongUser", "aut", mycrypto.EncDat2(testsession)+"|l|admin,123456")
	//check test data
	if rs.Status == 1 {
		log.Fatalf("Test fail: user logged in with wrong username pass")
	}
	fmt.Println("PASS")
}
func TestLoginSuccessUser() {
	rs := doCall("TestLoginSuccessUser", "aut", mycrypto.EncDat2(testsession)+"|l|demo,123")
	//check test data
	if rs.Status != 1 {
		log.Fatalf("Test fail: user cannot login with session and userpass")
	}
	fmt.Println("PASS")
}
func TestCallRPCWithoutSession() {
	rs := doCall("TestCallRPCWithoutSession", "shop", "lsi|abc,123")
	//check test data
	if rs.Status == 1 {
		log.Fatalf("Test fail: user can call rpc without session")
	}
	fmt.Println("PASS")
}
func TestCallRPCWithoutAuth() {
	rs := doCall("TestCallRPCWithoutAuth", "shop", "notlogginsession|lsi|abc,123")
	//check test data
	if rs.Status == 1 {
		log.Fatalf("Test fail: user can call rpc without auth (notloggedinsession)")
	}
	fmt.Println("PASS")
}
func TestCallUnkownRPCWithAuth() {
	rs := doCall("TestCallUnkownRPCWithAuth", "unknownauth", mycrypto.EncDat2(testsession)+"|lsi|abc,123")
	//check test data
	if rs.Status == 1 || rs.Error != "service not run" {
		log.Fatalf("Test fail: user can call unknow rpc ")
	}
	fmt.Println("PASS")
}
func TestCallRPCWithUnknownAction() {
	rs := doCall("TestCallRPCWithUnknownAction", "shop", mycrypto.EncDat2(testsession)+"|unknowaction|abc,123")
	//check test data
	if rs.Status == 1 || rs.Error != "Hello admin-portal" {
		log.Fatalf("Test fail: user can call unknow action ")
	}
	fmt.Println("PASS")
}
func TestCallRPCWithAuth() {
	rs := doCall("TestCallRPCWithAuth", "shop", mycrypto.EncDat2(testsession)+"|lsi|")
	//check test data
	if rs.Status != 1 {
		log.Fatalf("Test fail: user can not call rpc with action properly")
	}
	fmt.Println("PASS")
}
