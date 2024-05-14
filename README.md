# ERC20 transfer statistics calculator

Implemented method `/api/top_active` which returns top 5 most active addresses from last 100 blocks. 
Addresses activity calculated by count of ERC20 token transfers. For statistics used transferred and received tokens.

Service process each block transactions. If found ERC20 transaction source address adds to statistics. 
To determine target address service parse transaction input, determine target address from method call.  

Implemented rate limiter to support 60 RPC (limitations of GetBlock service).

## Config

Required environment variables:
```
GET_BLOCK_API_KEY=<API KEY>>
HTTP_SERVER_HOST=<SERVER LISTEN HOST>
HTTP_SERVER_PORT=<SERVER LISTEN PORT>
```

## Limitations

Supported `transfer(address _to, uint256 _value)` method. 

Not supported:
- `transferFrom(address _from, address _to, uint256 _value)`
- batched transfers in one transaction