package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/skratchdot/open-golang/open"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
)

const (
	windowWidth  = 650
	windowHeight = 480
	title        = "InstaCrypt"
	version      = "v0.1.0"
	web          = "https://instacrypt.io"
	api          = "https://api.instacrypt.io"
)

// Vars injected via ldflags by bundler
var (
	AppName            string
	BuiltAt            string
	VersionAstilectron string
	VersionElectron    string
)

// Application Vars
var (
	fs    = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	debug = fs.Bool("d", false, "enables the debug mode")
	w     *astilectron.Window
)

// User Environment
var uname string     // User Name
var homedir string   // Home Directory
var lsshdir string   // Local SSH Dir
var sshdir string    // Secrets Desktop SSH Dir
var configdir string // Database Dir
var tmpdir string    // Temp Dir for Processing Encrypt/Decrypt

// Authorized
var authed = false
var token = ""

// User Input
var file string
var key string
var direction = "encrypt"
var count = 0

// Key and Properties
type KeyProp struct {
	Direction string `json:"dir"`
	Key       string `json:"key"`
}

type FileType struct {
	Type string `json:"type"`
	File string `json:"file"`
}

// Locker DB Structs
type lockstruct struct {
	Name   string
	Locker string
}

type lcstruct struct {
	Contacts string
}

type access struct {
	Token string
}

type Contact struct {
	Id    uint
	Name  string
	Email string
	Lock  string
}

type Response struct {
	Data    []*Contact
	Message string
	Status  bool
}

type Message struct {
	Ok      bool
	Title   string
	Message string
}

type CList struct {
	Ok      bool
	Title   string
	Message string
	List    []string
}

// Get/Set User Configurations
func getUserEnv() {
	usr, err := user.Current() // Get Current User
	if err != nil {
		fmt.Print("Get User Error: ")
		fmt.Println(err)
	}
	uname = usr.Username
	homedir = usr.HomeDir
	lsshdir = homedir + "/.ssh"
	sshdir = homedir + "/.icfx/.ssh"
	configdir = homedir + "/.icfx"
	tmpdir = homedir + "/.icfx/tmp"
}

func checkAuthed() {

	// Instantiate DB
	lockers, err := scribble.New(configdir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	// Get Locker Count
	a := access{}
	if err := lockers.Read(".auth", "access", &a); err != nil {
		fmt.Println("ACCESS: DB READ Error: ", err)
	}

	token = a.Token
	if token != "" {
		authed = true
	}
}

func auth(secret string) bool {

	// Instantiate DB
	lockers, err := scribble.New(configdir, nil)
	if err != nil {
		fmt.Println("Error: ", err)
		return false
	}

	rerr := removeAuth()
	if rerr != nil {
		fmt.Print("Error: ", rerr)
		return false
	}

	// Insert app access secret
	a := &access{}
	a.Token = secret
	aerr := lockers.Write(".auth", "access", a)
	if aerr != nil {
		fmt.Println("Error: ", aerr)
		return false
	}
	fmt.Println("Successfully written secret to local storage...")
	checkAuthed()
	return true

}

func removeAuth() error {

	// Remove backup folder if it exists
	_, serr := os.Stat(configdir + "/.auth.old")
	if serr == nil {
		serr := os.RemoveAll(configdir + "/.auth.old")
		if serr != nil {
			return serr
		}
	}

	// Backup current auth folder if it exists
	_, aerr := os.Stat(configdir + "/.auth")
	if aerr == nil {
		rerr := os.Rename(configdir+"/.auth", configdir+"/.auth.old")
		if rerr != nil {
			return rerr
		}
	}

	// Create auth folder
	merr := os.Mkdir(configdir+"/.auth", 0700)
	if merr != nil {
		return merr
	}

	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func keyExists() bool {
	checkKey, _ := exists(sshdir + "/ic-lock")
	if checkKey == true {
		return true
	} else {
		return false
	}
}

func lKeyExists() bool {
	checkKey, _ := exists(lsshdir + "/id_rsa")
	if checkKey == true {
		return true
	} else {
		return false
	}
}

func keygen() (bool, string, string) {
	checkDir, _ := exists(sshdir)
	if checkDir == false {
		merr := os.MkdirAll(sshdir, 0700)
		if merr != nil {
			log.Fatal(merr)
			return false, "ERROR", fmt.Sprint(merr)
		}
	}
	checkKey, _ := exists(sshdir + "/ic-lock")
	if checkKey == true {
		fmt.Println("A key already exists...")
		return false, "ERROR", "A key already exists..."
	}
	cmd := exec.Command("sh", "-c", "ssh-keygen -t rsa -b 4096 -f "+sshdir+"/ic-lock -N ''")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "ssh-keygen -t rsa -b 4096 -f "+sshdir+"/ic-lock -N ''")
	}
	err := cmd.Run()
	if err != nil {
		fmt.Println("KeyGen Error:", err.Error())
		return false, "ERROR", err.Error()
	}

	// Check to make sure public key file is there
	if _, err := os.Stat(sshdir + "/ic-lock.pub"); os.IsNotExist(err) {
		fmt.Println("Looks like the lock and key didn't get generated properly...")
		return false, "ERROR", "Looks like the lock and key didn't get generated properly..."
	}

	// Read new public key into string
	content, err := ioutil.ReadFile(sshdir + "/ic-lock.pub")
	if err != nil {
		log.Fatal(err)
	}
	// Convert []byte to string and print to screen
	loc := string(content)
	fmt.Println(loc)
	status, _ := addLockToServer(loc)
	if status == false {
		fmt.Println("Could not register your lock with InstaCrypt...")
		return false, "ERROR", "Could not register your lock with InstaCrypt..."
	}
	return pop("keygened")
}

