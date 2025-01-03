package cryptoPlus

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

// GenerateRsaKey 生成rsa密钥对，并且保存到磁盘中
func GenerateRsaKey(keySize int) error {
	//使用rsa中的GenerateKey生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return err
	}

	//2. 通过x509标准将得到的rsa私钥序列化为ASN.1的DER编码格式
	derText := x509.MarshalPKCS1PrivateKey(privateKey)

	//3. 要组织一个pem.block
	block := pem.Block{
		Type:  "Prism X Private Key", //这个地方随便写什么
		Bytes: derText,               //ASN.1的DER编码格式
	}

	//4. pem编码
	file, err1 := os.Create("private.pem") //保存在磁盘的文件名
	if err1 != nil {
		return err
	}

	if err = pem.Encode(file, &block); err != nil {
		return err
	}

	if err = file.Close(); err != nil {
		return err
	}

	//生成rsa公钥 ===================================
	//1. 从私钥中取出公钥
	publicKey := privateKey.PublicKey
	//2.使用x509标准序列化  注意参数传参为地址
	derStream, err2 := x509.MarshalPKIXPublicKey(&publicKey)
	if err2 != nil {
		return err
	}

	//3. 将得到的数据放到pem.block中
	block = pem.Block{
		Type:    "Prism X Public Key",
		Headers: nil,
		Bytes:   derStream,
	}

	//pem编码
	file, err = os.Create("public.pem")
	if err != nil {
		return err
	}

	if err = pem.Encode(file, &block); err != nil {
		return err
	}

	if err = file.Close(); err != nil {
		return err
	}
	return nil
}

// RsaPublicEncode 使用rsa公钥加密文件
func RsaPublicEncode(plainText []byte, filename string) []byte {
	//1. 读取公钥信息 放到data变量中
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	stat, _ := file.Stat() //得到文件属性信息
	data := make([]byte, stat.Size())

	if _, err = file.Read(data); err != nil {
		return nil
	}

	if err = file.Close(); err != nil {
		return nil
	}
	//2. 将得到的字符串pem解码
	block, _ := pem.Decode(data)

	//3. 使用x509将编码之后的公钥解析出来
	pubInterface, err2 := x509.ParsePKIXPublicKey(block.Bytes)
	if err2 != nil {
		return nil
	}
	pubKey := pubInterface.(*rsa.PublicKey)

	//4. 使用公钥加密
	cipherText, err3 := rsa.EncryptPKCS1v15(rand.Reader, pubKey, plainText)
	if err3 != nil {
		return nil
	}
	return cipherText
}

// RsaPrivateDecode 使用rsa私钥解密
func RsaPrivateDecode(cipherText []byte, private []byte) []byte {
	//1. 将得到的字符串进行pem解码
	block, _ := pem.Decode(private)
	//2. 使用x509将编码之后的私钥解析出来
	privateKey, err3 := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err3 != nil {
		panic(err3)
	}
	//4. 使用私钥将数据解密
	plainText, err4 := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err4 != nil {
		panic(err4)
	}
	return plainText
}
