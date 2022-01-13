import FungibleToken from "../../contracts/FungibleToken.cdc"
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import CryptoZooNFTMinter from "../../contracts/CryptoZooNFTMinter.cdc"
import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"
import FlowToken from "../../contracts/FlowToken.cdc"

transaction(typeID: UInt64) {

  let paymentVault: @FungibleToken.Vault
  let receiverRef: &CryptoZooNFT.Collection{NonFungibleToken.CollectionPublic}

  prepare(signer: AuthAccount) {

    if (CryptoZooNFT.getNFTTemplateByTypeID(typeID: typeID).isLand) {
      panic("land cannot be purchased directly")
    }

    if signer.borrow<&CryptoZooNFT.Collection>(from: CryptoZooNFT.CollectionStoragePath) == nil {
      signer.save(<- CryptoZooNFT.createEmptyCollection(), to: CryptoZooNFT.CollectionStoragePath)
      signer.unlink(CryptoZooNFT.CollectionPublicPath)
      signer.link<&CryptoZooNFT.Collection{NonFungibleToken.CollectionPublic, CryptoZooNFT.CryptoZooNFTCollectionPublic}>(CryptoZooNFT.CollectionPublicPath, target: CryptoZooNFT.CollectionStoragePath)
    }

    self.receiverRef = signer.borrow<&CryptoZooNFT.Collection{NonFungibleToken.CollectionPublic}>(from: CryptoZooNFT.CollectionStoragePath)
      ?? panic ("Could not get receiver reference to the NFT Collection")
    let price = CryptoZooNFT.getNFTTemplateByTypeID(typeID: typeID).priceFlow
    let signerFlowTokenVault = signer.borrow<&FlowToken.Vault>(from: /storage/flowTokenVault)
      ?? panic("Cannot borrow FlowToken vault from signer storage")
    self.paymentVault <- signerFlowTokenVault.withdraw(amount: price)
  }

  execute {
    CryptoZooNFTMinter.mintNFTWithFlow(recipient: self.receiverRef, typeID: typeID, paymentVault: <- self.paymentVault)
  }
}