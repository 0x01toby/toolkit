package wallet

import "github.com/taorzhang/toolkit/client"

type Option func(wallet *Wallet)

func WithEthProvider(client client.Provider) Option {
	return func(wallet *Wallet) {
		wallet.Client = client
	}
}