func keydel() (bool, string, string) {
	checkKey, _ := exists(sshdir + "/ic-lock")
	if checkKey == true {
		status, _ := removeLockFromServer()
		if status == false {
			fmt.Println("Could not remove lock from InstaCrypt Vault...")
			return false, "ERROR", "Could not remove lock from InstaCrypt Vault..."
		}
		err := os.RemoveAll(sshdir + "/")
		if err != nil {
			fmt.Println("Error occurred when attempting to delete key pair...")
			return false, "ERROR", "Error occurred when attempting to delete key pair..."
		}
		return true, "Deleted", "Your Lock & Key has been Deleted from your device and InstaCrypt Vault!"

	} else {
		fmt.Println("Lock & Key does not exist!")
		return false, "ERROR", "Lock & Key does not exist!"
	}
}

func handleRPC(data string) {
	switch {
	case data == "openosi":
		openl("https://osi.3df.io")
	case data == "open3df":
		openl("https://3df.io")
	case data == "openic":
		openl("https://instacrypt.io")
	}
}

//Helper for Opening URL with Default Browser
func openl(uri string) {
	err := open.Run(uri)
	if err != nil {
		fmt.Println(err)
	}
}

func showFile() (bool, string) {

	var fname string

	if direction == "encrypt" {
		fname = file + ".icfx"
	}

	if direction == "decrypt" {
		fname = removeExt(file)
	}

	if runtime.GOOS == "linux" {
		err := open.StartWith(fname, "nautilus")
		if err != nil {
			return false, "Couldn't open folder: " + fmt.Sprint(err)
		}
	}

	if runtime.GOOS == "darwin" {
		cmd := exec.Command("open", "-R", fname)
		err := cmd.Run()
		if err != nil {
			return false, "Couldn't open folder: " + fmt.Sprint(err)
		}
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command(`cmd`, `/C`, `explorer`, `/select,`, fname)
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
			// This throws an error even though it successfully opens the file
			// pop("\nCouldn't open folder " + fname + ": " + fmt.Sprint(err))
		}
	}
	return true, ""

}

