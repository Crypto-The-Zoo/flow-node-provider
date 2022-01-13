import { deployContractByName, executeScript, mintFlow, sendTransaction } from "flow-js-testing";

import { getCryptoZooAdminAddress } from "./common";

export const nftTemplates = {
  nft: {
    typeId: 1,
    isPack: false,
    name: 'regular NFT',
    description: 'A regular NFT',
    mintLimit: 2,
    priceUSD: 10.00000000,
    priceFlow: 1.00000000,
    metadata: {
      mintAttributes: 'AnimalLover',
      uri: '456@zoo.com/images/nft.png',
      mimetype: 'image/png',
      cid: '123',
    },
    timestamps: {
      availableAt: "1000.00000000",
      expiresAt: "5000000000.00000000"
    },
    isLand: false,
  },
  land: {
    typeId: 2,
    isPack: false,
    name: 'land NFT',
    nftName: 'A 20x20 land NFT located at [1, 1]',
    description: 'A land NFT',
    mintLimit: 2,
    priceUSD: 100.00000000,
    priceFlow: 10.00000000,
    metadata: {
      mintAttributes: 'AnimalLover',
      uri: '456@zoo.com/images/land.png',
      mimetype: 'image/png',
      cid: '456',
    },
    timestamps: {
      availableAt: "1000.00000000",
      expiresAt: "5000000000.00000000"
    },
    isLand: true,
  },
}

export const landCoords = {
  1: [0, 0],
  2: [10, 10],
  3: [100, 100],
}

/*
 * Deploys NonFungibleToken and CryptoZoo contracts to cryptoZooAdmin.
 * @throws Will throw an error if transaction is reverted.
 * @returns {Promise<*>}
 * */
export const deployCryptoZooContracts = async () => {
	const cryptoZooAdmin = await getCryptoZooAdminAddress();
	await mintFlow(cryptoZooAdmin, "10.0");

	await deployContractByName({ to: cryptoZooAdmin, name: "NonFungibleToken" });
  await deployContractByName({ to: cryptoZooAdmin, name: "MetadataViews" });

	const addressMap = { NonFungibleToken: cryptoZooAdmin, MetadataViews: cryptoZooAdmin };
	return deployContractByName({ to: cryptoZooAdmin, name: "CryptoZooNFT", addressMap });
};

/*
 * Setups CryptoZooNFT collection on account and exposes public capability.
 * @param {string} account - account address
 * @throws Will throw an error if transaction is reverted.
 * @returns {Promise<*>}
 * */
export const setupCryptoZooNFTOnAccount = async (account) => {
	const name = "CryptoZoo/setup_account";
	const signers = [account];

	return sendTransaction({ name, signers });
};

/*
 * Returns CryptoZooNFT supply.
 * @throws Will throw an error if execution will be halted
 * @returns {UInt64} - number of NFT minted so far
 * */
export const getCryptoZooSupply = async () => {
	const name = "CryptoZoo/get_crypto_zoo_supply";

	return executeScript({ name });
};

export const getCryptoZooMetadata = async (ownerAddress, id) => {
  const name = "CryptoZoo/get_nft_metadata";
  const args = [ownerAddress, id];
  return executeScript({ name, args });
}

export const getCryptoZooTemplate = async(typeId) => {
  const name = "CryptoZoo/get_crypto_zoo_template";
  const args = [typeId];
  return executeScript({ name, args });
}

export const checkIsCoordMinted = async (typeId, coord) => {
  const name = "CryptoZoo/check_is_coord_minted";
  const args = [typeId, coord];
  return executeScript({ name, args });
}

export const getCryptoZooItemCount = async (account) => {
	const name = "CryptoZoo/get_collection_length";
	const args = [account];

	return executeScript({ name, args });
};

export const getCryptoZooItem = async (account, itemID) => {
	const name = "CryptoZoo/get_crypto_zoo_item";
	const args = [account, itemID];

	return executeScript({ name, args });
};

export const createCryptoZooNFTTemplates = async() => {
  await createCryptoZooNFTTemplate(
    nftTemplates.nft.typeId,
    nftTemplates.nft.isPack,
    nftTemplates.nft.name,
    nftTemplates.nft.description,
    nftTemplates.nft.mintLimit,
    nftTemplates.nft.priceUSD,
    nftTemplates.nft.priceFlow,
    nftTemplates.nft.metadata,
    nftTemplates.nft.timestamps,
    nftTemplates.nft.isLand,
  );

  return await createCryptoZooNFTTemplate(
    nftTemplates.land.typeId,
    nftTemplates.land.isPack,
    nftTemplates.land.name,
    nftTemplates.land.description,
    nftTemplates.land.mintLimit,
    nftTemplates.land.priceUSD,
    nftTemplates.land.priceFlow,
    nftTemplates.land.metadata,
    nftTemplates.land.timestamps,
    nftTemplates.land.isLand,
  );
}

export const createCryptoZooNFTTemplate = async (typeId, isPack, nftName, description, mintLimit, priceUSD, priceFlow, metadata, timestamps, isLand) => {
	const cryptoZooAdmin = await getCryptoZooAdminAddress();

	const name = "CryptoZoo/create_nft_template";
	const args = [typeId, isPack, nftName, description, mintLimit, priceUSD, priceFlow, metadata, timestamps, isLand];
	const signers = [cryptoZooAdmin];

	return sendTransaction({ name, args, signers });
};

export const mintCryptoZooNFT = async (recipient, typeId) => {
	const cryptoZooAdmin = await getCryptoZooAdminAddress();

	const name = "CryptoZoo/mint_crypto_zoo_nft";
	const args = [recipient, typeId];
	const signers = [cryptoZooAdmin];

	return sendTransaction({ name, args, signers });
};

export const mintCryptoZooLand = async (recipient, typeId, coord, nftName) => {
	const cryptoZooAdmin = await getCryptoZooAdminAddress();

	const name = "CryptoZoo/mint_crypto_zoo_land";
	const args = [recipient, typeId, coord, nftName];
	const signers = [cryptoZooAdmin];

	return sendTransaction({ name, args, signers });
};

export const transferCryptoZooNFT = async (sender, recipient, itemId) => {
	const name = "CryptoZoo/transfer_crypto_zoo_nft";
	const args = [recipient, itemId];
	const signers = [sender];

	return sendTransaction({ name, args, signers });
};

export const increaseCryptoZooNFTMintLimit = async (typeId, mintLimit) => {
  const cryptoZooAdmin = await getCryptoZooAdminAddress();

	const name = "CryptoZoo/increase_mint_limit";
	const args = [typeId, mintLimit];
	const signers = [cryptoZooAdmin];

	return sendTransaction({ name, args, signers });
};