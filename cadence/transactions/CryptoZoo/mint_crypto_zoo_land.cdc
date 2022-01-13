import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"

transaction(recipient: Address, typeID: UInt64, coord: [UInt64], nftName: String) {

    let admin: &CryptoZooNFT.Admin

    prepare(signer: AuthAccount) {

      self.admin = signer.borrow<&CryptoZooNFT.Admin>(from: CryptoZooNFT.AdminStoragePath)
        ?? panic("Could not borrow a reference to the Admin")
    }

    execute {
        let recipient = getAccount(recipient)

        let receiver = recipient
            .getCapability(CryptoZooNFT.CollectionPublicPath)!
            .borrow<&{NonFungibleToken.CollectionPublic}>()
            ?? panic("Could not get receiver reference to the NFT Collection")

        self.admin.mintLandNFT(recipient: receiver, typeID: typeID, coord: coord, nftName: nftName)
    }
}