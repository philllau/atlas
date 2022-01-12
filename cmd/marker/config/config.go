package config

import (
	"crypto/ecdsa"
	"github.com/mapprotocol/atlas/cmd/marker/mapprotocol"
	"gopkg.in/urfave/cli.v1"
	"math/big"

	"github.com/mapprotocol/atlas/accounts/abi"
	"github.com/mapprotocol/atlas/cmd/marker/account"
	blscrypto "github.com/mapprotocol/atlas/helper/bls"
	"github.com/mapprotocol/atlas/params"

	"github.com/ethereum/go-ethereum/common"
)

type LockedGoldParameters struct {
	LockedGoldABI     *abi.ABI
	LockedGoldAddress common.Address
}
type AccountsParameters struct {
	AccountsABI     *abi.ABI
	AccountsAddress common.Address
}
type ValidatorParameters struct {
	ValidatorABI     *abi.ABI
	ValidatorAddress common.Address
}
type ElectionParameters struct {
	ElectionABI     *abi.ABI
	ElectionAddress common.Address
}
type GoldTokenParameters struct {
	GoldTokenABI     *abi.ABI
	GoldTokenAddress common.Address
}

type Config struct {
	From                 common.Address
	PublicKey            []byte
	PrivateKey           *ecdsa.PrivateKey
	BlsPub               blscrypto.SerializedPublicKey
	BLSProof             []byte
	Value                uint64
	Duration             int64
	Commission           int64
	Lesser               common.Address
	Greater              common.Address
	VoteNum              *big.Int
	TopNum               *big.Int
	LockedNum            *big.Int
	WithdrawIndex        *big.Int
	RelockIndex          *big.Int
	ValidatorIndex       *big.Int
	TargetAddress        common.Address
	ip                   string
	port                 int
	Verbosity            string
	LockedGoldParameters LockedGoldParameters
	AccountsParameters   AccountsParameters
	ValidatorParameters  ValidatorParameters
	ElectionParameters   ElectionParameters
	GoldTokenParameters  GoldTokenParameters
}

func AssemblyConfig(ctx *cli.Context) *Config {
	config := Config{}
	//------------------ pre set --------------------------
	path := ""
	password := ""
	config.VoteNum = big.NewInt(int64(100))
	config.Lesser = params.ZeroAddress
	config.Greater = params.ZeroAddress
	config.TargetAddress = params.ZeroAddress
	config.Commission = 80
	config.Verbosity = "3"
	//-----------------------------------------------------

	if ctx.IsSet(KeyStoreFlag.Name) {
		path = ctx.String(KeyStoreFlag.Name)
	}
	if ctx.IsSet(PasswordFlag.Name) {
		password = ctx.String(PasswordFlag.Name)
	}

	if ctx.IsSet(CommissionFlag.Name) {
		config.Commission = ctx.Int64(CommissionFlag.Name)
	}
	if ctx.IsSet(LesserFlag.Name) {
		config.Lesser = common.HexToAddress(ctx.String(LesserFlag.Name))
	}
	if ctx.IsSet(GreaterFlag.Name) {
		config.Greater = common.HexToAddress(ctx.String(GreaterFlag.Name))
	}
	if ctx.IsSet(VoteNumFlag.Name) {
		config.VoteNum = big.NewInt(ctx.Int64(VoteNumFlag.Name))
	}
	if ctx.IsSet(TargetAddressFlag.Name) {
		config.TargetAddress = common.HexToAddress(ctx.String(TargetAddressFlag.Name))
	}
	if ctx.IsSet(ValueFlag.Name) {
		config.Value = ctx.Uint64(ValueFlag.Name)
	}
	if ctx.IsSet(DurationFlag.Name) {
		config.Duration = ctx.Int64(DurationFlag.Name)
	}
	if ctx.IsSet(TopNumFlag.Name) {
		config.TopNum = big.NewInt(ctx.Int64(TopNumFlag.Name))
	}
	if ctx.IsSet(LockedNumFlag.Name) {
		config.LockedNum = big.NewInt(ctx.Int64(LockedNumFlag.Name))
	}
	if ctx.IsSet(WithdrawIndexFlag.Name) {
		config.WithdrawIndex = big.NewInt(ctx.Int64(WithdrawIndexFlag.Name))
	}
	if ctx.IsSet(RelockIndexFlag.Name) {
		config.RelockIndex = big.NewInt(ctx.Int64(RelockIndexFlag.Name))
	}
	if ctx.IsSet(ValidatorIndexFlag.Name) {
		config.ValidatorIndex = big.NewInt(ctx.Int64(ValidatorIndexFlag.Name))
	}
	if ctx.IsSet(VerbosityFlag.Name) {
		config.Verbosity = ctx.String(VerbosityFlag.Name)
	}
	account := account.LoadAccount(path, password)
	blsPub, err := account.BLSPublicKey()
	if err != nil {
		return nil
	}
	config.PublicKey = account.PublicKey()
	config.From = account.Address
	config.PrivateKey = account.PrivateKey
	config.BlsPub = blsPub
	config.BLSProof = account.MustBLSProofOfPossession()

	ValidatorAddress := mapprotocol.MustProxyAddressFor("Validators")
	LockedGoldAddress := mapprotocol.MustProxyAddressFor("LockedGold")
	AccountsAddress := mapprotocol.MustProxyAddressFor("Accounts")
	ElectionAddress := mapprotocol.MustProxyAddressFor("Election")
	GoldTokenAddress := mapprotocol.MustProxyAddressFor("GoldToken")
	config.ValidatorParameters.ValidatorAddress = ValidatorAddress
	config.LockedGoldParameters.LockedGoldAddress = LockedGoldAddress
	config.AccountsParameters.AccountsAddress = AccountsAddress
	config.ElectionParameters.ElectionAddress = ElectionAddress
	config.GoldTokenParameters.GoldTokenAddress = GoldTokenAddress

	abiValidators := mapprotocol.AbiFor("Validators")
	abiLockedGold := mapprotocol.AbiFor("LockedGold")
	abiAccounts := mapprotocol.AbiFor("Accounts")
	abiElection := mapprotocol.AbiFor("Election")
	abiGoldToken := mapprotocol.AbiFor("GoldToken")
	config.ValidatorParameters.ValidatorABI = abiValidators
	config.LockedGoldParameters.LockedGoldABI = abiLockedGold
	config.AccountsParameters.AccountsABI = abiAccounts
	config.ElectionParameters.ElectionABI = abiElection
	config.GoldTokenParameters.GoldTokenABI = abiGoldToken

	return &config
}