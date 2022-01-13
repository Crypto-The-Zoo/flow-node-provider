import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"
import MetadataViews from "../../contracts/MetadataViews.cdc"

pub struct NFTResult {
    pub(set) var name: String
    pub(set) var description: String
    pub(set) var uri: String
    pub(set) var cid: String
    pub(set) var mimetype: String
    pub(set) var owner: Address
    pub(set) var type: String

    init() {
        self.name = ""
        self.description = ""
        self.mimetype = ""
        self.uri = ""
        self.cid = ""
        self.owner = 0x0
        self.type = ""
    }
}

pub fun main(address: Address, id: UInt64): NFTResult {
    let account = getAccount(address)

    let collection = account
        .getCapability(CryptoZooNFT.CollectionPublicPath)
        .borrow<&{CryptoZooNFT.CryptoZooNFTCollectionPublic}>()
        ?? panic("Could not borrow a reference to the collection")

    let nft = collection.borrowCryptoZooNFT(id: id)!

    var data = NFTResult()

    // Get the basic display information for this NFT
    if let view = nft.resolveView(Type<MetadataViews.Display>()) {
        let display = view as! MetadataViews.Display

        data.name = display.name
        data.description = display.description
    }

    // Get the http display information for this NFT
    if let view = nft.resolveView(Type<MetadataViews.HTTPThumbnail>()) {
        let display = view as! MetadataViews.HTTPThumbnail

        data.uri = display.uri
        data.mimetype = display.mimetype
    }

    // Get the i[fs] display information for this NFT
    if let view = nft.resolveView(Type<MetadataViews.IPFSThumbnail>()) {
        let display = view as! MetadataViews.IPFSThumbnail

        data.cid = display.cid
        data.mimetype = display.mimetype
    }

    // The owner is stored directly on the NFT object
    let owner: Address = nft.owner!.address!

    data.owner = owner

    // Inspect the type of this NFT to verify its origin
    let nftType = nft.getType()

    data.type = nftType.identifier
    // `data.type` is `A.f3fcd2c1a78f5eee.CryptoZooNFT.NFT`

    return data
}