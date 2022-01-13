import FungibleToken from "../../contracts/FungibleToken.cdc"
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import FlowToken from "../../contracts/FlowToken.cdc"
import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"
import NFTStorefront from "../../contracts/NFTStorefront.cdc"

pub fun getOrCreateCollection(account: AuthAccount): &CryptoZooNFT.Collection{NonFungibleToken.Receiver} {
    if let collectionRef = account.borrow<&CryptoZooNFT.Collection>(from: CryptoZooNFT.CollectionStoragePath) {
        return collectionRef
    }

    let collection <- CryptoZooNFT.createEmptyCollection() as! @CryptoZooNFT.Collection
    let collectionRef = &collection as &CryptoZooNFT.Collection
    account.save(<-collection, to: CryptoZooNFT.CollectionStoragePath)
    account.link<&CryptoZooNFT.Collection{NonFungibleToken.CollectionPublic, CryptoZooNFT.CryptoZooNFTCollectionPublic}>(CryptoZooNFT.CollectionPublicPath, target: CryptoZooNFT.CollectionStoragePath)

    return collectionRef
}

transaction(listingResourceID: UInt64, storefrontAddress: Address) {

    let paymentVault: @FungibleToken.Vault
    let CryptoZooNFTCollection: &CryptoZooNFT.Collection{NonFungibleToken.Receiver}
    let storefront: &NFTStorefront.Storefront{NFTStorefront.StorefrontPublic}
    let listing: &NFTStorefront.Listing{NFTStorefront.ListingPublic}

    prepare(account: AuthAccount) {
        self.storefront = getAccount(storefrontAddress)
            .getCapability<&NFTStorefront.Storefront{NFTStorefront.StorefrontPublic}>(
                NFTStorefront.StorefrontPublicPath
            )!
            .borrow()
            ?? panic("Could not borrow Storefront from provided address")

        self.listing = self.storefront.borrowListing(listingResourceID: listingResourceID)
            ?? panic("No Listing with that ID in Storefront")
        
        let price = self.listing.getDetails().salePrice

        let mainFLOWVault = account.borrow<&FlowToken.Vault>(from: /storage/flowTokenVault)
            ?? panic("Cannot borrow FLOW vault from account storage")
        
        self.paymentVault <- mainFLOWVault.withdraw(amount: price)

        self.CryptoZooNFTCollection = getOrCreateCollection(account: account)
    }

    execute {
        let item <- self.listing.purchase(
            payment: <-self.paymentVault
        )

        self.CryptoZooNFTCollection.deposit(token: <-item)

        self.storefront.cleanup(listingResourceID: listingResourceID)
    }
}