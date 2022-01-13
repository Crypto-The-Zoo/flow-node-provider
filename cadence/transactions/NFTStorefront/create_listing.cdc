import FungibleToken from "../../contracts/FungibleToken.cdc"
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import FlowToken from "../../contracts/FlowToken.cdc"
import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"
import NFTStorefront from "../../contracts/NFTStorefront.cdc"

pub fun getOrCreateStorefront(account: AuthAccount): &NFTStorefront.Storefront {
    if let storefrontRef = account.borrow<&NFTStorefront.Storefront>(from: NFTStorefront.StorefrontStoragePath) {
        return storefrontRef
    }

    let storefront <- NFTStorefront.createStorefront()

    let storefrontRef = &storefront as &NFTStorefront.Storefront

    account.save(<-storefront, to: NFTStorefront.StorefrontStoragePath)

    account.link<&NFTStorefront.Storefront{NFTStorefront.StorefrontPublic}>(NFTStorefront.StorefrontPublicPath, target: NFTStorefront.StorefrontStoragePath)

    return storefrontRef
}

transaction(saleItemID: UInt64, saleItemPrice: UFix64) {

    let flowReceiver: Capability<&FlowToken.Vault{FungibleToken.Receiver}>
    let CryptoZooNFTProvider: Capability<&CryptoZooNFT.Collection{NonFungibleToken.Provider, NonFungibleToken.CollectionPublic}>
    let storefront: &NFTStorefront.Storefront

    prepare(account: AuthAccount) {
        // We need a provider capability, but one is not provided by default so we create one if needed.
        let CryptoZooNFTCollectionProviderPrivatePath = /private/CryptoZooNFTCollectionProvider

        self.flowReceiver = account.getCapability<&FlowToken.Vault{FungibleToken.Receiver}>(/public/flowTokenReceiver)!

        assert(self.flowReceiver.borrow() != nil, message: "Missing or mis-typed FLOW receiver")

        if !account.getCapability<&CryptoZooNFT.Collection{NonFungibleToken.Provider, NonFungibleToken.CollectionPublic}>(CryptoZooNFTCollectionProviderPrivatePath)!.check() {
            account.link<&CryptoZooNFT.Collection{NonFungibleToken.Provider, NonFungibleToken.CollectionPublic}>(CryptoZooNFTCollectionProviderPrivatePath, target: CryptoZooNFT.CollectionStoragePath)
        }

        self.CryptoZooNFTProvider = account.getCapability<&CryptoZooNFT.Collection{NonFungibleToken.Provider, NonFungibleToken.CollectionPublic}>(CryptoZooNFTCollectionProviderPrivatePath)!

        assert(self.CryptoZooNFTProvider.borrow() != nil, message: "Missing or mis-typed CryptoZooNFT.Collection provider")

        self.storefront = getOrCreateStorefront(account: account)
    }

    execute {
        let saleCut = NFTStorefront.SaleCut(
            receiver: self.flowReceiver,
            amount: saleItemPrice
        )
        self.storefront.createListing(
            nftProviderCapability: self.CryptoZooNFTProvider,
            nftType: Type<@CryptoZooNFT.NFT>(),
            nftID: saleItemID,
            salePaymentVaultType: Type<@FlowToken.Vault>(),
            saleCuts: [saleCut]
        )
    }
}