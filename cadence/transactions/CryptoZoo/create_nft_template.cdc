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