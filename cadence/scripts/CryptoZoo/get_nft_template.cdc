import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"

pub fun main(typeID: UInt64): CryptoZooNFT.CryptoZooNFTTemplate {

  return CryptoZooNFT.getNFTTemplateByTypeID(typeID: typeID)

}