func keygened() (bool, string) {
	fname := sshdir + "/ic-lock.pub"
	if runtime.GOOS == "linux" {
		err := open.StartWith(fname, "nautilus")
		if err != nil {
			return false, "\nCouldn't open folder: " + fmt.Sprint(err)
		}
	}

	if runtime.GOOS == "darwin" {
		cmd := exec.Command("open", "-R", fname)
		err := cmd.Run()
		if err != nil {
			return false, "\nCouldn't open folder: " + fmt.Sprint(err)
		}
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command(`cmd`, `/C`, `explorer`, `/select,`, fname)
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
			// This throws an error even though it successfully opens the file
			// pop("\nCouldn't open folder " + fname + ": " + fmt.Sprint(err))
		}
	}
	return true, ""
}

func pop(msg string) (bool, string, string) {
	if msg == "done" {
		message := filepath.Base(file) + " has been " + direction + "ed successfully!"
		action := direction + "ed!"
		fmt.Println(message)
		return true, action, message + "\n\nWould you like to show the file in its folder?"
	} else if msg == "keygened" {
		return true, "Lock & Key Generated!", "\nYour Lock and Key has been generated and registered to your InstaCrypt Vault!\n\nWould you like to show the files so you can back them up?"
	} else if msg == "refreshed" {
		return true, "Refreshed!", "\nYour contact list has been refreshed!"
	}
	return false, "ERROR", "pop function error!"
}

func getCurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func getFileName(f string) string {
	var basename string
	if f == "file" {
		basename = filepath.Base(file)
	} else {
		basename = filepath.Base(key)
	}
	return basename
}

