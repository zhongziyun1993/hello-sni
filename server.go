package main

import (
	"crypto/tls"
	"fmt"
	"crypto/x509"
	"net"
)

var ExampleCert = `-----BEGIN CERTIFICATE-----
MIIDIjCCAgqgAwIBAgIBATANBgkqhkiG9w0BAQUFADCBoTELMAkGA1UEBhMCVVMx
CzAJBgNVBAgMAlNDMRUwEwYDVQQHDAxEZWZhdWx0IENpdHkxHDAaBgNVBAoME0Rl
ZmF1bHQgQ29tcGFueSBMdGQxEDAOBgNVBAsMB1Rlc3QgQ0ExGjAYBgNVBAMMEXd3
dy5leGFtcGxlY2EuY29tMSIwIAYJKoZIhvcNAQkBFhNleGFtcGxlQGV4YW1wbGUu
Y29tMB4XDTE1MDExMjE0MTk0MVoXDTE2MDExMjE0MTk0MVowfDEYMBYGA1UEAwwP
d3d3LmV4YW1wbGUuY29tMQswCQYDVQQIDAJTQzELMAkGA1UEBhMCVVMxIjAgBgkq
hkiG9w0BCQEWE2V4YW1wbGVAZXhhbXBsZS5jb20xEDAOBgNVBAoMB0V4YW1wbGUx
EDAOBgNVBAsMB0V4YW1wbGUwgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAMrv
gu6ZTTefNN7jjiZbS/xvQjyXjYMN7oVXv76jbX8gjMOmg9m0xoVZZFAE4XyQDuCm
47VRx5Qrf/YLXmB2VtCFvB0AhXr5zSeWzPwaAPrjA4ebG+LUo24ziS8KqNxrFs1M
mNrQUgZyQC6XIe1JHXc9t+JlL5UZyZQC1IfaJulDAgMBAAGjDTALMAkGA1UdEwQC
MAAwDQYJKoZIhvcNAQEFBQADggEBAFCi7ZlkMnESvzlZCvv82Pq6S46AAOTPXdFd
TMvrh12E1sdVALF1P1oYFJzG1EiZ5ezOx88fEDTW+Lxb9anw5/KJzwtWcfsupf1m
V7J0D3qKzw5C1wjzYHh9/Pz7B1D0KthQRATQCfNf8s6bbFLaw/dmiIUhHLtIH5Qc
yfrejTZbOSP77z8NOWir+BWWgIDDB2//3AkDIQvT20vmkZRhkqSdT7et4NmXOX/j
jhPti4b2Fie0LeuvgaOdKjCpQQNrYthZHXeVlOLRhMTSk3qUczenkKTOhvP7IS9q
+Dzv5hqgSfvMG392KWh5f8xXfJNs4W5KLbZyl901MeReiLrPH3w=
-----END CERTIFICATE-----`

var ExampleKey = `-----BEGIN PRIVATE KEY-----
MIICeAIBADANBgkqhkiG9w0BAQEFAASCAmIwggJeAgEAAoGBAMrvgu6ZTTefNN7j
jiZbS/xvQjyXjYMN7oVXv76jbX8gjMOmg9m0xoVZZFAE4XyQDuCm47VRx5Qrf/YL
XmB2VtCFvB0AhXr5zSeWzPwaAPrjA4ebG+LUo24ziS8KqNxrFs1MmNrQUgZyQC6X
Ie1JHXc9t+JlL5UZyZQC1IfaJulDAgMBAAECgYEAnxOjEj/vrLNLMZE1Q9H7PZVF
WdP/JQVNvQ7tCpZ3ZdjxHwkvf//aQnuxS5yX2Rnf37BS/TZu+TIkK4373CfHomSx
UTAn2FsLmOJljupgGcoeLx5K5nu7B7rY5L1NHvdpxZ4YjeISrRtEPvRakllENU5y
gJE8c2eQOx08ZSRE4TkCQQD7dws2/FldqwdjJucYijsJVuUdoTqxP8gWL6bB251q
elP2/a6W2elqOcWId28560jG9ZS3cuKvnmu/4LG88vZFAkEAzphrH3673oTsHN+d
uBd5uyrlnGjWjuiMKv2TPITZcWBjB8nJDSvLneHF59MYwejNNEof2tRjgFSdImFH
mi995wJBAMtPjW6wiqRz0i41VuT9ZgwACJBzOdvzQJfHgSD9qgFb1CU/J/hpSRIM
kYvrXK9MbvQFvG6x4VuyT1W8mpe1LK0CQAo8VPpffhFdRpF7psXLK/XQ/0VLkG3O
KburipLyBg/u9ZkaL0Ley5zL5dFBjTV2Qkx367Ic2b0u9AYTCcgi2DsCQQD3zZ7B
v7BOm7MkylKokY2MduFFXU0Bxg6pfZ7q3rvg8gqhUFbaMStPRYg6myiDiW/JfLhF
TcFT4touIo7oriFJ
-----END PRIVATE KEY-----`

