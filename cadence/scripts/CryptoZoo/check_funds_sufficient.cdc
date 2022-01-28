import FungibleToken from "../../contracts/FungibleToken.cdc"
import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"
import FlowToken from "../../contracts/FlowToken.cdc"

pub fun main(address: Address, typeID: UInt64): Bool {
  let priceFlow = CryptoZooNFT.getNFTTemplateByTypeID(typeID: typeID).priceFlow
  let flowVaultRef = getAccount(address)
    .getCapability(/public/flowTokenBalance)
    .borrow<&FlowToken.Vault{FungibleToken.Balance}>()
    ?? panic("Could not borrow Balance reference to the Vault")
  return flowVaultRef.balance >= priceFlow
}