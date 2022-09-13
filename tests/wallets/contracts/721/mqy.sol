// SPDX-License-Identifier: MIT
pragma solidity 0.8.11;

import "./src/contracts/tokens/nf-token-metadata.sol";
import "./src/contracts/ownership/ownable.sol";

contract MqyFT is NFTokenMetadata, Ownable {

    constructor() {
        nftName = "MQY NFT";
        nftSymbol = "MQY";
    }

    function mint(address _to, uint256 _tokenId, string calldata _uri) external onlyOwner {
        super._mint(_to, _tokenId);
        super._setTokenUri(_tokenId, _uri);
    }

}