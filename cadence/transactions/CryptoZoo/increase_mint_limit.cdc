import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"

transaction(typeID: UInt64, mintLimit: UInt64) {

    let admin: &CryptoZooNFT.Admin

    prepare(signer: AuthAccount) {

      self.admin = signer.borrow<&CryptoZooNFT.Admin>(from: CryptoZooNFT.AdminStoragePath)
        ?? panic("Could not borrow a reference to the Admin")
    }

    execute {
      self.admin.updateNFTTemplateMintLimit(typeID: typeID, newMintLimit: mintLimit)
    }
}