package wallet

import "github.com/taorzhang/toolkit/client"

type Option func(wallet *Account)

func WithEthProvider(client client.Provider) Option {
	return func(wallet *Account) {
		wallet.Client = client
	}
}