func getContactCount() int {
	// Instantiate DB
	lockers, err := scribble.New(configdir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	// Get Locker Count
	locker := lcstruct{}
	if err := lockers.Read("lockers", "count", &locker); err != nil {
		fmt.Println("COUNT: DB READ Error: ", err)
	}
	count, cerr := strconv.Atoi(locker.Contacts)
	if cerr != nil {
		fmt.Println("COUNT: String Conversion Error: ", cerr)
	} else if count == 0 {
		fmt.Println("No Contacts?")
	}

	return count
}

func getContactList() (bool, []string) {
	// Instantiate DB
	lockers, err := scribble.New(configdir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	// Get Locker Count
	locker := lcstruct{}
	if err := lockers.Read("lockers", "count", &locker); err != nil {
		fmt.Println("GET: DB READ Error: ", err)
		return false, nil
	}
	count, cerr := strconv.Atoi(locker.Contacts)
	if cerr != nil {
		fmt.Println("GET: String Conversion Error", cerr)
		return false, nil
	} else if count == 0 {
		fmt.Println("No Contacts?")
	}

	// Feed contact list into array
	var clist = make([]string, count)
	for i := 0; i < count; i++ {
		centry := lockstruct{}
		if err := lockers.Read("lockers", strconv.Itoa(i), &centry); err != nil {
			fmt.Println("GET: Loop Error: ", err)
		}
		clist[i] = centry.Name
	}
	return true, clist
}

func setContactSelected(conid int) (bool, string) {
	// Instantiate DB
	lockers, err := scribble.New(configdir, nil)
	if err != nil {
		fmt.Println("Error", err)
		return false, err.Error()
	}

	// Get Public Key
	con := lockstruct{}
	if err := lockers.Read("lockers", strconv.Itoa(conid), &con); err != nil {
		fmt.Println("GET: Loop Error: ", err)
		return false, err.Error()
	}

	if con.Locker == "" {
		fmt.Println("It looks like this user hasn't registered a lock yet. Contact them and ask them to do so and then try again.")
		return false, "It looks like this user hasn't registered a lock yet. Contact them and ask them to do so and then try again."
	}
	lock := con.Locker

	// fmt.Println("Using the Locker of: " + con.Name)

	// Check if Temp Directory Exist
	checkDir, _ := exists(tmpdir)
	if checkDir == false {
		merr := os.MkdirAll(tmpdir, 0700)
		if merr != nil {
			log.Fatal(merr)
			return false, merr.Error()
		}
	}

	// Check if a Temp Lock is Already There, If So, Delete
	checkLock, _ := exists(tmpdir + "/lock.pub")
	if checkLock == true {
		err := os.Remove(tmpdir + "/lock.pub")
		if err != nil {
			fmt.Println("Couldn't Delete Temp Lock File...")
			return false, "Couldn't Delete Temp Lock File..."
		}
	}

	// Create Key File
	lfile, err := os.Create(tmpdir + "/lock.pub")
	if err != nil {
		log.Fatal(err)
		return false, err.Error()
	}
	writer := bufio.NewWriter(lfile)
	_, werr := writer.WriteString(lock)
	if werr != nil {
		fmt.Println("Got error while writing to a file. Err: %s", err.Error())
		return false, err.Error()
	}
	writer.Flush()

	// Set the Temp Public Key to Use
	key = tmpdir + "/lock.pub"
	// fmt.Println(key)

	return true, ""
}

func removeExt(input string) string {
	if len(input) > 0 {
		if i := strings.LastIndex(input, "."); i > 0 {
			input = input[:i]
		}
	}
	return input
}

func checkLogic() (bool, string) {
	// check that the user specified a file.
	if file == "" {
		return false, "\nPlease select a file to " + direction + "."
	}

	// Check that the user specified a key
	if key == "" {
		if direction == "encrypt" {
			return false, "\nPlease select a public key to " + direction + "."
		} else {
			return false, "\nPlease select a private key to " + direction + "."
		}
	}

	// Check that the file does not have a "\"
	if strings.Contains(getFileName("file"), "\\") {
		return false, "File name must not contain a \"\\\"... Please change the file name and try again."
	}

	// Check that the key does not have a "`"
	if strings.Contains(file, "`") {
		return false, "File name must not contain a \"`\"... Please change the file name and try again."
	}

	// Check that the user specified file has a .ssh extension if encrypt is selected
	if direction == "encrypt" && strings.Contains(file, ".icfx") {
		return false, "You are trying to encrypt with a lock. You probably want to decrypt instead?"
	}

	// Check that the user specified file has a .ssh extension if decrypt is selected
	if direction == "decrypt" && strings.Contains(file, ".icfx") == false {
		return false, "You can only decrypt .icfx files..."
	}

	// Read key file to check
	k, err := ioutil.ReadFile(key)
	if err != nil {
		return false, "\nCan't read Lock or Key: " + err.Error()
	}
	ks := string(k)

	// check that the user specified key is not a private key if encrypt is selected
	if direction == "encrypt" && strings.Contains(ks, "PRIVATE") == true {
		return false, "\nTo encrypt a file, you should be using the lock of the receiver instead. Please try again."
	}

	// check that the user specified key is not a public key if decrypt is selected
	if direction == "decrypt" && strings.Contains(ks, "PRIVATE") == false {
		return false, "\nTo decrypt a file, you should be using your key instead. Please try again."
	}

	return true, ""
}

func process() (bool, string, string) {

	check, msg := checkLogic()
	if check == false {
		// Logic check did not pass. Don't continue...
		return false, "ERROR", msg
	}

	if direction == "encrypt" {
		if runtime.GOOS == "linux" {
			cmd := exec.Command("bash", "-c", "ssh-vault -k \""+key+"\" create < \""+file+"\" \""+file+".icfx\"")
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			err := cmd.Run()
			// fmt.Println(cmd)
			if err != nil {
				return false, "ERROR", fmt.Sprint(err) + "\n\n" + stderr.String()
			} else {
				return pop("done")
			}
		}
		if runtime.GOOS == "darwin" {
			dir := getCurrentPath()
			cmd := exec.Command("sh", "-c", dir+"/ssh-vault -k \""+key+"\" create < \""+file+"\" \""+file+".icfx\"")
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			cerr := cmd.Run()
			// fmt.Println(cmd)
			if cerr != nil {
				return false, "ERROR", fmt.Sprint(cerr) + "\n\n" + stderr.String()
			} else {
				return pop("done")
			}
		}
		if runtime.GOOS == "windows" {
			dir := getCurrentPath()
			k := filepath.ToSlash(key)
			f := filepath.ToSlash(file)
			cmd := exec.Command("cmd", "/C", dir+"\\ssh-vault.exe", "-k", k, "create", "<", f, f+".icfx")
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			cerr := cmd.Run()
			if cerr != nil {
				return false, "ERROR", fmt.Sprint(cmd) + "\n" + fmt.Sprint(cerr) + "\n\n" + stderr.String()
			} else {
				return pop("done")
			}
		}
	}
	if direction == "decrypt" {
		if runtime.GOOS == "linux" {
			cmd := exec.Command("bash", "-c", "ssh-vault -k \""+key+"\" -o \""+removeExt(file)+"\" view \""+file+"\"")
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			err := cmd.Run()
			// fmt.Println(cmd)
			if err != nil {
				return false, "ERROR", fmt.Sprint(err) + "\n\n" + stderr.String()
			} else {
				return pop("done")
			}
		}
		if runtime.GOOS == "darwin" {
			dir := getCurrentPath()
			cmd := exec.Command("sh", "-c", dir+"/ssh-vault -k \""+key+"\" -o \""+removeExt(file)+"\" view \""+file+"\"")
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			err := cmd.Run()
			// fmt.Println(cmd)
			if err != nil {
				return false, "ERROR", fmt.Sprint(err) + "\n\n" + stderr.String()
			} else {
				return pop("done")
			}
		}
		if runtime.GOOS == "windows" {
			dir := getCurrentPath()
			cmd := exec.Command("cmd", "/C", dir+"\\ssh-vault", "-k", key, "-o", removeExt(file), "view", file)
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			cerr := cmd.Run()
			// fmt.Println(cmd)
			if cerr != nil {
				return false, "ERROR", fmt.Sprint(cerr) + "\n\n" + stderr.String()
			} else {
				return pop("done")
			}
		}
	}
	return false, "ERROR", "Something went wrong..."
}

func getServerContacts() (bool, string) {

	fmt.Println("Attempting to retrieve contact list from server")

	// Setup HTTP Request
	request, _ := http.NewRequest("GET", api+"/v1/app/contacts/get", nil)
	fmt.Println("Token: " + token)
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", web)

	// Call API
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return false, "\nServer unreachable... Please check your internet connection and then try again."
	}

	// Read Response
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))

	// Ouptut Response to STDOUT
	var response Response
	json.Unmarshal([]byte(string(data)), &response)
	fmt.Println("Status: ", response.Status)
	fmt.Println("Message: ", response.Message)
	fmt.Println("Data: ", response.Data)

	if response.Status == false {
		return false, "\nThere's something wrong with your secret. Generate a new one and re-authorize your app to try again."
	}

	for i := 0; i < len(response.Data); i++ {
		fmt.Print(i)
		fmt.Println(" Name : " + response.Data[i].Name)
	}

	// Remove Contact List
	rerr := removeContactList()
	if rerr != nil {
		fmt.Println("Couldn't clear contact list")
		fmt.Println(rerr)
		return false, "\nCould not refresh local list with the server's list..."
	}

	// Instantiate DB to load new contacts
	lockers, lerr := scribble.New(configdir, nil)
	if lerr != nil {
		fmt.Println("Error", lerr)
		return false, "\nError reading local database..."
	}

	// Load contacts
	for i := 0; i < len(response.Data); i++ {
		c := &lockstruct{}
		c.Name = response.Data[i].Name + " - " + response.Data[i].Email
		c.Locker = response.Data[i].Lock
		cerr := lockers.Write("lockers", strconv.Itoa(i), c)
		if cerr != nil {
			fmt.Println("Error: ", cerr)
			return false, "\nError refreshing local list with the server's list"
		}
	}

	// Insert contact list length
	count := &lcstruct{}
	count.Contacts = strconv.Itoa(len(response.Data))
	counterr := lockers.Write("lockers", "count", count)
	if counterr != nil {
		fmt.Println("Error: ", counterr)
		return false, "\nLocal DB Error..."
	}

	fmt.Println("Finished retrieving server contacts...")
	return true, ""

}

