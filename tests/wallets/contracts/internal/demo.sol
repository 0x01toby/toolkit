// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.11;

contract MyTest {
    function forceRevert(address payable to) public {
        to.transfer(0.05 ether);

        require(1!=1, "force revert");
    }

    function forceRevert2(address payable to) public {
        to.transfer(0.03 ether);

        require(1!=1, "force revert");
    }

    receive() external payable {

    }

    function partialFailure(address payable to) payable public {

        to.transfer(0.01 ether);
        try this.forceRevert(to){
        } catch {
            to.transfer(0.02 ether);
        }

        try this.forceRevert(to){
        } catch {
            to.transfer(0.02 ether);
        }


        to.transfer(0.03 ether);
    }
}