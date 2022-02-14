import CryptoZooNFT from 0x8ea44ab931cac762

pub fun main(typeID: UInt64): CryptoZooNFT.CryptoZooNFTTemplate {

  return CryptoZooNFT.getNFTTemplateByTypeID(typeID: typeID)

}