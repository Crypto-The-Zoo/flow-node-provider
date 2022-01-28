
import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"

pub fun main(typeID: UInt64): Bool {    
  return CryptoZooNFT.isNFTTemplateExpired(typeID: typeID)
}