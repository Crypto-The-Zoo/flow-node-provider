import NFTStorefront from "../../contracts/NFTStorefront.cdc"

transaction {
    prepare(account: AuthAccount) {
        if account.borrow<&NFTStorefront.Storefront>(from: NFTStorefront.StorefrontStoragePath) == nil {

            let storefront <- NFTStorefront.createStorefront()
            account.save(<-storefront, to: NFTStorefront.StorefrontStoragePath)
            account.link<&NFTStorefront.Storefront{NFTStorefront.StorefrontPublic}>(NFTStorefront.StorefrontPublicPath, target: NFTStorefront.StorefrontStoragePath)
        }
    }
}