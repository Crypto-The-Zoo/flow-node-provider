import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"

transaction(recipient: Address, typeID: UInt64) {

    let admin: &CryptoZooNFT.Admin

    prepare(signer: AuthAccount) {

      self.admin = signer.borrow<&CryptoZooNFT.Admin>(from: CryptoZooNFT.AdminStoragePath)!
    }

    execute {
        let recipient = getAccount(recipient)

        let receiver = recipient
            .getCapability(CryptoZooNFT.CollectionPublicPath)!
            .borrow<&{NonFungibleToken.CollectionPublic}>()!

        self.admin.mintNFT(recipient: receiver, typeID: typeID)
    }
}