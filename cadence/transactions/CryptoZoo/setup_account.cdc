// flow transactions send ./cadence/transactions/CryptoZoo/setup_account.cdc --signer testnet-mintee --network testnet
// flow transactions send ./cadence/transactions/CryptoZoo/setup_account.cdc --signer testnet-mintee --network testnet

import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"

transaction {
  prepare(signer: AuthAccount) {
    let hasCollection = getAccount(signer.address)
      .getCapability<&CryptoZooNFT.Collection>(CryptoZooNFT.CollectionPublicPath)
      .check()
    if hasCollection {
      panic("A Collection is already setup for this account")
    }
    if signer.borrow<&CryptoZooNFT.Collection>(from: CryptoZooNFT.CollectionStoragePath) == nil {
      signer.save(<- CryptoZooNFT.createEmptyCollection(), to: CryptoZooNFT.CollectionStoragePath)
    }
    signer.unlink(CryptoZooNFT.CollectionPublicPath)
    signer.link<&CryptoZooNFT.Collection{NonFungibleToken.CollectionPublic, CryptoZooNFT.CryptoZooNFTCollectionPublic}>(CryptoZooNFT.CollectionPublicPath, target: CryptoZooNFT.CollectionStoragePath)
  }
}