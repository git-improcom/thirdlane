/*
 * Copyright 2019. Improcom Inc
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 *  IMPROCOM, Inc
 *  1803 Gravesend Neck road
 *  Brooklyn, NY, 11229 USA
 *  Phone: +1-718-769-3400
 *  @author Alex Romanchuk - alexr@improcom.com - Improcom INC
 */

package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/improcom/thirdlane/thdbutil"
	"github.com/improcom/thirdlane/yealink"
	"github.com/spf13/viper"
	"math/rand"
	"os"
	"strconv"
)

type PhoneDirectoryEntry struct {
	ContactName   string
	ContactPhone1 string
	ContactPhone2 string
	ContactPhone3 string
}

var (
	vipe, vipe_dev *viper.Viper
	dev_env        bool
)

func init() {
	vipe = viper.New()
	vipe_dev = viper.New()
	vipe.SetConfigName("config")
	vipe.SetConfigType("yaml")
	vipe.AddConfigPath(".")
	dev_env, _ = strconv.ParseBool(os.Getenv("DEV_ENVIRONMENT"))
	if !dev_env {
		vipe.AddConfigPath("/usr/local/utils/thirdlane_modules/thirdlane_directory_generator")
	}
	err := vipe.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if dev_env {
		vipe_dev.SetConfigName(".env")
		vipe_dev.SetConfigType("yaml")
		vipe_dev.AddConfigPath(".")
		err := vipe_dev.ReadInConfig()
		if err != nil {
			panic(err)
		}
	}

}

func main() {
	tenant := vipe.GetString("tenant")
	dirEntries, _ := thdbutil.GetDirectoryEntries(tenant)
	yealinkProcessing(dirEntries, tenant)

	//spew.Dump(dirEntries)
}

func yealinkProcessing(dirEntries []thdbutil.DirectoryEntry, tenant string) {
	yPhoneConfig := mysql2yealink(dirEntries)
	xmlContent := yealink.GenerateXml(yPhoneConfig)

	myDomain := vipe.GetString("mydomain")
	fileName := createFileName(tenant)
	dirPath := "/home/PlcmSpIp/CONTACTS/yealink_directory/"
	if dev_env {
		dirPath = ""
	}
	fullPath := fmt.Sprintf("%v%v.xml", dirPath, fileName)
	f, err := os.Create(fullPath)

	if err != nil {
		panic(err)
	}
	defer f.Close() // make sure it gets closed after
	fmt.Fprintf(f, xmlContent.String())
	fmt.Printf("remote_phonebook.data.1.url = http://%v/provisioning/CONTACTS/yealink_directory/%v.xml\n", myDomain, fileName)
}

func mysql2yealink(dirEntries []thdbutil.DirectoryEntry) (yPhoneConfig yealink.PhoneDirectoryConfig) {
	for _, phEntry := range dirEntries {
		var yealinkContact yealink.ConfigDirectoryEntry
		yealinkContact.Name = fmt.Sprintf("%v %v", phEntry.Firstname, phEntry.Lastname)
		yealinkContact.Phone1 = phEntry.Extension
		yealinkContact.Phone2 = phEntry.Phone_mobile
		yealinkContact.Phone3 = phEntry.Phone_home
		yPhoneConfig.Entries = append(yPhoneConfig.Entries, yealinkContact)
	}
	return
}

func createFileName(tenant string) string {
	dict := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	salt := make([]rune, 33)
	for i := range salt {
		salt[i] = dict[rand.Intn(len(dict))]
	}
	return fmt.Sprintf("%v_%v", tenant, string(salt))
}
