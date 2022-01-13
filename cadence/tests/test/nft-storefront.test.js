import path from "path";

import { 
	emulator,
	init,
	getAccountAddress,
	shallPass,
	mintFlow,
} from "flow-js-testing";

import { toUFix64 } from "../src/common";
import {
  nftTemplates,
	landCoords,
  
	mintCryptoZooLand,
  getCryptoZooItemCount,
  getCryptoZooItem,
} from "../src/crypto-zoo";
import {
	deployNFTStorefront,
	purchaseItemListing,
	createItemListing,
	removeItemListing,
	setupStorefrontOnAccount,
	getListingCount,
} from "../src/nft-storefront";

// We need to set timeout for a higher number, because some transactions might take up some time
jest.setTimeout(500000);

describe("NFT Storefront", () => {
	beforeEach(async () => {
		const basePath = path.resolve(__dirname, "../../");
		const port = 7003;
		await init(basePath, { port });
		return emulator.start(port, false);
	});

	// Stop emulator, so it could be restarted
	afterEach(async () => {
		return emulator.stop();
	});

	it("shall deploy NFTStorefront contract", async () => {
		await shallPass(deployNFTStorefront());
	});

	it("shall be able to create an empty Storefront", async () => {
		// Setup
		await deployNFTStorefront();
		const AnimalLover = await getAccountAddress("AnimalLover");

		await shallPass(setupStorefrontOnAccount(AnimalLover));
	});

	it("shall be able to create a listing", async () => {
		// Setup
		await deployNFTStorefront();
		const AnimalLover = await getAccountAddress("AnimalLover");
		await setupStorefrontOnAccount(AnimalLover);

		// Mint CryptoZoo Land for AnimalLover's account
		await shallPass(mintCryptoZooLand(AnimalLover, nftTemplates.land.typeId, landCoords[1], nftTemplates.land.nftName));

		const itemID = 0;

		await shallPass(createItemListing(AnimalLover, itemID, toUFix64(1.11)));
	});

	it("shall be able to accept a listing", async () => {
		// Setup
		await deployNFTStorefront();

		// Setup seller account
		const AnimalLover = await getAccountAddress("AnimalLover");
		await setupStorefrontOnAccount(AnimalLover);
		await mintCryptoZooLand(AnimalLover, nftTemplates.land.typeId, landCoords[1], nftTemplates.land.nftName);

		const itemId = 0;

		// Setup buyer account
		const AnotherAnimalLover = await getAccountAddress("AnotherAnimalLover");
		await setupStorefrontOnAccount(AnotherAnimalLover);

		await shallPass(mintFlow(AnotherAnimalLover, toUFix64(100)));

		// AnotherAnimalLover shall be able to buy from AnimalLover
		const sellItemTransactionResult = await shallPass(createItemListing(AnimalLover, itemId, toUFix64(1.11)));

		const listingAvailableEvent = sellItemTransactionResult.events[0];
		const listingResourceID = listingAvailableEvent.data.listingResourceID;

		await shallPass(purchaseItemListing(AnotherAnimalLover, listingResourceID, AnimalLover));

		const itemCount = await getCryptoZooItemCount(AnotherAnimalLover);
		expect(itemCount).toBe(1);

		const listingCount = await getListingCount(AnimalLover);
		expect(listingCount).toBe(0);
	});

	it("shall be able to remove a listing", async () => {
		// Deploy contracts
		await shallPass(deployNFTStorefront());

		// Setup AnimalLover account
		const AnimalLover = await getAccountAddress("AnimalLover");
		await shallPass(setupStorefrontOnAccount(AnimalLover));

		// Mint instruction shall pass
		await shallPass(mintCryptoZooLand(AnimalLover, nftTemplates.land.typeId, landCoords[1], nftTemplates.land.nftName));

		const itemId = 0;

		await getCryptoZooItem(AnimalLover, itemId);

		// Listing item for sale shall pass
		const sellItemTransactionResult = await shallPass(createItemListing(AnimalLover, itemId, toUFix64(1.11)));

		const listingAvailableEvent = sellItemTransactionResult.events[0];
		const listingResourceID = listingAvailableEvent.data.listingResourceID;

		// AnimalLover shall be able to remove item from sale
		await shallPass(removeItemListing(AnimalLover, listingResourceID));

		const listingCount = await getListingCount(AnimalLover);
		expect(listingCount).toBe(0);
	});
});