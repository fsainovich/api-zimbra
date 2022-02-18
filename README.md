# api-zimbra
Validate if an email address exists in Zimbra´s LDAP

Made with GOLANG.

### Generate the binary:

```bash
go mod init api-zimbra
go build -x api-zimbra.go
cp api-zimbra /usr/local/bin
```

Access log file: /var/log/api-zimbra-access.log

Error log file: /var/log/api-zimbra-error.log
