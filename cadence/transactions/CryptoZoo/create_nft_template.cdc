// flow transactions send ./cadence/transactions/CryptoZoo/create_nft_template.cdc 1 false 'Spirit Animal x Hippo' 'You can use this Hippo as your character in our metaverse' 10 10.00000000 1000000.00000000 '{"uri": "https://www.inceptionanimals.com/hippo.png", "mimetype": "image/png"}' '{"availableAt": 1000.00000000, "expiresAt": 9999999999.00000000}' false --signer testnet-mintee --network testnet

import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"

transaction(typeID: UInt64, isPack: Bool, name: String, description: String, mintLimit: UInt64, priceUSD: UFix64, priceFlow: UFix64, metadata: {String: String}, timestamps: {String: UFix64}, isLand: Bool) {

  let admin: &CryptoZooNFT.Admin

  prepare(signer: AuthAccount) {
    self.admin = signer.borrow<&CryptoZooNFT.Admin>(from: CryptoZooNFT.AdminStoragePath)
      ?? panic("Could not borrow a reference to the CryptoZooNFT Admin")
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