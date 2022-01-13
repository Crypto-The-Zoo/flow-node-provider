import path from "path";

import { emulator, init, getAccountAddress, shallPass, shallThrow, shallResolve, shallRevert } from "flow-js-testing";

import { getCryptoZooAdminAddress } from "../src/common";
import {
	// static
	nftTemplates,
	landCoords,
	// transactions
	deployCryptoZooContracts,
	mintCryptoZooNFT,
	mintCryptoZooLand,
	setupCryptoZooNFTOnAccount,
	createCryptoZooNFTTemplates,
	increaseCryptoZooNFTMintLimit,
	// scripts
	getCryptoZooSupply,
	getCryptoZooTemplate,
	checkIsCoordMinted,
	getCryptoZooItemCount,
	transferCryptoZooNFT,
  getCryptoZooMetadata,
} from "../src/crypto-zoo";

// We need to set timeout for a higher number, because some transactions might take up some time
jest.setTimeout(50000);

describe("CryptoZooNFT Contract", () => {
	// Instantiate emulator and path to Cadence files
	beforeEach(async () => {
		const basePath = path.resolve(__dirname, "../../");
		const port = 7002;
		await init(basePath, { port });
		return emulator.start(port, false);
	});

	// Stop emulator, so it could be restarted
	afterEach(async () => {
		return emulator.stop();
	});

	it("shall deploy CryptoZooNFT contract", async () => {
		await deployCryptoZooContracts();
	});

	it("supply shall be 0 after contract is deployed", async () => {
		// Setup
		await deployCryptoZooContracts();
		const cryptoZooAdmin = await getCryptoZooAdminAddress();
		await shallPass(setupCryptoZooNFTOnAccount(cryptoZooAdmin));

		await shallResolve(async () => {
			const supply = await getCryptoZooSupply();
			expect(supply).toBe(0);
		});
	});

	it("shall be able to create NFT templates", async () => {
		// Setup
		await deployCryptoZooContracts();
		await shallPass(createCryptoZooNFTTemplates());

		await shallResolve(async () => {
			const regularNFTTemplate = await getCryptoZooTemplate(nftTemplates.nft.typeId);
			expect(regularNFTTemplate.typeID).toBe(nftTemplates.nft.typeId);
			expect(regularNFTTemplate.name).toBe(nftTemplates.nft.name);
			expect(regularNFTTemplate.isPack).toBe(nftTemplates.nft.isPack);
			expect(regularNFTTemplate.priceUSD).toBe(nftTemplates.nft.priceUSD.toFixed(8));
			expect(regularNFTTemplate.priceFlow).toBe(nftTemplates.nft.priceFlow.toFixed(8));
			expect(regularNFTTemplate.isLand).toBe(nftTemplates.nft.isLand);
			expect(regularNFTTemplate.metadata).toStrictEqual(nftTemplates.nft.metadata);
			expect(regularNFTTemplate.timestamps).toStrictEqual(nftTemplates.nft.timestamps);
			expect(regularNFTTemplate.isExpired).toBe(false);
			expect(regularNFTTemplate.coordMinted).toStrictEqual({});
		});

		await shallResolve(async () => {
			const landNFTTemplate = await getCryptoZooTemplate(nftTemplates.land.typeId);
			expect(landNFTTemplate.typeID).toBe(nftTemplates.land.typeId);
			expect(landNFTTemplate.name).toBe(nftTemplates.land.name);
			expect(landNFTTemplate.isPack).toBe(nftTemplates.land.isPack);
			expect(landNFTTemplate.priceUSD).toBe(nftTemplates.land.priceUSD.toFixed(8));
			expect(landNFTTemplate.priceFlow).toBe(nftTemplates.land.priceFlow.toFixed(8));
			expect(landNFTTemplate.isLand).toBe(nftTemplates.land.isLand);
			expect(landNFTTemplate.metadata).toStrictEqual(nftTemplates.land.metadata);
			expect(landNFTTemplate.timestamps).toStrictEqual(nftTemplates.land.timestamps);
			expect(landNFTTemplate.isExpired).toBe(false);
			expect(landNFTTemplate.coordMinted).toStrictEqual({});
		});
	});

	it("shall be able to mint a CryptoZooNFT", async () => {
		// Setup
		await deployCryptoZooContracts();
		await createCryptoZooNFTTemplates();

		const AnimalLover = await getAccountAddress("AnimalLover");
		await setupCryptoZooNFTOnAccount(AnimalLover);

		// Mint instruction for AnimalLover account shall be resolved
		await shallPass(mintCryptoZooNFT(AnimalLover, nftTemplates.nft.typeId));
	});

  it("shall be able to resolve metadata of a CryptoZooNFT", async () => {
		// Setup
		await deployCryptoZooContracts();
		await createCryptoZooNFTTemplates();

		const AnimalLover = await getAccountAddress("AnimalLover");
		await setupCryptoZooNFTOnAccount(AnimalLover);

		// Mint instruction for AnimalLover account shall be resolved
		await shallPass(mintCryptoZooNFT(AnimalLover, nftTemplates.nft.typeId));

    // Metadata shall be resolved
    await shallResolve(async () => {
      const resolvedMetadata = await getCryptoZooMetadata(AnimalLover, 0)

      expect(resolvedMetadata.name).toBe(nftTemplates.nft.name)
      expect(resolvedMetadata.description).toBe(nftTemplates.nft.description)
      expect(resolvedMetadata.uri).toBe(nftTemplates.nft.metadata.uri)
      expect(resolvedMetadata.cid).toBe(nftTemplates.nft.metadata.cid)
      expect(resolvedMetadata.mimetype).toBe(nftTemplates.nft.metadata.mimetype)
      expect(resolvedMetadata.owner).toBe(AnimalLover)
      expect(resolvedMetadata.type).toBe('A.01cf0e2f2f715450.CryptoZooNFT.NFT')
    })
	});

	it("shall be able to mint a CryptoZooNFT Land", async () => {
		// Setup
		await deployCryptoZooContracts();
		await createCryptoZooNFTTemplates();

		const AnimalLover = await getAccountAddress("AnimalLover");
		await setupCryptoZooNFTOnAccount(AnimalLover);

		// NFT template shall not have coords marked as minted before mint
		await shallResolve(async () => {
			const land1IsMinted = await checkIsCoordMinted(nftTemplates.land.typeId, landCoords[1]);
			const land2IsMinted = await checkIsCoordMinted(nftTemplates.land.typeId, landCoords[2]);
			expect(land1IsMinted).toBe(false);
			expect(land2IsMinted).toBe(false);
		});

		// Mint instruction for AnimalLover account shall be resolved
		await shallPass(mintCryptoZooLand(AnimalLover, nftTemplates.land.typeId, landCoords[1], nftTemplates.land.nftName));
		await shallPass(mintCryptoZooLand(AnimalLover, nftTemplates.land.typeId, landCoords[2], nftTemplates.land.nftName));

		// NFT template shall have coords marked as minted
		await shallResolve(async () => {
			const land1IsMinted = await checkIsCoordMinted(nftTemplates.land.typeId, landCoords[1]);
			const land2IsMinted = await checkIsCoordMinted(nftTemplates.land.typeId, landCoords[2]);
			expect(land1IsMinted).toBe(true);
			expect(land2IsMinted).toBe(true);
		});

    // Metadata shall be resolved
    await shallResolve(async () => {
      const resolvedMetadata = await getCryptoZooMetadata(AnimalLover, 0)
      expect(resolvedMetadata.name).toBe(nftTemplates.land.nftName)
      expect(resolvedMetadata.description).toBe(nftTemplates.land.description)
      expect(resolvedMetadata.uri).toBe(nftTemplates.land.metadata.uri)
      expect(resolvedMetadata.cid).toBe(nftTemplates.land.metadata.cid)
      expect(resolvedMetadata.mimetype).toBe(nftTemplates.land.metadata.mimetype)
      expect(resolvedMetadata.owner).toBe(AnimalLover)
      expect(resolvedMetadata.type).toBe('A.01cf0e2f2f715450.CryptoZooNFT.NFT')
    })
	});

	it("shall be able to create a new empty NFT Collection", async () => {
		// Setup
		await deployCryptoZooContracts();
		const AnimalLover = await getAccountAddress("AnimalLover");
		await setupCryptoZooNFTOnAccount(AnimalLover);

		// shall be able te read AnimalLover collection and ensure it's empty
		await shallResolve(async () => {
			const itemCount = await getCryptoZooItemCount(AnimalLover);
			expect(itemCount).toBe(0);
		});
	});

	it("shall not be able to withdraw an NFT that doesn't exist in a collection", async () => {
		// Setup
		await deployCryptoZooContracts();
		const AnimalLover = await getAccountAddress("AnimalLover");
		const AnotherAnimalLover = await getAccountAddress("AnotherAnimalLover");
		await setupCryptoZooNFTOnAccount(AnimalLover);
		await setupCryptoZooNFTOnAccount(AnotherAnimalLover);

		// Transfer transaction shall fail for non-existent item
		await shallRevert(transferCryptoZooNFT(AnimalLover, AnotherAnimalLover, 1337));
	});

	it("shall be able to withdraw an NFT and deposit to another accounts collection", async () => {
		await deployCryptoZooContracts();
		await shallPass(createCryptoZooNFTTemplates());
		const AnimalLover = await getAccountAddress("AnimalLover");
		const AnotherAnimalLover = await getAccountAddress("AnotherAnimalLover");
		await setupCryptoZooNFTOnAccount(AnimalLover);
		await setupCryptoZooNFTOnAccount(AnotherAnimalLover);

		// Land Mint instruction for AnimalLover account shall be resolved
		await shallPass(mintCryptoZooLand(AnimalLover, nftTemplates.land.typeId, landCoords[1], nftTemplates.land.nftName));

		// Transfer land transaction shall pass
		await shallPass(transferCryptoZooNFT(AnimalLover, AnotherAnimalLover, 0));
	});

	it("shall not be able to mint the same coord twice", async () => {
		// Setup
		await deployCryptoZooContracts();
		await createCryptoZooNFTTemplates();

		const AnimalLover = await getAccountAddress("AnimalLover");
		await setupCryptoZooNFTOnAccount(AnimalLover);

		// Mint instruction for AnimalLover account shall be resolved
		await shallPass(mintCryptoZooLand(AnimalLover, nftTemplates.land.typeId, landCoords[1], nftTemplates.land.nftName));
		await shallRevert(mintCryptoZooLand(AnimalLover, nftTemplates.land.typeId, landCoords[1], nftTemplates.land.nftName));
	});

	it("template shall expire when mint limit is hit", async () => {
		// Setup
		await deployCryptoZooContracts();
		await createCryptoZooNFTTemplates();

		const AnimalLover = await getAccountAddress("AnimalLover");
		await setupCryptoZooNFTOnAccount(AnimalLover);

		// Mint instruction for AnimalLover account shall be resolved
		await shallPass(mintCryptoZooLand(AnimalLover, nftTemplates.land.typeId, landCoords[1], nftTemplates.land.nftName));
		await shallPass(mintCryptoZooLand(AnimalLover, nftTemplates.land.typeId, landCoords[2], nftTemplates.land.nftName));

		await shallRevert(mintCryptoZooLand(AnimalLover, nftTemplates.land.typeId, landCoords[3]));

		await shallResolve(async () => {
			const landNFTTemplate = await getCryptoZooTemplate(nftTemplates.land.typeId);
			expect(landNFTTemplate.isExpired).toBe(true);
		});
	});

	it("shall be able to increase template mint limit, which unexpires the template", async () => {
		// Setup
		await deployCryptoZooContracts();
		await createCryptoZooNFTTemplates();

		const AnimalLover = await getAccountAddress("AnimalLover");
		await setupCryptoZooNFTOnAccount(AnimalLover);

		// Mint instruction for AnimalLover account shall be resolved
		await shallPass(mintCryptoZooLand(AnimalLover, nftTemplates.land.typeId, landCoords[1], nftTemplates.land.nftName));
		await shallPass(mintCryptoZooLand(AnimalLover, nftTemplates.land.typeId, landCoords[2], nftTemplates.land.nftName));

		// Mint instruction for AnimalLover account shall be reverted when mintLimit is hit
		await shallRevert(mintCryptoZooLand(AnimalLover, nftTemplates.land.typeId, landCoords[3], nftTemplates.land.nftName));

		// Shall be able to unexpire template by increasing mint limit
		await shallPass(increaseCryptoZooNFTMintLimit(nftTemplates.land.typeId, 3));

		// Mint instruction for AnimalLover account shall be resolved
		await shallResolve(mintCryptoZooLand(AnimalLover, nftTemplates.land.typeId, landCoords[3], nftTemplates.land.nftName));
	});
});