func addLockToServer(lock string) (bool, string) {

	fmt.Println("Attempting to upload lock to InstaCrypt")

	jsonData := map[string]string{"lock": lock}
	jsonValue, _ := json.Marshal(jsonData)

	// Setup HTTP Request
	request, _ := http.NewRequest("POST", api+"/v1/app/lock/add", bytes.NewBuffer(jsonValue))
	fmt.Println("Token: " + token)
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", web)

	// Call API
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return false, ""
	}

	// Read Response
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))

	// Ouptut Response to STDOUT
	var response Response
	json.Unmarshal([]byte(string(data)), &response)
	fmt.Println("Status: ", response.Status)
	fmt.Println("Message: ", response.Message)
	fmt.Println("Data: ", response.Data)

	if response.Status == false {
		return false, "Unable to submit lock to InstaCrypt Vault"
	}

	return true, ""

}

func removeLockFromServer() (bool, string) {

	fmt.Println("Attempting to upload lock to InstaCrypt")

	// Setup HTTP Request
	request, _ := http.NewRequest("POST", api+"/v1/app/lock/remove", nil)
	fmt.Println("Token: " + token)
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", web)

	// Call API
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return false, "Are you connected to the internet?"
	}

	// Read Response
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))

	// Ouptut Response to STDOUT
	var response Response
	json.Unmarshal([]byte(string(data)), &response)
	fmt.Println("Status: ", response.Status)
	fmt.Println("Message: ", response.Message)
	fmt.Println("Data: ", response.Data)

	if response.Status == false {
		return false, "Unable to remove lock from InstaCrypt Vault"
	}

	return true, ""

}

