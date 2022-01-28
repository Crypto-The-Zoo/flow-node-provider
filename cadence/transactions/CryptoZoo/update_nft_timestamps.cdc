// flow transactions send ./cadence/transactions/CryptoZoo/update_nft_timestamps.cdc 2 '{"availableAt": 0.00000000, "expiresAt": 1643691600.00000000}' --network testnet --signer testnet-minter

import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import CryptoZooNFT from "../../contracts/CryptoZooNFT.cdc"

transaction(typeID: UInt64, timestamps: {String: UFix64}) {

    let admin: &CryptoZooNFT.Admin

    prepare(signer: AuthAccount) {

      self.admin = signer.borrow<&CryptoZooNFT.Admin>(from: CryptoZooNFT.AdminStoragePath)
        ?? panic("Could not borrow a reference to the Admin")
    }

    execute {
      CryptoZooNFT.updateTimestamps(newTimestamps: timestamps)
    }
}