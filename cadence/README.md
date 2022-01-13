# Crypto Zoo Metaverse Cadence Smart Contract

# Background

- Crypto The Zoo project is a open-world browser 3D game platform for animal lovers to collect collectibles, including but not limited to lands, avatars, characters, and in-game inventory.
- This repository contains the cadence smart contracts for three core functionalities that the game requires
  - CryptoZooNFT.cdc for administrative management of NFTs and lands in the form of NFTs
  - CryptoZooNFTStorefront.cdc for peer-to-peer trading purposes in the form of secondary market
  - CryptoZooNFTMinter.cdc for users to mint CryptoZooNFTs and lands by paying Flow
# Run tests
Please run the script below to trigger three suites of smart contract tests. 
```
cd tests \
yarn test
```

# Tests

- `crypto-zoo.test.js`
  -  Deploy CryptoZooNFT
  -  Setup account for CryptoZooNFT
  -  Create CryptoZooNFT templates: for both NFTs and land type of NFTs
  -  Admin minter to mint a CryptoZooNFT NFT
  -  Admin minter to mint a CryptoZooNFT Land
  -  Setup new empty CryptoZooNFT collection
  -  Withdraw can only be performed to NFT resources that exists
  -  A owned NFT can be transferred to another account that has CryptoZooNFT collection
  -  A single coordinate can only be minted once for each CryptoZooNFT land
  -  Increase mint limit for a give template

- `crypt-zoo-minter.test.js`
  - Deploy CryptoZooNFTMinter
  - Purchase fails when trying to purchase a CryptoZooNFT land
  - Purchase fails when insufficient funds
  - Purchase with minter

- `nft-storefront.test.js`
  - Deploy NFTStorefront
  - Setup NFTStorefront on any account
  - Create a listing
  - Accept a listing
  - Remove a listing
## 