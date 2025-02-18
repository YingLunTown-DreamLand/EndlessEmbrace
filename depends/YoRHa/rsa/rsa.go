package yorha_rsa

// Note: Most of the codes are origin copied from
// https://www.cnblogs.com/akidongzi/p/12036165.html
import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
)

// ...
func SplitSlice(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf)
	}
	return chunks
}

// 公钥加密
func EncryptPKCS1v15(key *rsa.PublicKey, data []byte) ([]byte, error) {
	partLen := key.N.BitLen()/8 - 11
	chunks := SplitSlice(data, partLen)
	buffer := bytes.NewBuffer([]byte{})
	for _, chunk := range chunks {
		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, key, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(bytes)
	}
	return buffer.Bytes(), nil
}

// 私钥解密
func DecryptPKCS1v15(key *rsa.PrivateKey, encrypted []byte) ([]byte, error) {
	partLen := key.N.BitLen() / 8
	chunks := SplitSlice(encrypted, partLen)
	buffer := bytes.NewBuffer([]byte{})
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, key, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(decrypted)
	}
	return buffer.Bytes(), nil
}
