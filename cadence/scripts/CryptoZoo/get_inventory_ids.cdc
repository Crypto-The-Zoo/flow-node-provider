// flow scripts execute ./cadence/scripts/CryptoZoo/get_inventory_ids.cdc 0x1122aee5915f7fee --network testnet

import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"

pub fun main(ownerAddress: Address): {String: [UInt64]} {
  let owner = getAccount(ownerAddress)
  let ids: {String: [UInt64]} = {}

  if let col = owner.getCapability(CryptoZooNFT.CollectionPublicPath)
  .borrow<&{CryptoZooNFT.CryptoZooNFTCollectionPublic}>() {
      ids["CryptoZooNFT"] = col.getIDs()
  }

  return ids
}