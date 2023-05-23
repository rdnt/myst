package credentials

type Credentials struct {
	Address    string
	Username   string
	Password   string
	PublicKey  []byte
	PrivateKey []byte
	UserKeys   map[string][]byte
}
