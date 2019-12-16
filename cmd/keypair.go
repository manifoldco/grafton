package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"golang.org/x/crypto/ed25519"

	"github.com/urfave/cli/v2"

	"github.com/manifoldco/go-base64"
	"github.com/manifoldco/go-signature"
)

// keypair represents a Master Public Keypair used for generating and endorsing
// Live Keypairs which sign HTTP Requests
type keypair struct {
	PublicKey  ed25519.PublicKey  `json:"public_key"`
	PrivateKey ed25519.PrivateKey `json:"private_key"`
}

// liveKeypair represents an endorsed keypair used for signing requests
type liveKeypair struct {
	keypair
	Endorsement *base64.Value
}

func getKeypair() (*keypair, error) {
	keyFile, err := getKeyFilePath()
	if err != nil {
		return nil, cli.NewExitError("Could not determine working directory: "+err.Error(), -1)
	}

	if _, err = os.Stat(keyFile); os.IsNotExist(err) {
		return nil, cli.NewExitError(
			"Master key file does not exist; generate one using 'grafton generate'", -1)
	}

	k, err := loadKeypair(keyFile)
	if err != nil {
		return nil, cli.NewExitError("Could not load master key file: "+err.Error(), -1)
	}

	return k, err
}

// GetKeyFilePath returns the filepath to where the master key should belong
func getKeyFilePath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(cwd, "masterkey.json"), nil
}

// NewKeypair returns a NewKeypair struct
func newKeypair() (*keypair, error) {
	pubKey, privKey, err := ed25519.GenerateKey(nil) // crypto.Rand is used
	if err != nil {
		return nil, err
	}

	return &keypair{PublicKey: pubKey, PrivateKey: privKey}, err
}

// LoadKeypair loads a Keypair from a JSON file into a Keypair struct
func loadKeypair(file string) (*keypair, error) {
	k := &keypair{}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, k)
	if err != nil {
		return nil, err
	}

	return k, err
}

// Save writes the Keypair to a file in JSON
func (k *keypair) save(file string) error {
	b, err := json.Marshal(k)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, b, 0644)
}

// LiveKeypair creates and endorses a Keypair for signing requests
func (k *keypair) liveKeypair() (*liveKeypair, error) {
	pubKey, privKey, err := ed25519.GenerateKey(nil) // crypto.Rand is used
	if err != nil {
		return nil, err
	}

	sig := ed25519.Sign(k.PrivateKey, []byte(pubKey))
	c := &liveKeypair{
		keypair: keypair{
			PublicKey:  pubKey,
			PrivateKey: privKey,
		},
		Endorsement: base64.New(sig),
	}

	return c, nil
}

func emptyKeypair() (*liveKeypair, error) {
	pubKey, privKey, err := ed25519.GenerateKey(nil) // crypto.Rand is used
	if err != nil {
		return nil, err
	}

	return &liveKeypair{
		keypair: keypair{
			PublicKey:  pubKey,
			PrivateKey: privKey,
		},
		Endorsement: base64.New([]byte("not-valid")),
	}, nil
}

// Sign generates a signature using the Live Keypair
func (l *liveKeypair) Sign(b []byte) (*signature.Signature, error) {
	sig := ed25519.Sign(l.PrivateKey, b)

	return &signature.Signature{
		Value:       base64.New(sig),
		PublicKey:   base64.New([]byte(l.PublicKey)),
		Endorsement: l.Endorsement,
	}, nil
}
