package bchatcommon

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"

	"bytes"
	"io"
	"math/rand"
	"os"

	"bufio"

	"io/ioutil"
	"log"
	//"runtime"

	"crypto/sha256"

	"fmt"
	"math"
	"strconv"
	"strings"
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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func LineCounter(r io.Reader, numBytes int) (int, error) {
	buf := make([]byte, numBytes*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}

func ReadStream(f string) {
	file, err := os.Open(f)
	//fmt.Println(LineCounter(file, 2))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains("deanandkayla", scanner.Text()) {
			fmt.Println("FOUND")
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func ReadMem(f string) {
	var wg sync.WaitGroup
	b, e := ioutil.ReadFile(f)
	if e != nil {
		panic(e)
	}
	array := bytes.Split(b, []byte("\n"))
	//fmt.Println(b)

	N := 32
	perGo := len(array) / N
	//leftOver := len(array) % runtime.NumCPU()
	//rem := len(array) % perGo

	for j := 0; j < N; j++ {

		var goChunk [][]byte
		//fmt.Println(j, j * (len(array) / runtime.NumCPU()), perGo)
		if j+1 == N {
			goChunk = array[j * (len(array) / N):]
		} else {
			goChunk = array[j * (len(array) / N): perGo]
		}
		wg.Add(1)
		go func() {
			for j := 0; j < len(goChunk); j++ {
				//if bytes.Contains(goChunk[j], []byte("asdf")) {
				encrpt := sha256.Sum256(goChunk[j])
				_=encrpt
			}
			wg.Done()
			//fmt.Println()
		}()
		perGo = perGo + len(array) / N
	}
	fmt.Println("waiting on grps")
	wg.Wait()
}

func ChunkFIle() {
	         fileToBeChunked := "../biggerread.txt"

         file, err := os.Open(fileToBeChunked)

         if err != nil {
                 fmt.Println(err)
                 os.Exit(1)
         }

         defer file.Close()

         fileInfo, _ := file.Stat()

         var fileSize int64 = fileInfo.Size()

         const fileChunk = 1000 * (1 << 20) // 1 MB, change this to your requirement

         // calculate total number of parts the file will be chunked into

         totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

         fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)

         for i := uint64(0); i < totalPartsNum; i++ {

                 partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
                 partBuffer := make([]byte, partSize)

                 file.Read(partBuffer)

                 // write to disk
                 fileName := "somebigfile_" + strconv.FormatUint(i, 10)
                 _, err := os.Create(fileName)

                 if err != nil {
                         fmt.Println(err)
                         os.Exit(1)
                 }

                 // write/save buffer to disk
                 ioutil.WriteFile(fileName, partBuffer, os.ModeAppend)

                 fmt.Println("Split to : ", fileName)
         }
}

//func main() {
//	readStream("smallread2.txt")
//}

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
