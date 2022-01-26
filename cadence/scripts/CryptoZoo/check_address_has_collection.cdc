import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"

pub fun main(address: Address): Bool {
  let account = getAccount(address)

  var collectionCap = account.getCapability<&{NonFungibleToken.CollectionPublic}>(CryptoZooNFT.CollectionPublicPath)

  if collectionCap.check() {
    return true
  } else {
    return false
  }
}