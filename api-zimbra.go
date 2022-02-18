package main

import (
        "fmt"
        "net/http"
        "os"
        "io"
        "time"
        "github.com/gin-gonic/gin"
        "github.com/go-ldap/ldap/v3"
)

const (
        ldapServer   = "ldap://127.0.0.1:389"
        baseDN       = "uid=zimbra,cn=admins,cn=zimbra"
        baseSearch   = "ou=people,dc=example,dc=com"
        
        //CHANGE 
        ldapPassword = "EXECUTE_CODE_BELLOW_TO_RETRIEVE_ZIMBRA_LDAP_PASSWORD"
		
	      // su - zimbra        
        // source ~/bin/zmshutil ; zmsetvars
        // env | grep ldap_root_password | cut -d= -f2
)

func main() {

        gin.DisableConsoleColor()
        gin.SetMode(gin.ReleaseMode)

        accesslogfile, _ := os.Create("/var/log/api-zimbra-access.log")
        errlogfile, _ := os.Create("/var/log/api-zimbra-error.log")
        gin.DefaultWriter = io.MultiWriter(accesslogfile)
        gin.DefaultErrorWriter = io.MultiWriter(errlogfile)

        router := gin.Default()
        router.GET("/getemail/:login", getEmailByID)
  
        //CHANGE
        router.SetTrustedProxies([]string{"CHANGE_FOR_YOUR_ZIMBRA_IP"})
        
        router.Run(":8090")
}

func getEmailByID(c *gin.Context) {

        user := c.Param("login")
        filterDN := fmt.Sprintf("(&(objectClass=zimbraAccount)(uid=%s))", ldap.EscapeFilter(user))
        searchReq := ldap.NewSearchRequest(baseSearch, ldap.ScopeWholeSubtree, 0, 0, 0, false, filterDN, []string{"uid"}, []ldap.Control{})

        errLogger := gin.DefaultErrorWriter

        l, err := ldap.DialURL(ldapServer)
        if err != nil {
                c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "LDAP SERVER ERROR - Connect"})
                t := time.Now()
                error := t.Format("01-02-2006 15:04:05 Mon") + " - " +  "500 - LDAP SERVER ERROR - Connect\n"
                errLogger.Write([]byte(error))
                return
        }

        err = l.Bind(baseDN, ldapPassword)
        if err != nil {
                c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "LDAP SERVER ERROR - BIND"})
                t := time.Now()
                error := t.Format("01-02-2006 15:04:05 Mon") + " - " +  "500 - LDAP SERVER ERROR - BIND\n"
                errLogger.Write([]byte(error))
                return
        }

        result, err := l.Search(searchReq)
        if err != nil {
                c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "LDAP SERVER ERROR - Search"})
                t := time.Now()
                error := t.Format("01-02-2006 15:04:05 Mon") + " - " +  "500 - LDAP SERVER ERROR - Search\n"
                errLogger.Write([]byte(error))
                return
        }

        defer l.Close()

        if len(result.Entries) > 0 {
                c.IndentedJSON(http.StatusOK, gin.H{"message": "TRUE"})
        } else {
                c.IndentedJSON(http.StatusOK, gin.H{"message": "FALSE"})
        }
}