func removeContactList() error {

	// Remove backup folder if it exists
	_, serr := os.Stat(configdir + "/lockers.old")
	if serr == nil {
		serr := os.RemoveAll(configdir + "/lockers.old")
		if serr != nil {
			return serr
		}
	}

	// Backup current contact list if it exists
	_, lerr := os.Stat(configdir + "/lockers")
	if lerr == nil {
		rerr := os.Rename(configdir+"/lockers", configdir+"/lockers.old")
		if rerr != nil {
			return rerr
		}
	}

	// Create configDir
	merr := os.Mkdir(configdir+"/lockers", 0700)
	if merr != nil {
		return merr
	}

	return nil
}

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "sleep":
		time.Sleep(3 * time.Second)
		payload = true
	case "getVersion":
		payload = version
	case "getAuthed":
		payload = authed
	case "authorize":
		var s string
		if err = json.Unmarshal(m.Payload, &s); err != nil {
			payload = err.Error()
			return
		}
		fmt.Println("Auth Check...")
		payload = auth(s)
	case "toConsole":
		var line string
		if err = json.Unmarshal(m.Payload, &line); err != nil {
			payload = err.Error()
			return
		}
		fmt.Println(line)
		payload = "Success!"
	case "handle":
		var data string
		if err = json.Unmarshal(m.Payload, &data); err != nil {
			payload = err.Error()
			return
		}
		handleRPC(data)
		payload = "Success!"
	case "setDirection":
		var d string
		if err = json.Unmarshal(m.Payload, &d); err != nil {
			payload = err.Error()
			return
		}
		direction = d
		payload = d
	case "setKey":
		var res = &KeyProp{}
		if err = json.Unmarshal(m.Payload, res); err != nil {
			payload = err.Error()
			return
		}

		d := res.Direction
		k := res.Key
		fmt.Println("Direction: ", d)
		fmt.Println("Key: ", k)

		msg := &Message{}
		var kext string
		if d == "encrypt" {
			kext = ".pub"
		}
		if d == "decrypt" {
			kext = ""
		}
		if k == "3df" {
			if keyExists() {
				key = sshdir + "/key" + kext
			} else {
				msg.Ok = false
				msg.Title = "ERROR"
				msg.Message = "Lock & Key does not exists! Please generate a new lock & key"
				payload = msg
			}
		} else if k == "local" {
			if lKeyExists() {
				key = lsshdir + "/id_rsa" + kext
			} else {
				msg.Ok = false
				msg.Title = "ERROR"
				msg.Message = "Local private key does not exists! Please generate a new key pair"
				payload = msg
			}
		}
		fmt.Println("Processed Key: ", key)
		msg.Ok = true
		msg.Title = ""
		msg.Message = ""
		payload = msg
	case "setLocalKey":
		var k string
		if err = json.Unmarshal(m.Payload, &k); err != nil {
			payload = err.Error()
		}
		key = k
		name := filepath.Base(key)
		payload = name
	case "setFile":
		var f string
		if err = json.Unmarshal(m.Payload, &f); err != nil {
			payload = err.Error()
			return
		}
		file = f
		name := filepath.Base(file)
		payload = name
	case "clearFile":
		file = ""
		payload = true
	case "clearKey":
		key = ""
		payload = true
	case "getServerContacts":
		msg := &Message{}
		msg.Ok, msg.Message = getServerContacts()
		if msg.Ok == false {
			msg.Title = "ERROR"
		}
		payload = msg
	case "setContactSelected":
		var con int
		if err = json.Unmarshal(m.Payload, &con); err != nil {
			payload = err.Error()
			return
		}
		msg := &Message{}
		msg.Ok, msg.Message = setContactSelected(con)
		if msg.Ok == false {
			msg.Title = "Oops!"
		}
		payload = msg
	case "getContactList":
		msg := &CList{}
		msg.Ok, msg.List = getContactList()
		if msg.Ok == false {
			msg.Title = "ERROR"
			msg.Message = "\nCould not get contact list. Please try again..."
		} else {
			msg.Title = "Refreshed"
			msg.Message = "\nYour contact list has been refreshed!"
		}
		payload = msg
	case "keygen":
		msg := &Message{}
		if authed == false {
			msg.Ok = false
			msg.Title = "Unauthorized"
			msg.Message = "You have to first authorize this app to use this feature."
			payload = msg
			return
		}
		msg.Ok, msg.Title, msg.Message = keygen()
		payload = msg
	case "keygened":
		msg := &Message{}
		msg.Ok, msg.Message = keygened()
		if msg.Ok == false {
			msg.Title = "ERROR"
		}
		payload = msg
	case "keydel":
		msg := &Message{}
		if authed == false {
			msg.Ok = false
			msg.Title = "Unauthorized"
			msg.Message = "You have to first authorize this app to use this feature."
			payload = msg
			return
		}
		msg.Ok, msg.Title, msg.Message = keydel()
		payload = msg
	case "keyExists":
		payload = keyExists()
	case "submit":
		msg := &Message{}
		msg.Ok, msg.Title, msg.Message = process()
		payload = msg
	case "done":
		msg := &Message{}
		msg.Ok, msg.Message = showFile()
		if msg.Ok == false {
			msg.Title = "ERROR"
		}
		payload = msg
	}
	return
}

func main() {

	// Set logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	// Get User Env Virables
	getUserEnv()

	// Check if the app has been authorized
	checkAuthed()

	// Check to see if a 3DF Key Pair Exists
	keyExists()

	// Parse flags
	fs.Parse(os.Args[1:])

	// Run bootstrap
	l.Printf("Running app built at %s\n", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            "InstaCrypt",
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/logo.png",
			SingleInstance:     true,
			VersionAstilectron: VersionAstilectron,
			VersionElectron:    VersionElectron,
		},
		Debug:  false, //true, //*debug,
		Logger: l,
		// Uncomment this block for Mac compiling
		//MenuOptions: []*astilectron.MenuItemOptions{{
		//	Label: astikit.StrPtr(""),
		//	SubMenu: []*astilectron.MenuItemOptions{
		//		{
		//			Label: astikit.StrPtr("Paste"),
		//		},
		//		{Role: astilectron.MenuItemRolePaste},
		//	},
		//}},
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("#333"),
				Center:          astikit.BoolPtr(true),
				Height:          astikit.IntPtr(windowHeight),
				Width:           astikit.IntPtr(windowWidth),
			},
		}},
	}); err != nil {
		l.Fatal(fmt.Errorf("running bootstrap failed: %w", err))
	}

}
