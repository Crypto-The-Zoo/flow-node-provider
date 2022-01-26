
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"

pub fun main(typeID: UInt64): Bool {
  return CryptoZooNFT.isNFTTemplateExist(typeID: typeID)
}