
import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"

pub fun main(typeID: UInt64): UInt64 {
  return CryptoZooNFT.totalSupply
}