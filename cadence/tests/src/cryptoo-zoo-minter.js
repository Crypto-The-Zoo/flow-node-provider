import { deployContractByName, executeScript, mintFlow, sendTransaction } from "flow-js-testing";

import { getCryptoZooAdminAddress } from "./common";

import { deployCryptoZooContracts, createCryptoZooNFTTemplates } from "./crypto-zoo";

export const deployCryptoZooMinterContract = async () => {
	const cryptoZooAdmin = await getCryptoZooAdminAddress();
	await mintFlow(cryptoZooAdmin, "10.0");

  await deployCryptoZooContracts()
  await createCryptoZooNFTTemplates()

  const addressMap = {
    FungibleToken: "0xee82856bf20e2aa6",
    FlowToken: "0x0ae53cb6e3f42a79",
    NonFungibleToken: cryptoZooAdmin,
    CryptoZooNFT: cryptoZooAdmin,
  };
	return deployContractByName({ to: cryptoZooAdmin, name: "CryptoZooNFTMinter", addressMap });
};


export const externalMintCryptoZooNFT = async (typeId, account) => {

  const name = "CryptoZoo/mint_crypto_zoo_nft_with_flow";
  const args = [typeId];
  const signers = [account];

  return sendTransaction({ name, args, signers });
};
