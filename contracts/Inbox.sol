pragma solidity >=0.5.2 <0.6.0;

contract Inbox{
    string public testMessage;

    constructor(string memory initialMessage) public{
        testMessage = initialMessage;
    }
    function setMessage(string memory newTestMessage) public{
        testMessage = newTestMessage;
    }
}