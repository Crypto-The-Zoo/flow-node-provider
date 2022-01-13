
import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"

pub fun main(typeID: UInt64, coord: [UInt64]): Bool {    
  return CryptoZooNFT.checkIsCoordMinted(typeID: typeID, coord: coord)
}