var ExampleCACert = `-----BEGIN CERTIFICATE-----
MIIEFzCCAv+gAwIBAgIJALK1iUpF2VQLMA0GCSqGSIb3DQEBBQUAMIGhMQswCQYD
VQQGEwJVUzELMAkGA1UECAwCU0MxFTATBgNVBAcMDERlZmF1bHQgQ2l0eTEcMBoG
A1UECgwTRGVmYXVsdCBDb21wYW55IEx0ZDEQMA4GA1UECwwHVGVzdCBDQTEaMBgG
A1UEAwwRd3d3LmV4YW1wbGVjYS5jb20xIjAgBgkqhkiG9w0BCQEWE2V4YW1wbGVA
ZXhhbXBsZS5jb20wHhcNMTUwMTEyMTQxNTAxWhcNMjUwMTA5MTQxNTAxWjCBoTEL
MAkGA1UEBhMCVVMxCzAJBgNVBAgMAlNDMRUwEwYDVQQHDAxEZWZhdWx0IENpdHkx
HDAaBgNVBAoME0RlZmF1bHQgQ29tcGFueSBMdGQxEDAOBgNVBAsMB1Rlc3QgQ0Ex
GjAYBgNVBAMMEXd3dy5leGFtcGxlY2EuY29tMSIwIAYJKoZIhvcNAQkBFhNleGFt
cGxlQGV4YW1wbGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA
w2rK1J2NMtQj0KDug7g7HRKl5jbf0QMkMKyTU1fBtZ0cCzvsF4CqV11LK4BSVWaK
rzkaXe99IVJnH8KdOlDl5Dh/+cJ3xdkClSyeUT4zgb6CCBqg78ePp+nN11JKuJlV
IG1qdJpB1J5O/kCLsGcTf7RS74MtqMFo96446Zvt7YaBhWPz6gDaO/TUzfrNcGLA
EfHVXkvVWqb3gqXUztZyVex/gtP9FXQ7gxTvJml7UkmT0VAFjtZnCqmFxpLZFZ15
+qP9O7Q2MpsGUO/4vDAuYrKBeg1ZdPSi8gwqUP2qWsGd9MIWRv3thI2903BczDc7
r8WaIbm37vYZAS9G56E4+wIDAQABo1AwTjAdBgNVHQ4EFgQUugLrSJshOBk5TSsU
ANs4+SmJUGwwHwYDVR0jBBgwFoAUugLrSJshOBk5TSsUANs4+SmJUGwwDAYDVR0T
BAUwAwEB/zANBgkqhkiG9w0BAQUFAAOCAQEAaMJ33zAMV4korHo5aPfayV3uHoYZ
1ChzP3eSsF+FjoscpoNSKs91ZXZF6LquzoNezbfiihK4PYqgwVD2+O0/Ty7UjN4S
qzFKVR4OS/6lCJ8YncxoFpTntbvjgojf1DEataKFUN196PAANc3yz8cWHF4uvjPv
WkgFqbIjb+7D1YgglNyovXkRDlRZl0LD1OQ0ZWhd4Ge1qx8mmmanoBeYZ9+DgpFC
j9tQAbS867yeOryNe7sEOIpXAAqK/DTu0hB6+ySsDfMo4piXCc2aA/eI2DCuw08e
w17Dz9WnupZjVdwTKzDhFgJZMLDqn37HQnT6EemLFqbcR0VPEnfyhDtZIQ==
-----END CERTIFICATE-----`

func makeCertPool() (tls.Certificate, *x509.CertPool){
	privateKey := []byte(ExampleKey)
	cert := []byte(ExampleCert)
	caCert := []byte(ExampleCACert)
	combinedCert := append(cert, caCert...)

	tlsCert, err := tls.X509KeyPair(combinedCert, privateKey)

	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(combinedCert) {
		panic("Failed appending certs")
	}

	return tlsCert, certPool
}

func handleClient(conn net.Conn) {
	_, err := fmt.Fprintf(conn, "Hello TLS\n")
	if err != nil {
		fmt.Printf("Error on connection:%v", err)
	}
	conn.Close()

}

func main() {
	tlsCert, certPool := makeCertPool()

	cfg := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		RootCAs: certPool,
		ClientCAs: certPool,
		ClientAuth: tls.NoClientCert,
	}

	listener, _ := tls.Listen("tcp", ":9999", cfg)

	fmt.Printf("Listening for TCP connections on 9999\n")

	for {
		conn, err := listener.Accept()

		if err != nil {
			panic(err)
		}
		fmt.Printf("Accepted connection\n")
		go handleClient(conn)
	}
}
