import { deployContractByName, executeScript, sendTransaction } from "flow-js-testing";
import { getCryptoZooAdminAddress } from "./common";
import { deployCryptoZooContracts, createCryptoZooNFTTemplates, setupCryptoZooNFTOnAccount } from "./crypto-zoo";

export const deployNFTStorefront = async () => {
	const CryptoZooAdmin = await getCryptoZooAdminAddress();

	await deployCryptoZooContracts();
  await createCryptoZooNFTTemplates();

	const addressMap = {
		NonFungibleToken: CryptoZooAdmin,
		CryptoZooNFT: CryptoZooAdmin,
	};

	return deployContractByName({ to: CryptoZooAdmin, name: "NFTStorefront", addressMap });
};

export const setupStorefrontOnAccount = async (account) => {
	// Account shall be able to store CryptoZoo NFTs
	await setupCryptoZooNFTOnAccount(account);

	const name = "NFTStorefront/setup_account";
	const signers = [account];

	return sendTransaction({ name, signers });
};

export const createItemListing = async (seller, itemId, price) => {
	const name = "NFTStorefront/create_listing";
	const args = [itemId, price];
	const signers = [seller];

	return sendTransaction({ name, args, signers });
};

export const purchaseItemListing = async (buyer, resourceId, seller) => {
	const name = "NFTStorefront/purchase_listing";
	const args = [resourceId, seller];
	const signers = [buyer];

	return sendTransaction({ name, args, signers });
};

export const removeItemListing = async (owner, itemId) => {
	const name = "NFTStorefront/remove_listing";
	const signers = [owner];
	const args = [itemId];

	return sendTransaction({ name, args, signers });
};

export const getListingCount = async (account) => {
	const name = "NFTStorefront/get_listings_length";
	const args = [account];

	return executeScript({ name, args });
};
