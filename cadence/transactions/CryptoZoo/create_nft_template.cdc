// flow transactions send ./cadence/transactions/CryptoZoo/create_nft_template.cdc 3 false "Inception LAND VOUCHER" "An Inception LAND VOUCHER that can redeem in exchange for an Inception REGULAR PLOT" 100 1000000.00000000 30.00000000 '{"quality": "rare", "uri": "https://storage.googleapis.com/inception_public/regular_plot.png", "mimetype": "image/png"}' '{"availableAt": 0.00000000,"expiresAt": 1706763600.00000000}' false --network testnet --signer testnet-minter

import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"

transaction(typeID: UInt64, isPack: Bool, name: String, description: String, mintLimit: UInt64, priceUSD: UFix64, priceFlow: UFix64, metadata: {String: String}, timestamps: {String: UFix64}, isLand: Bool) {

  let admin: &CryptoZooNFT.Admin

  prepare(signer: AuthAccount) {
    self.admin = signer.borrow<&CryptoZooNFT.Admin>(from: CryptoZooNFT.AdminStoragePath)!
  }

  execute {
    self.admin.createNFTTemplate(
      typeID: typeID,
      isPack: isPack,
      name: name,
      description: description,
      mintLimit: mintLimit,
      priceUSD: priceUSD,
      priceFlow: priceFlow,
      metadata: metadata,
      timestamps: timestamps,
      isLand: isLand,
    )
  }

}