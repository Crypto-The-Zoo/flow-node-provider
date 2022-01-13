import path from "path";

import {
	emulator,
	init,
	getAccountAddress,
	mintFlow,
	shallPass,
	shallThrow,
	shallResolve,
	shallRevert,
} from "flow-js-testing";

import { nftTemplates, getCryptoZooItem } from "../src/crypto-zoo";

import { deployCryptoZooMinterContract, externalMintCryptoZooNFT } from "../src/cryptoo-zoo-minter";

// We need to set timeout for a higher number, because some transactions might take up some time
jest.setTimeout(50000);

describe("CryptoZooNFT Contract", () => {
	// Instantiate emulator and path to Cadence files
	beforeEach(async () => {
		const basePath = path.resolve(__dirname, "../../");
		const port = 7004;
		await init(basePath, { port });
		return emulator.start(port, false);
	});

	// Stop emulator, so it could be restarted
	afterEach(async () => {
		return emulator.stop();
	});

	it("shall deploy CryptoZooNFTMinter contract", async () => {
		await deployCryptoZooMinterContract();
	});

	it("Shall revert if trying to purchase land", async () => {
		await deployCryptoZooMinterContract();
		const AnimalLover = await getAccountAddress("AnimalLover");
		shallRevert(externalMintCryptoZooNFT(nftTemplates.land.typeId, AnimalLover));
	});

	it("Shall revert if insufficient funds", async () => {
		await deployCryptoZooMinterContract();
		const AnimalLover = await getAccountAddress("AnimalLover");
		shallRevert(externalMintCryptoZooNFT(nftTemplates.nft.typeId, AnimalLover));
	});

	it("Shall be able to purchase an NFT with Flow", async () => {
		await deployCryptoZooMinterContract();
		const AnimalLover = await getAccountAddress("AnimalLover");
		await mintFlow(AnimalLover, "10.0");

		await shallPass(externalMintCryptoZooNFT(nftTemplates.nft.typeId, AnimalLover));

		await shallResolve(async () => {
			const nft = await getCryptoZooItem(AnimalLover, 0);
			expect(nft.typeID).toBe(nftTemplates.nft.typeId);
		});
	});
});
