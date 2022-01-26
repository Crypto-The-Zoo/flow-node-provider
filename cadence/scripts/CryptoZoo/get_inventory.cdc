// flow scripts execute ./cadence/scripts/CryptoZoo/get_inventory.cdc 0x1122aee5915f7fee --network testnet

import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"
  
  pub struct NFTCollection {
    pub let owner: Address
    pub let nfts: [NFTData]

    init(owner: Address) {
        self.owner = owner
        self.nfts = []
    }
  }

  pub struct NFTData {
    pub let contract: NFTContract
    pub let id: UInt64
    pub let uuid: UInt64?
    pub let name: String?
    pub let description: String?
    pub let external_domain_view_url: String?
    pub let token_uri: String?
    pub let media: [NFTMedia?]
    pub let metadata: {String: AnyStruct}
    pub let serialNumber: UInt64
    pub let mintLimit: UInt64
    pub let quality: String?

    init(
        contract: NFTContract,
        id: UInt64,
        uuid: UInt64?,
        name: String?,
        description: String?,
        external_domain_view_url: String?,
        token_uri: String?,
        media: [NFTMedia?],
        metadata: {String: AnyStruct},
        serialNumber: UInt64,
        mintLimit: UInt64,
        quality: String?,
    ) {
        self.contract = contract
        self.id = id
        self.uuid = uuid
        self.name = name
        self.description = description
        self.external_domain_view_url = external_domain_view_url
        self.token_uri = token_uri
        self.media = media
        self.metadata = metadata
        self.serialNumber = serialNumber
        self.mintLimit = mintLimit
        self.quality = quality
    }
  }

  pub struct NFTContract {
    pub let name: String
    pub let address: Address
    pub let storage_path: String
    pub let public_path: String
    pub let public_collection_name: String
    pub let external_domain: String

    init(
        name: String,
        address: Address,
        storage_path: String,
        public_path: String,
        public_collection_name: String,
        external_domain: String
    ) {
        self.name = name
        self.address = address
        self.storage_path = storage_path
        self.public_path = public_path
        self.public_collection_name = public_collection_name
        self.external_domain = external_domain
    }
  }

  pub struct NFTMedia {
    pub let uri: String?
    pub let mimetype: String?

    init(
        uri: String?,
        mimetype: String?
    ) {
        self.uri = uri
        self.mimetype = mimetype
    }
  }

  pub fun main(ownerAddress: Address): [NFTData?] {

    let owner = getAccount(ownerAddress)
    let ids: {String: [UInt64]} = {}
    let NFTs: [NFTData?] = []

    if let col = owner.getCapability(CryptoZooNFT.CollectionPublicPath)
    .borrow<&{CryptoZooNFT.CryptoZooNFTCollectionPublic}>() {
        ids["InceptionAnimals"] = col.getIDs()
    }

    for key in ids.keys {
        for id in ids[key]! {
            var d: NFTData? = nil

            switch key {
                case "InceptionAnimals": d = getInceptionAnimals(owner: owner, id: id)
                default:
                    panic("adapter for NFT not found: ".concat(key))
            }

            NFTs.append(d)
        }
    }

    return NFTs
  }

  // https://flow-view-source.com/testnet/account/0xd60702f03bcafd46
  pub fun getInceptionAnimals(owner: PublicAccount, id: UInt64): NFTData? {
      let contract = NFTContract(
          name: "InceptionAnimals",
          address: 0xd60702f03bcafd46,
          storage_path: "CryptoZooNFT.CollectionStoragePath",
          public_path: "CryptoZooNFT.CollectionPublicPath",
          public_collection_name: "CryptoZooNFT.CryptoZooNFTCollectionPublic",
          external_domain: "https://www.inceptionanimals.com/"
      )
  
      let col = owner.getCapability(CryptoZooNFT.CollectionPublicPath)
          .borrow<&{CryptoZooNFT.CryptoZooNFTCollectionPublic}>()
      if col == nil { return nil }
  
      let nft = col!.borrowCryptoZooNFT(id: id)
      if nft == nil { return nil }
  
      return NFTData(
          contract: contract,
          id: nft!.id,
          uuid: nft!.uuid,
          name: nft!.name,
          description: nft!.getNFTTemplate()!.description,
          external_domain_view_url: nil,
          token_uri: nft!.getNFTTemplate()!.getMetadata()["uri"]!,
          media: [NFTMedia(uri: nft!.getNFTTemplate()!.getMetadata()["uri"]!, mimetype: nft!.getNFTTemplate()!.getMetadata()["mimetype"]!)],
          metadata: nft!.getNFTTemplate()!.getMetadata()!,
          serialNumber: nft!.serialNumber,
          mintLimit: nft!.getNFTTemplate()!.mintLimit,
          quality: nft!.getNFTTemplate()!.getMetadata()["quality"],
      )
  }