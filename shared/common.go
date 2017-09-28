package bchatcommon

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

const B_CONNECT = "B_CONNECT"
const B_MESSAGE = "B_MESSAGE"
const B_NAMECHANGE = "B_NAMECHANGE"
const B_ROOMCHANGE = "B_ROOMCHANGE"
const B_DISCONNECT = "B_DISCONNECT"
const MAIN_ROOM = "MAIN_ROOM"

type BChatClient struct {
	Name   string
	Room   string
	WsConn *websocket.Conn
	Uid    string
	mu     sync.Mutex
}

func (c *BChatClient) ChangeName(s string) {
	c.Name = s
}

func (c *BChatClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.WsConn.Close()
}

func (c *BChatClient) ChangeRoom(s string) {
	c.Room = s
}

func (c *BChatClient) SendMessage(s BMessage) {
	c.mu.Lock()
	c.WsConn.WriteJSON(s)
	c.mu.Unlock()
}

func (c *BChatClient) ReadMessage() (BMessage, error) {
	c.mu.Lock()
	var bMessage BMessage
	err := c.WsConn.ReadJSON(&bMessage)
	c.mu.Unlock()
	return bMessage, err
}

type BMessage struct {
	MsgType    string
	Name       string
	Room       string
	Payload    string
	TimeStamp  time.Time
	Uid        string
	OnlineData string
	RoomData   string
}

//func GenCreds() (credentials.TransportCredentials){
//	crt, key := CreatePemKey()
//	certificate, err := tls.X509KeyPair(crt, key)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	certPool := x509.NewCertPool()
//	ca, err := ioutil.ReadFile("GIAG3.crt")
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	if ok := certPool.AppendCertsFromPEM(ca); !ok {
//		log.Fatalln("unable to append certificate")
//	}
//
//	// Create the TLS credentials
//	creds := credentials.NewTLS(&tls.Config{
//		ServerName:         "ASDF",
//		Certificates:       []tls.Certificate{certificate},
//		RootCAs:            certPool,
//		InsecureSkipVerify: true,
//	})
//	return creds
//}
//func GetMD5Hash(text string) string {
//    hasher := md5.New()
//    hasher.Write([]byte(text))
//    return hex.EncodeToString(hasher.Sum(nil))
//}
//
//func GetOutboundIP() string {
//    conn, err := net.Dial("udp", "8.8.8.8:80")
//    if err != nil {
//        log.Fatal(err)
//    }
//    defer conn.Close()
//    localAddr := conn.LocalAddr().(*net.UDPAddr)
//    return localAddr.IP.String()
//}
//
//func publicKey(priv interface{}) interface{} {
//	switch k := priv.(type) {
//	case *rsa.PrivateKey:
//		return &k.PublicKey
//	case *ecdsa.PrivateKey:
//		return &k.PublicKey
//	default:
//		return nil
//	}
//}
//
//func pemBlockForKey(priv interface{}) *pem.Block {
//	switch k := priv.(type) {
//	case *rsa.PrivateKey:
//		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
//	case *ecdsa.PrivateKey:
//		b, err := x509.MarshalECPrivateKey(k)
//		if err != nil {
//			log.Fatalf("Unable to marshal ECDSA private key: %v", err)
//			os.Exit(2)
//		}
//		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
//	default:
//		return nil
//	}
//}
//
//func CreatePemKey() (certpem, keypem []byte) {
//	priv, _ := rsa.GenerateKey(rand.Reader, 1024)
//	notBefore := time.Now()
//	notAfter := notBefore.AddDate(1, 0, 0)
//	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
//	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)
//
//	template := x509.Certificate{
//		SerialNumber: serialNumber,
//		Subject: pkix.Name{
//			Organization: []string{"Acme Co"},
//		},
//		NotBefore:             notBefore,
//		NotAfter:              notAfter,
//		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
//		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
//		BasicConstraintsValid: true,
//	}
//	template.IsCA = true
//	derbytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
//	certpem = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derbytes})
//	keypem = pem.EncodeToMemory(pemBlockForKey(priv))
//	return certpem, keypem
//}
